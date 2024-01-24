package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	file, err := os.ReadFile("client.go")
	if err != nil {
		panic(err)
	}

	conn, err := net.Dial("tcp", "10.42.0.64:8888")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()
	binary.Write(conn, binary.LittleEndian, int64(len(file)))
	qtsBytes, err := io.CopyN(conn, bytes.NewReader(file), int64(len(file)))
	if err != nil {
		panic(err)
	}
	fmt.Println("Sent", qtsBytes, "bytes to server")
}
