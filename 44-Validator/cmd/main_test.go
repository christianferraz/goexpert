package main

import (
	"context"
	"testing"
	"time"

	"github.com/christianferraz/goexpert/44-Validator/internal/validator"
)

// 1. TESTE BÁSICO - Usuário válido
func TestUsuarioValido(t *testing.T) {
	user := Usuario{
		Nome:  "João Silva",
		Email: "joao@email.com",
		CPF:   "82084610168", // CPF válido
		Site:  "https://example.com",
		Idade: 25,
	}

	ctx := context.Background()
	erros := validator.ValidateStruct(ctx, user)

	// Verifica se NÃO há erros
	if len(erros) != 0 {
		t.Errorf("Esperado 0 erros, mas obteve %d erros:", len(erros))
		for key, msg := range erros {
			t.Errorf("  Campo '%s': %s", key, msg)
		}
	}
}

// 2. TESTE BÁSICO - Usuário inválido
func TestUsuarioInvalido(t *testing.T) {
	user := Usuario{
		Nome:  "C",             // muito curto
		Email: "invalid-email", // formato inválido
		CPF:   "123",           // CPF inválido
		Site:  "invalid-url",   // URL inválida
		Idade: 15,              // menor que 18
	}

	ctx := context.Background()
	erros := validator.ValidateStruct(ctx, user)

	// Verifica se HÁ erros (deve ter 5 erros)
	expectedErrors := 5
	if len(erros) != expectedErrors {
		t.Errorf("Esperado %d erros, mas obteve %d erros", expectedErrors, len(erros))
	}

	// Verifica erros específicos
	if _, exists := erros["Nome"]; !exists {
		t.Error("Esperado erro no campo 'Nome'")
	}

	if _, exists := erros["Email"]; !exists {
		t.Error("Esperado erro no campo 'Email'")
	}

	if _, exists := erros["CPF"]; !exists {
		t.Error("Esperado erro no campo 'CPF'")
	}

	if _, exists := erros["Site"]; !exists {
		t.Error("Esperado erro no campo 'Site'")
	}

	if _, exists := erros["Idade"]; !exists {
		t.Error("Esperado erro no campo 'Idade'")
	}
}

// 3. TESTE COM CONTEXT CANCELADO
func TestContextCancelado(t *testing.T) {
	user := Usuario{
		Nome:  "João Silva",
		Email: "joao@email.com",
		CPF:   "82084610168",
		Site:  "https://example.com",
		Idade: 25,
	}

	// Cria um context já cancelado
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancela imediatamente

	// A validação ainda deve funcionar (context não é usado internamente)
	erros := validator.ValidateStruct(ctx, user)

	if len(erros) != 0 {
		t.Errorf("Validação deveria ter funcionado mesmo com context cancelado. Erros: %v", erros)
	}
}

// 4. TESTE COM TIMEOUT
func TestContextComTimeout(t *testing.T) {
	user := Usuario{
		Nome:  "João Silva",
		Email: "joao@email.com",
		CPF:   "820.846.101-68",
		Site:  "https://example.com",
		Idade: 125,
	}

	// Context com timeout muito curto
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Aguarda o timeout
	time.Sleep(1 * time.Millisecond)

	// Validação ainda funciona (context não usado)
	erros := validator.ValidateStruct(ctx, user)

	if len(erros) != 0 {
		t.Errorf("Validação deveria ter funcionado mesmo com timeout. Erros: %v", erros)
	}
}

// 5. TESTE COM VALUES NO CONTEXT
func TestContextComValues(t *testing.T) {
	user := Usuario{
		Nome:  "Joao Silva",
		Email: "joao@email.com",
		CPF:   "82084610168",
		Site:  "https://example.com",
		Idade: 25,
	}

	// Context com valores customizados
	ctx := context.WithValue(context.Background(), "requestID", "12345")
	ctx = context.WithValue(ctx, "userID", "67890")

	erros := validator.ValidateStruct(ctx, user)

	if len(erros) != 0 {
		t.Errorf("Validação falhou inesperadamente. Erros: %v", erros)
	}
}

// 6. TESTE DE PERFORMANCE COM CONTEXT
func TestPerformanceComContext(t *testing.T) {
	user := Usuario{
		Nome:  "Joao Silva",
		Email: "joao@email.com",
		CPF:   "82084610168",
		Site:  "https://example.com",
		Idade: 25,
	}

	ctx := context.Background()

	// Mede o tempo de execução
	start := time.Now()

	// Executa validação múltiplas vezes
	for i := 0; i < 1000; i++ {
		erros := validator.ValidateStruct(ctx, user)
		if len(erros) != 0 {
			t.Errorf("Validação falhou na iteração %d: %v", i, erros)
		}
	}

	duration := time.Since(start)
	t.Logf("1000 validações executadas em %v", duration)

	// Verifica se não demorou muito (ajuste conforme necessário)
	if duration > 100*time.Millisecond {
		t.Errorf("Validação muito lenta: %v", duration)
	}
}

