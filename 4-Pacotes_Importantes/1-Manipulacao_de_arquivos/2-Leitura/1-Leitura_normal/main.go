package main

import (
	"fmt"
	"os"
)

func main() {

	//leitura
	arquivo, err := os.ReadFile("abc.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(arquivo))
}
