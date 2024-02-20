package tests

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type LoadTestReport struct {
	TotalTime      time.Duration
	TotalRequests  int
	SuccessCount   int
	ErrorsCount    int
	StatusCodeDist map[int]int
	mu             sync.Mutex // Mutex para sincronizar acesso à estrutura
}

func RunLoadTest(url string, totalRequests int, concurrency int) LoadTestReport {
	if totalRequests < concurrency {
		concurrency = totalRequests
	}
	var report LoadTestReport
	report.StatusCodeDist = make(map[int]int)

	wg := sync.WaitGroup{}
	requestsChan := make(chan struct{}, concurrency)

	startTime := time.Now()

	for i := 0; i < totalRequests; i++ {
		requestsChan <- struct{}{}
		wg.Add(1)

		go func() {
			defer wg.Done()
			resp, err := MakeRequest(url)
			if err != nil {
				fmt.Println(err.Error())
				report.ErrorsCount++
			} else {
				report.mu.Lock()
				report.StatusCodeDist[resp.StatusCode]++
				if resp.StatusCode == http.StatusOK {
					report.SuccessCount++
				}
				report.mu.Unlock() // Desbloqueia após a modificação
			}
			<-requestsChan
		}()
	}

	wg.Wait() // Espera todas as requisições serem completadas
	report.TotalRequests = totalRequests
	report.TotalTime = time.Since(startTime) // Calcula o tempo total após todas as requisições
	return LoadTestReport{
		TotalTime:      report.TotalTime,
		TotalRequests:  report.TotalRequests,
		SuccessCount:   report.SuccessCount,
		ErrorsCount:    report.ErrorsCount,
		StatusCodeDist: report.StatusCodeDist,
	}
}
