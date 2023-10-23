package main

import (
	"fmt"
	"sync"
	"time"
)

func task(name string, wg *sync.WaitGroup) {
	for i := 0; i < 10; i++ {
		fmt.Printf("%d: Task %s is running\n", i, name)
		time.Sleep(1 * time.Second)
		wg.Done()
	}
}

// Thread 1
func main() {
	waitgroups := sync.WaitGroup{}
	waitgroups.Add(25)
	// Thread 2
	go task("A", &waitgroups)
	// Thread 3
	go task("B", &waitgroups)
	// Thread 4
	go func() {
		for i := 0; i < 5; i++ {
			fmt.Printf("%d: Task %s is running, poha\n", i, "anonymous")
			time.Sleep(1 * time.Second)
			waitgroups.Done()
		}
	}()
	waitgroups.Wait()
}