// 7. TESTE DE CAMPOS OPCIONAIS (omitempty)
func TestCamposOpcionais(t *testing.T) {
	// Usuário sem site (opcional)
	user := Usuario{
		Nome:  "João Silva",
		Email: "joao@email.com",
		CPF:   "82084610168",
		Site:  "", // vazio, mas é omitempty
		Idade: 25,
	}

	ctx := context.Background()
	erros := validator.ValidateStruct(ctx, user)

	// Não deve ter erros
	if len(erros) != 0 {
		t.Errorf("Esperado 0 erros para campo opcional vazio, mas obteve: %v", erros)
	}

	// Usuário com site inválido (não deve passar)
	userComSiteInvalido := Usuario{
		Nome:  "João Silva",
		Email: "joao@email.com",
		CPF:   "82084610168",
		Site:  "site-invalido", // inválido
		Idade: 25,
	}

	erros = validator.ValidateStruct(ctx, userComSiteInvalido)

	// Deve ter erro no site
	if _, exists := erros["Site"]; !exists {
		t.Error("Esperado erro no campo 'Site' quando URL é inválida")
	}
}

// 8. TESTE COM DIFERENTES TIPOS DE CONTEXT
func TestDiferentesTiposDeContext(t *testing.T) {
	user := Usuario{
		Nome:  "João Silva",
		Email: "joao@email.com",
		CPF:   "82084610168",
		Site:  "https://example.com",
		Idade: 25,
	}

	tests := []struct {
		name string
		ctx  context.Context
	}{
		{
			name: "Background Context",
			ctx:  context.Background(),
		},
		{
			name: "TODO Context",
			ctx:  context.TODO(),
		},
		{
			name: "Context com Deadline",
			ctx: func() context.Context {
				ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(1*time.Hour))
				return ctx
			}(),
		},
		{
			name: "Context com Cancel",
			ctx: func() context.Context {
				ctx, _ := context.WithCancel(context.Background())
				return ctx
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			erros := validator.ValidateStruct(tt.ctx, user)

			if len(erros) != 0 {
				t.Errorf("Validação falhou com %s: %v", tt.name, erros)
			}
		})
	}
}

// 9. BENCHMARK PARA MEDIR PERFORMANCE
func BenchmarkValidateStruct(b *testing.B) {
	user := Usuario{
		Nome:  "João Silva",
		Email: "joao@email.com",
		CPF:   "82084610168",
		Site:  "https://example.com",
		Idade: 25,
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.ValidateStruct(ctx, user)
	}
}

// 10. EXEMPLO DE TESTE TABLE-DRIVEN
func TestValidacaoTableDriven(t *testing.T) {
	tests := []struct {
		name           string
		usuario        Usuario
		expectedErrors int
		expectedFields []string
	}{
		{
			name: "Usuario Completamente Válido",
			usuario: Usuario{
				Nome:  "João Silva",
				Email: "joao@email.com",
				CPF:   "82084610168",
				Site:  "https://example.com",
				Idade: 25,
			},
			expectedErrors: 0,
			expectedFields: []string{},
		},
		{
			name: "Nome Muito Curto",
			usuario: Usuario{
				Nome:  "J",
				Email: "joao@email.com",
				CPF:   "82084610168",
				Site:  "https://example.com",
				Idade: 25,
			},
			expectedErrors: 1,
			expectedFields: []string{"Nome"},
		},
		{
			name: "Múltiplos Erros",
			usuario: Usuario{
				Nome:  "J",
				Email: "invalid",
				CPF:   "123",
				Site:  "invalid",
				Idade: 15,
			},
			expectedErrors: 5,
			expectedFields: []string{"Nome", "Email", "CPF", "Site", "Idade"},
		},
		{
			name: "Site Opcional Vazio",
			usuario: Usuario{
				Nome:  "João Silva",
				Email: "joao@email.com",
				CPF:   "82084610168",
				Site:  "", // vazio mas opcional
				Idade: 25,
			},
			expectedErrors: 0,
			expectedFields: []string{},
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			erros := validator.ValidateStruct(ctx, tt.usuario)

			// Verifica número de erros
			if len(erros) != tt.expectedErrors {
				t.Errorf("Esperado %d erros, obteve %d: %v", tt.expectedErrors, len(erros), erros)
			}

			// Verifica campos específicos com erro
			for _, field := range tt.expectedFields {
				if _, exists := erros[field]; !exists {
					t.Errorf("Esperado erro no campo '%s'", field)
				}
			}
		})
	}
}
