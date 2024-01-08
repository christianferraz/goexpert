package main

import (
	"fmt"

	"github.com/christianferraz/goexpert/2-Alura/1-ContaCorrente/contas"
)

func PagarBoleto(conta verificarConta, valorBoleto float64) {
	conta.Sacar(valorBoleto)
}

type verificarConta interface {
	Sacar(valor float64) string
}

func main() {
	contaChristian := contas.ContaPoupanca{}
	contaChristian.Depositar(100.3)
	PagarBoleto(&contaChristian, 60)
	fmt.Printf("%v\n", contaChristian.ObterSaldo())
}
