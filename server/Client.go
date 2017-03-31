package main

import (
	"net"
	"sync"
	"encoding/binary"
	"server/proto_struct"
	"github.com/golang/protobuf/proto"
	"fmt"
	"os"
	"bufio"
	"bytes"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	fmt.Println("What to send to the server? Type Q to quit.")
	inputReader := bufio.NewReader(os.Stdin)
	input, _ := inputReader.ReadString('\n')

	param := &proto_struct.RoleLoginC2S{
		Account: proto.String(input),
	}

	data := make([]byte, 2)
	msgData, _ := proto.Marshal(param)

	binary.BigEndian.PutUint16(data, uint16(1))
	data = append(data, msgData...)

	conn, _ := net.Dial("tcp", "127.0.0.1:8080")
	conn.Write(data)

	bytes := make([]byte, 32768)

	for{
		i, _ := conn.Read(bytes)

		recvdata := bytes[:i]
		mId, _ := praseData(recvdata)
		fmt.Println("recv msg:", mId)
	}

	wg.Wait()
}

func praseData(data []byte) (uint16, []byte) {
	b_buf := bytes.NewBuffer(data[:2])

	var mId uint16
	binary.Read(b_buf, binary.BigEndian, &mId)

	return mId, data[2:]
}