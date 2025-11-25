package main

import (
	"fmt"
	"runtime"
	"time"
)

func rotina(ch chan<- int) {
	ch <- 1
	fmt.Println("executou o 1")
	ch <- 2
	fmt.Println("executou o 2")
	ch <- 3
	fmt.Println("executou o 3")
	ch <- 4
	ch <- 5
	fmt.Println("Executou!")
	time.Sleep(time.Second * 5)
	ch <- 6
}
func main() {
	ch := make(chan int, 128)
	go rotina(ch)
	fmt.Println(<-ch)
	fmt.Println("recebeu o 1")
	time.Sleep(time.Second * 5)
	fmt.Println(<-ch)
	fmt.Println("recebeu o 2")
	time.Sleep(time.Second)
	fmt.Println(<-ch)
	fmt.Println("recebeu o 3")
	fmt.Println(<-ch)
	fmt.Println("recebeu o 4")
	fmt.Println(<-ch)
	fmt.Println("recebeu o 5")
	fmt.Println(<-ch)
	fmt.Println("recebeu o 6")
	fmt.Println("nr cpus ", runtime.NumCPU())
	close(ch)
}
