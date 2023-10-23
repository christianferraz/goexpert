package main

import (
	"fmt"
	"reflect"
)

func main() {
	numeros := []int{1, 2, 3, 4, 5} // compilador conta!
	a1 := [3]int{1, 2, 3}           // array
	s1 := []int{1, 2, 3}            // slice
	fmt.Println(reflect.TypeOf(a1), reflect.TypeOf(s1))
	for i, numero := range numeros {
		fmt.Printf("%d) %d\n", i+1, numero)
	}
	for _, num := range numeros {
		fmt.Println(num)
	}
}
