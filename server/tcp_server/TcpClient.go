package tcp_server

import (
	"net"
	"sync"
	"sync/atomic"
	"server/logger"
	"time"
	"bytes"
	"encoding/binary"
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

type Client struct {
	id                    uint64
	uuid                  bson.ObjectId
	mu                    sync.RWMutex
	nc                    net.Conn
	sendBuff              chan *[]byte
	srv                   *Server
	lastDispatchedAt      int64
	lastDispatchedCounter int
	running               bool
}

func initClient(conn net.Conn, sev *Server) *Client {
	c := &Client{nc: conn, srv: sev}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.uuid = bson.NewObjectId()
	c.sendBuff = make(chan *[]byte)
	c.id = atomic.AddUint64(&sev.gcid, 1)
	c.running = true

	go c.readLoop()
	go c.writeLoop()

	return c
}

func (c *Client) CloseConnection(after time.Duration, wait bool) {
	defer func() {
		if r := recover(); r != nil {
			logger.Logger().Error("Panic recovered from CloseConnection")
		}
	}()
	if after > 0 {
		time.Sleep(after)
	}

	c.SetClosed()
	if c.srv != nil {
		c.srv.removeClient(c)
	}

	close(c.sendBuff)
	c.nc.Close()

	logger.Logger().Notice("The client", c.id, "close connection or timeout")
}

func (c *Client) readLoop() {
	defer func() {
		if r := recover(); r != nil {
			logger.Logger().Error("client read loop error:", r)
		}
	}()

	if c.netConn() == nil {
		return
	}

	bytes := make([]byte, 32768)

	for c.srv.isRunning() {
		if c.isRunning() {
			c.netConn().SetReadDeadline(time.Now().Add(5 * time.Minute))
			i, err := c.netConn().Read(bytes)

			if err != nil {
				if err.Error() == "EOF" {
					c.CloseConnection(0*time.Second, false)
					break
				}
			}
			data := bytes[:i]

			var e error

			if e == nil && len(data) > 0 {
				if c.lastDispatchedAt == time.Now().Unix() {
					c.lastDispatchedCounter += 1
					if c.lastDispatchedCounter >= 6 {
						time.Sleep(2 * time.Second)
					}
				} else {
					c.lastDispatchedAt = time.Now().Unix()
					c.lastDispatchedCounter = 0
				}

				mId, pack := praseData(data)
				go c.srv.dispatch(c, mId, pack)
			}
		}
	}
}

func (c *Client) isRunning() bool {
	return c.running
}

func (c *Client) SetClosed() {
	c.running = false
}

func (c *Client) writeLoop() {
	defer func() {
		if r := recover(); r != nil {
			logger.Logger().Error("client write loop error:", r)
		}
	}()

	for data := range c.sendBuff {

		if !c.isRunning() {
			return
		}

		c.send2Client(data)
	}

}

func (c *Client) Write(mId uint16, data []byte) {
	defer func() {
		if r := recover(); r != nil {
			logger.Logger().Error("client write loop error:", r)
		}
	}()

	idData := make([]byte, 2)

	binary.BigEndian.PutUint16(idData, uint16(mId))
	data = append(idData, data...)

	if c.isRunning() {
		c.WriteToSendBuff(&data)
	}
}

func (c *Client) WriteToSendBuff(data *[]byte) {
	if c.sendBuff != nil {
		c.sendBuff <- data
	}
}

func (c *Client) send2Client(data *[]byte) {
	if c.netConn() == nil {
		return
	}

	c.netConn().SetWriteDeadline(time.Now().Add(10 * time.Minute))

	if n, err := c.nc.Write(*data); err != nil {
		logger.Logger().Error(fmt.Sprintf("Write ", n, " bytes to nc error: ", err))
		go c.CloseConnection(0*time.Second, false)
	}
}

func (c *Client) netConn() net.Conn {
	return c.nc
}

func praseData(data []byte) (uint16, []byte) {
	b_buf := bytes.NewBuffer(data[:2])

	var mId uint16
	binary.Read(b_buf, binary.BigEndian, &mId)

	return mId, data[2:]
}
