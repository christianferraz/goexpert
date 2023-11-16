package main

import (
	"fmt"
	"os"
)

func main() {
	f, error := os.Create("abc.txt")
	if error != nil {
		panic(error)
	}
	// gravar string
	// tam, err := f.WriteString("Hello World!")
	// gravar bytes
	tam, err := f.Write([]byte("Hello World!"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("o tamanho do arquivo é: %v bytes\n", tam)
	f.Close()

	//leitura
	arquivo, err := os.ReadFile("abc.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(arquivo))
}
