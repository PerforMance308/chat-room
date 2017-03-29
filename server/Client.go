package main

import (
	"net"
	"sync"
	"encoding/binary"
	"time"
	"bytes"
	"server/proto_struct"
	"github.com/golang/protobuf/proto"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	conn, _ := net.Dial("tcp", "127.0.0.1:8080")

	param := &proto_struct.RoleLoginC2S{
		Account: proto.String("test"),
	}

	data := make([]byte, 2)
	msgData, _ := proto.Marshal(param)

	binary.BigEndian.PutUint16(data, uint16(1))
	data = append(data, msgData...)

	time.Sleep(time.Second * 2)
	conn.Write(data)

	wg.Wait()
}

func praseData(data []byte) (uint16, []byte) {
	b_buf := bytes.NewBuffer(data[:2])

	var mId uint16
	binary.Read(b_buf, binary.BigEndian, &mId)

	return mId, data[2:]
}
