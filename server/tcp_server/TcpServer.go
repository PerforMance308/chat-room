package tcp_server

import (
	"sync"
	"reflect"
	"server/logger"
	"net"
	"fmt"
	"os"
	"github.com/golang/protobuf/proto"
	"server/options"
)

type Server struct {
	wg              *sync.WaitGroup
	mu              sync.Mutex
	clients         map[uint64]*Client
	requestHandlers map[uint16]reflect.Value
	listener        net.Listener
	running         bool
	gcid            uint64
	reqParamsType   map[uint16]reflect.Type
	config   options.TCPConf
}

func NewServer(conf options.TCPConf, wg *sync.WaitGroup) *Server {
	wg.Add(1)

	s := &Server{}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.clients = make(map[uint64]*Client)
	s.requestHandlers = make(map[uint16]reflect.Value)
	s.running = true
	s.reqParamsType = make(map[uint16]reflect.Type)
	s.config = conf

	logger.Logger().Debug("Init server")
	return s
}

func (s *Server) isRunning() bool {
	return s.running
}

func (s *Server) Start() {
	defer s.wg.Done()

	host := s.config.Host
	port := s.config.Port

	logger.Logger().Debug("TCP server will listen", host, "at", port)

	hp := fmt.Sprintf("%s:%d", host, port)
	if l, err := net.Listen("tcp", hp); err != nil {
		logger.Logger().Critical("Failed to listen:", hp+err.Error())
		os.Exit(-1)
	} else {
		s.mu.Lock()
		s.listener = l
		s.mu.Unlock()

		logger.Logger().Notice("TCP server started to accept connection.")

		for s.isRunning() {
			conn, err := l.Accept()

			if err != nil {
				logger.Logger().Error("Temporary Client accept Error ", err)
			}

			go s.handleConn(conn)
		}
	}
	logger.Logger().Info("Sever exiting")
}

func (s *Server) handleConn(conn net.Conn) {
	logger.Logger().Debug("server clients counter ", len(s.clients))
	c := initClient(conn, s)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.clients[c.id] = c
}

func (s *Server) removeClient(c *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, c.id)
}

func (s *Server) dispatch(client *Client, mId uint16, request []byte) {
	defer func() {
		if r := recover(); r != nil {
			logger.Logger().Error("Panic recovered from Dispatch Error",r)
		}
	}()

	rqParamsType := s.reqParamsType[mId]
	rqParams := reflect.New(rqParamsType).Interface()
	if err := proto.Unmarshal(request, rqParams.(proto.Message)); err != nil {
		return
	}

	if handler, ok := s.requestHandlers[mId]; ok {
		in := make([]reflect.Value, 2)
		in[0] = reflect.ValueOf(client)
		in[1] = reflect.ValueOf(rqParams).Elem()

		handler.Call(in)
	} else {
		logger.Logger().Notice(fmt.Sprintf("No Handler for the handle id:", mId))
	}
}

func (s *Server) AddRequestHandler(id uint16, handler interface{}) {
	if _, ok := s.requestHandlers[id]; ok {
		logger.Logger().Error(fmt.Sprintf("There were a handler for", id))
	} else {

		hId, hFun := ParseRequestHandlerFun(handler)

		s.reqParamsType[id] = hId
		s.requestHandlers[id] = hFun
	}
}

func ParseRequestHandlerFun(handler interface{}) (reflect.Type, reflect.Value) {
	hType := reflect.TypeOf(handler)
	if hType.Kind() != reflect.Func {
		panic("Request handler need be a func")
	}

	numArgs := hType.NumIn()

	var testFunc func(*Client)
	testType := reflect.TypeOf(testFunc)
	numIn := testType.NumIn()

	if numArgs != numIn+1 {
		panic(fmt.Sprintf("Request handler func need args", numIn+1))
	}

	for i := 0; i < numIn-1; i++ {
		if hType.In(i) != testType.In(i) {
			panic(fmt.Sprintf("Request handler func params", hType.In(i), "should be ", testType.In(i), "args"))
		}

	}

	return hType.In(numIn), reflect.ValueOf(handler)
}