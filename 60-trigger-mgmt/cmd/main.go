package main

import (
	"time"

	"github.com/christianferraz/goexpert/60-trigger-mgmt/internal/trigger"
)

func main() {
	tm := trigger.NewTriggerManager(10, 1*time.Minute, "https://uol.com.br")

	for {
		tm.RegisterRequest()
		time.Sleep(1 * time.Second)
	}
}
