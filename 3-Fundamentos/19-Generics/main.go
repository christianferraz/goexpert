package main

import (
	"fmt"
)

// criar constraints
type Number interface{ int | float64 }

// antes - func Soma(m map[string]int) int {
// func generic
// func Soma[T int | float64](m map[string]T) T {
// func generic com constraints
func Soma[T Number](m map[string]T) T {
	var soma T
	for _, v := range m {
		soma += v
	}
	return soma
}

func main() {
	m := map[string]float64{"a": 1.0, "b": 2.0, "c": 3.0}
	fmt.Println(Soma(m))
	fmt.Printf("O resultado da soma é: %v\n", Soma(m))

}
