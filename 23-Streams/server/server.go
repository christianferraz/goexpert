package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

type StreamServer struct{}

func (ss *StreamServer) ConnectAndReadFile() {
	ln, err := net.Listen("tcp", ":8888")
	if err != nil {
		panic(err.Error())
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err.Error())
		}
		go ss.Process(conn)
	}
}

func (ss *StreamServer) Process(conn net.Conn) {
	// buf := make([]byte, 1024)
	buf := new(bytes.Buffer)
	var size int64
	for {
		err := binary.Read(conn, binary.LittleEndian, &size)
		if err != nil {
			panic(err.Error())
		}
		qtdBytes, err := io.CopyN(buf, conn, size)

		if err != nil {
			panic(err)
		}
		fmt.Println(buf.String())
		fmt.Println("Received", qtdBytes, "bytes from client")
		break
	}
}

func main() {
	ss := &StreamServer{}
	ss.ConnectAndReadFile()
}
