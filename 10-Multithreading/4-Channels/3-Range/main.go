package main

import (
	"fmt"
	"time"
)

// Thread 1
func main() {
	ch := make(chan int)
	go publish(ch)
	go reader(ch)
	time.Sleep(5 * time.Minute)
	// for x := range ch {
	// 	fmt.Printf("Received %d\n", x)
	// }
}

func reader(ch chan int) {
	for x := range ch {
		fmt.Printf("Received %d\n", x)
	}
}

func publish(ch chan int) {
	for i := 0; i < 10; i++ {
		time.Sleep(3 * time.Second)
		ch <- i
	}
	//evitar o deadlock
	close(ch)
}
