package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func FuzzJSONUnmarshal(f *testing.F) {
	// Casos iniciais
	f.Add(`{"nome":"Christian"}`)
	f.Add(`{"idade":30}`)

	f.Fuzz(func(t *testing.T, data string) {
		var v map[string]any
		// Testa se a lib consegue lidar sem panic
		_ = json.Unmarshal([]byte(data), &v)
		fmt.Println(v)
	})
}

// executar
// go test -fuzz=Fuzz
