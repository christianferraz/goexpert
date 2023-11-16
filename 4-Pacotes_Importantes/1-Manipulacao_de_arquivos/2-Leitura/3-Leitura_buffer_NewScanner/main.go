package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("abc.txt")
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		linha := scanner.Bytes()
		fmt.Println(string(linha))
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
	}
}
