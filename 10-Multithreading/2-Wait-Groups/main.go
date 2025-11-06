package main

import (
	"fmt"
	"sync"
	"time"
)

// Thread 1
func main() {
	waitgroups := sync.WaitGroup{}
	waitgroups.Go(func() {
		for i := range 10 {
			fmt.Printf("%d: Task %s is running\n", i, "A")
			time.Sleep(1 * time.Second)
		}
	})
	// Thread 2
	waitgroups.Go(func() {
		for i := range 10 {
			fmt.Printf("%d: Task %s is running\n", i, "B")
			time.Sleep(1 * time.Second)
		}
	})
	// Thread 4
	waitgroups.Go(func() {
		for i := range 5 {
			fmt.Printf("%d: Task %s is running, poha\n", i, "anonymous")
			time.Sleep(1 * time.Second)
			waitgroups.Done()
		}
	})
	waitgroups.Wait()
}
