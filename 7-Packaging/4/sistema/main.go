package main

import (
	"github.com/christianferraz/goexpert/7-Packaging/4/math"
	"github.com/google/uuid"
)

func main() {
	m := math.NewMath(1, 2)
	println(m.Add())
	println(uuid.New().String())
}
