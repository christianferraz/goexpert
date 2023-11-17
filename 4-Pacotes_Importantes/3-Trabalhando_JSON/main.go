package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Conta struct {
	Numero int     `json:"n"`
	Saldo  float64 `json:"s"`
}

func main() {
	conta := Conta{123, 55.5}
	res, err := json.Marshal(conta)
	if err != nil {
		panic(err)
	}
	println(string(res))
	// trabalhando com Encoding
	// arq, err := os.Create("eu.txt")
	// if err != nil {
	// 	panic(err)
	// }
	// json.NewEncoder(arq).Encode(conta)
	json.NewEncoder(os.Stdout).Encode(conta)
	jsonPuro := `{"n":1233,"s":55.5}`
	var conta2 Conta
	err = json.Unmarshal([]byte(jsonPuro), &conta2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%.1f\n", conta2.Saldo)

}
