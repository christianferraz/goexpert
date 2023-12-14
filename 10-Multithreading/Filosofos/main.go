package main

import (
	"fmt"
	"sync"
	"time"
)

type Filosofo struct {
	id                  int
	leftFork, rightFork *sync.Mutex
	refeicoes           int
}

const Refeicoes_desejadas = 1
const Num_filosofos = 5

func (f *Filosofo) comer() {
	fmt.Printf("Filosofo %d comendo\n", f.id)
	f.refeicoes++
	time.Sleep(time.Duration(30) * time.Second) // Tempo aleatório para comer
	fmt.Printf("Filosofo %d terminou de comer\n", f.id)
}

func (f *Filosofo) pensar() {
	fmt.Printf("Filosofo %d pensando\n", f.id)
	time.Sleep(time.Duration(5) * time.Second) // Tempo aleatório para pensar
}

func (f *Filosofo) tentarComer() {
	for f.refeicoes < Refeicoes_desejadas {
		f.pensar()

		f.leftFork.Lock()
		fmt.Printf("Filosofo %d pegou o garfo esquerdo \n", f.id)

		if !f.rightFork.TryLock() {
			f.leftFork.Unlock()
			continue
		}

		fmt.Printf("Filosofo %d pegou o garfo direito\n", f.id)

		f.comer()

		f.rightFork.Unlock()
		fmt.Printf("Filosofo %d devolveu o garfo direito\n", f.id)

		f.leftFork.Unlock()
		fmt.Printf("Filosofo %d devolveu o garfo esquerdo\n", f.id)
	}
}

func main() {
	var wg sync.WaitGroup
	garfos := make([]*sync.Mutex, Num_filosofos)
	filosofos := make([]*Filosofo, Num_filosofos)

	for i := range garfos {
		garfos[i] = &sync.Mutex{}
	}
	for i := 0; i < Num_filosofos; i++ {
		filosofos[i] = &Filosofo{
			id:        i + 1,
			leftFork:  garfos[i],
			rightFork: garfos[(i+1)%Num_filosofos],
		}
	}

	for i := 0; i < Num_filosofos; i++ {
		wg.Add(1)
		go func(filosofo *Filosofo) {
			defer wg.Done()
			filosofo.tentarComer()
		}(filosofos[i])
	}
	wg.Wait()
}
