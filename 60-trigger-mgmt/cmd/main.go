package main

import (
	"time"

	"github.com/christianferraz/goexpert/60-trigger-mgmt/internal/trigger"
)

func main() {
	tm := trigger.NewTriggerManager(20*time.Second, "https://10.42.0.32:30443/api/util/SyncAll")

	for {
		tm.RegisterRequest()
		tm.RegisterRequest()
		tm.RegisterRequest()
		tm.RegisterRequest()
		tm.RegisterRequest()
		tm.RegisterRequest()
		tm.RegisterRequest()
		tm.RegisterRequest()
		tm.RegisterRequest()
		tm.RegisterRequest()
		time.Sleep(1000 * time.Second)
		tm.RegisterRequest()
	}
}
