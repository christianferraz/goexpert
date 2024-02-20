package tests

import "fmt"

func (r *LoadTestReport) Print() {
	fmt.Println("Load Test Report")
	fmt.Println("----------------")
	fmt.Println("Tempo total gasto na execução - ", r.TotalTime)
	fmt.Println("Quantidades de requests com sucesso - ", r.SuccessCount)
	for i := range r.StatusCodeDist {
		fmt.Println("Código de retorno ", i, " - ", r.StatusCodeDist[i])
	}
	fmt.Println("Quantidade de erros - ", r.ErrorsCount)
}
