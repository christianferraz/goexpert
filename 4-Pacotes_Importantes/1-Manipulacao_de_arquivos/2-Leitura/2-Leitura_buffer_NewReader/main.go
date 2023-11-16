package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	//leitura
	arquivo, err := os.Open("abc.txt")

	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(arquivo)
	buffer := make([]byte, 10)
	for {
		n, err := reader.Read(buffer)
		if err != nil {
			break
		}
		// o :n é a posiçao da onde está fazendo a leitura
		fmt.Println(string(buffer[:n]))
	}
}
