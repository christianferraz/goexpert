package trigger

import (
	"crypto/tls"
	"log"
	"net/http"
	"sync"
	"time"
)

type TriggerManager struct {
	mu           sync.Mutex
	requestCount int
	triggerURL   string
	client       *http.Client
	interval     time.Duration
}

func NewTriggerManager(interval time.Duration, triggerURL string) *TriggerManager {
	tm := &TriggerManager{
		interval:   interval,
		triggerURL: triggerURL,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			Timeout: 5 * time.Second,
		},
	}

	go tm.startTriggerLoop()

	return tm
}

// Método chamado externamente quando uma requisição for recebida
func (tm *TriggerManager) RegisterRequest() {
	tm.mu.Lock()
	tm.requestCount++
	tm.mu.Unlock()
	log.Printf("📥 Requisição recebida. Total acumulado: %d", tm.requestCount)
}

// Loop que dispara a trigger a cada X segundos, se houver requisições
func (tm *TriggerManager) startTriggerLoop() {
	ticker := time.NewTicker(tm.interval)
	defer ticker.Stop()

	for {
		<-ticker.C

		tm.mu.Lock()
		if tm.requestCount > 0 {
			log.Printf("🚨 %d requisições detectadas no intervalo. Enviando trigger...", tm.requestCount)
			go tm.executeTriggerWithRetry()
			tm.requestCount = 0
		} else {
			log.Println("⏱️ Nenhuma requisição registrada. Trigger não será enviada.")
		}
		tm.mu.Unlock()
	}
}

// Envia a trigger com tentativas de retry até obter sucesso
func (tm *TriggerManager) executeTriggerWithRetry() {
	for {
		log.Printf("🔁 Tentando executar trigger: %s", tm.triggerURL)
		resp, err := tm.client.Get(tm.triggerURL)
		if err != nil {
			log.Printf("❌ Erro ao conectar: %v", err)
		} else {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				log.Printf("✅ Trigger executada com sucesso! Status: %d", resp.StatusCode)
				return
			}
			log.Printf("⚠️ Trigger respondeu com status: %d", resp.StatusCode)
		}
		time.Sleep(10 * time.Second) // Espera antes de tentar novamente
	}
}
