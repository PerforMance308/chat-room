package tcp_server

import (
	"sync"
	"reflect"
	"logger"
	"net"
	"fmt"
	"os"
	"github.com/golang/protobuf/proto"
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
}

func NewServer(wg *sync.WaitGroup) *Server {
	wg.Add(1)

	s := &Server{}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.clients = make(map[uint64]*Client)
	s.requestHandlers = make(map[uint16]reflect.Value)
	s.running = true
	s.reqParamsType = make(map[uint16]reflect.Type)

	logger.Debug("Initialize server struct")
	return s
}

func (s *Server) isRunning() bool {
	return s.running
}

func (s *Server) Start() {
	defer s.wg.Done()

	host := "127.0.0.1"
	port := 8080
	logger.Debug("TCP server will listen %s at %d", host, port)

	hp := fmt.Sprintf("%s:%d", host, port)
	if l, err := net.Listen("tcp", hp); err != nil {
		logger.Critical("Failed to listen: %s", hp+err.Error())
		os.Exit(-1)
	} else {
		s.mu.Lock()
		s.listener = l
		s.mu.Unlock()

		logger.Notice("TCP server started to accept connection.")

		for s.isRunning() {
			conn, err := l.Accept()

			if err != nil {
				if ne, ok := err.(net.Error); ok && ne.Temporary() {
					logger.Error("Temporary Client accept Error %v", err)
					//time.Sleep(10)
				} else if s.isRunning() {
					logger.Critical("Accept error: %v", err)
				}
				continue
			}

			s.handleConn(conn)
		}
	}
	logger.Info("Sever exiting")
}

func (s *Server) handleConn(conn net.Conn) {
	logger.Debug("server clients counter %d", len(s.clients))
	c := initClient(conn, s)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.clients[c.id] = c

}

func (s *Server) removeClient(c *Client) {
	delete(s.clients, c.id)
}

func (s *Server) dispatch(client *Client, mId uint16, request []byte) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Panic recovered from Dispatch Error")
		}
	}()

	rqParamsType := s.reqParamsType[mId]
	rqParams := reflect.New(rqParamsType).Interface()
	if err := proto.Unmarshal(request, rqParams.(proto.Message)); err != nil {
		return
	}

	if handler, ok := s.requestHandlers[mId]; ok {
		logger.Notice(fmt.Sprintf("Dispatch a pack to handler %d.", mId))
		in := make([]reflect.Value, 2)
		in[0] = reflect.ValueOf(client)
		in[1] = reflect.ValueOf(rqParams).Elem()

		handler.Call(in)
	} else {
		logger.Notice(fmt.Sprintf("No Handler for the handle id: %d", mId))
	}
}

func (s *Server) AddRequestHandler(id uint16, handler interface{}) {
	if _, ok := s.requestHandlers[id]; ok {
		logger.Error(fmt.Sprintf("There were a handler for %d", id))
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
		panic(fmt.Sprintf("Request handler func need  %d args", numIn+1))
	}

	for i := 0; i < numIn-1; i++ {
		if hType.In(i) != testType.In(i) {
			panic(fmt.Sprintf("Request handler func params  %v should be %v args", hType.In(i), testType.In(i)))
		}

	}

	return hType.In(numIn), reflect.ValueOf(handler)
}