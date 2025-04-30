package trigger

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type TriggerManager struct {
	mu             sync.Mutex
	requestCount   int
	lastTrigger    time.Time
	triggerPending bool
	triggerRunning bool // impede múltiplas triggers simultâneas
	minuteLimit    time.Duration
	threshold      int
	triggerURL     string
	client         *http.Client
}

func NewTriggerManager(threshold int, minuteLimit time.Duration, triggerURL string) *TriggerManager {
	tm := &TriggerManager{
		threshold:   threshold,
		minuteLimit: minuteLimit,
		triggerURL:  triggerURL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
	go tm.monitorTrigger()
	return tm
}

func (tm *TriggerManager) RegisterRequest() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.requestCount++
	log.Printf("Requisição registrada. Total: %d", tm.requestCount)

	if tm.requestCount >= tm.threshold && !tm.triggerRunning {
		if time.Since(tm.lastTrigger) >= tm.minuteLimit {
			tm.lastTrigger = time.Now()
			tm.requestCount = 0
			tm.triggerRunning = true
			go tm.executeTriggerWithRetry()
		} else {
			tm.triggerPending = true
		}
	}
}

func (tm *TriggerManager) monitorTrigger() {
	for {
		time.Sleep(5 * time.Second)
		tm.mu.Lock()
		if tm.triggerPending && !tm.triggerRunning && time.Since(tm.lastTrigger) >= tm.minuteLimit {
			tm.lastTrigger = time.Now()
			tm.requestCount = 0
			tm.triggerPending = false
			tm.triggerRunning = true
			go tm.executeTriggerWithRetry()
		}
		tm.mu.Unlock()
	}
}

func (tm *TriggerManager) executeTriggerWithRetry() {
	defer func() {
		tm.mu.Lock()
		tm.triggerRunning = false
		tm.mu.Unlock()
	}()

	for {
		log.Printf("🔁 Tentando executar trigger: %s", tm.triggerURL)
		resp, err := tm.client.Get(tm.triggerURL)
		if err != nil {
			log.Printf("❌ Erro de conexão: %v", err)
		} else {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				log.Printf("✅ Trigger executada com sucesso! Status: %d", resp.StatusCode)
				return
			}
			log.Printf("⚠️ Trigger respondeu com status: %d (esperado 200)", resp.StatusCode)
		}

		// Espera 10s antes de tentar novamente
		time.Sleep(10 * time.Second)
	}
}
