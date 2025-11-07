package main

import (
	"fmt"
	"runtime"
	"time"
)

func rotina(ch chan<- int) {
	ch <- 1
	ch <- 2
	ch <- 3
	ch <- 4
	ch <- 5
	fmt.Println("Executou!")
	time.Sleep(time.Second * 5)
	ch <- 6
}
func main() {
	ch := make(chan int, 6)
	go rotina(ch)
	fmt.Println(<-ch)
	fmt.Println(<-ch)
	time.Sleep(time.Second)
	fmt.Println(<-ch)
	fmt.Println(<-ch)
	fmt.Println(<-ch)
	fmt.Println(<-ch)
	fmt.Println("nr cpus ", runtime.NumCPU())
	close(ch)
}
