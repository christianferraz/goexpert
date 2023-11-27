package main

import (
	"fmt"
	"os"
)

func main() {
	for i := 0; i < 100; i++ {
		f, err := os.Create(fmt.Sprintf("file-%d.txt", i))
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		f.WriteString("Arquivo de teste")
	}
}
