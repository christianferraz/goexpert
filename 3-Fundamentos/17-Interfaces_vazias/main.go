package main

import "fmt"

func main() {
	var x interface{} = 3.14
	showType(x)
	var y interface{} = "Olá mundo"
	showType(y)
}

func showType(t interface{}) {
	fmt.Printf("O tipo da variável é %T e o valor é %v\n", t, t)
}
