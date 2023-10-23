package main

import (
	"fmt"
	"time"
)

// Thread 1
func main() {
	canal := make(chan string) // Vazio

	// Thread 2
	go func() {
		time.Sleep(10 * time.Second)
		canal <- "Olá Mundo!" // Está cheio
	}()

	// Thread 1
	fmt.Println("antes do canal")
	msg := <-canal // Canal esvazia
	fmt.Println("antes do canal")
	fmt.Println(msg)
}
