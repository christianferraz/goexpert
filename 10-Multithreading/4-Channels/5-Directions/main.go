package main

import (
	"fmt"
	"sync"
)

func envia(nome string, hello chan<- string, wg *sync.WaitGroup) {
	hello <- nome
	wg.Done()
}

func ler(data <-chan string) {
	fmt.Println(<-data)
}

func main() {
	wg := sync.WaitGroup{}
	hello := make(chan string)
	wg.Add(1)
	go envia("Hello", hello, &wg)
	ler(hello)
	wg.Wait()
}
