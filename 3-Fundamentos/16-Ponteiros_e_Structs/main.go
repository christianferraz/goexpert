package main

type Conta struct {
	saldo int
}

func NewConta() *Conta {
	return &Conta{saldo: 0}
}

func (c *Conta) Add(saldo int) int {
	c.saldo += saldo
	return c.saldo
}

func main() {
	c := NewConta()
	c.Add(100)
	c.Add(100)
	println(c.saldo)
}
