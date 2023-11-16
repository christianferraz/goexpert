package main

import "fmt"

func main() {
	var minhaVar interface{} = "Olá mundo"
	println(minhaVar.(string))
	res, ok := minhaVar.(int)
	fmt.Printf("O valor de res é %v, o valor de ok é %v", res, ok)
	// sem receber o ok, qdo der erro, vai gerar um panic
	res2 := minhaVar.(int)
	fmt.Printf("o valor de res2 é %v\n", res2)

}
