package main

import "fmt"

// -----------------------------------------------------------------
// PASSO 1: Definir o Contrato (A "Tomada de Parede")
// -----------------------------------------------------------------
// A interface Notificador define um contrato. Qualquer tipo que
// tiver um método Enviar(mensagem string) será considerado um Notificador.
// Ela não se importa COMO a notificação é enviada, apenas QUE ELA POSSA SER ENVIADA.
// type Notificador interface {
// 	Enviar(mensagem string) error
// }

// -----------------------------------------------------------------
// PASSO 2: Criar as Implementações Concretas (Os "Aparelhos")
// -----------------------------------------------------------------
// -----------------------------------------------------------------
// A interface Notificador define um contrato. Qualquer tipo que
// tiver um método Enviar(mensagem string) será considerado um Notificador.
// Ela não se importa COMO a notificação é enviada, apenas QUE ELA POSSA SER ENVIADA.
type Notificador interface {
	Enviar(mensagem string) error
}

// -----------------------------------------------------------------
// PASSO 2: Criar as Implementações Concretas (Os "Aparelhos")
// -----------------------------------------------------------------

// Implementação #1: EmailNotificador
// Este é um tipo concreto que sabe como enviar um e-mail.
type EmailNotificador struct {
	Destinatario string
}

// Ao implementar este método, EmailNotificador automaticamente
// satisfaz a interface Notificador. Ele cumpre o contrato.
func (e EmailNotificador) Enviar(mensagem string) error {
	fmt.Printf("Enviando E-MAIL para '%s': %s\n", e.Destinatario, mensagem)
	return nil // Em um caso real, poderia haver um erro aqui.
}

// Implementação #2: SMSNotificador
// Este é outro tipo concreto, que sabe como enviar um SMS.
type SMSNotificador struct {
	NumeroTelefone string
}

// SMSNotificador também cumpre o contrato da interface Notificador,
// pois também possui o método Enviar(mensagem string).
func (s *SMSNotificador) Enviar(mensagem string) error {
	fmt.Printf("Enviando SMS para '%s': %s\n", s.NumeroTelefone, mensagem)
	return nil
}

// -----------------------------------------------------------------
// PASSO 3: Criar uma Função que Usa a Interface (O "Consumidor")
// -----------------------------------------------------------------
// Esta função é o nosso "quarto". Ela não sabe se vai receber um E-mail
// ou um SMS. Ela só sabe que vai receber um "Notificador".
// Ela trabalha apenas com a interface (o contrato), não com os tipos concretos.
func EnviarAlerta(n Notificador, alerta string) {
	fmt.Println("--- Disparando um novo alerta ---")
	err := n.Enviar(alerta)
	if err != nil {
		fmt.Println("Ocorreu um erro ao enviar o alerta:", err)
	}
	fmt.Println("-------------------------------")
}

// -----------------------------------------------------------------
// PASSO 4: Juntar Tudo (O "Eletricista")
// -----------------------------------------------------------------
func main() {
	// Criamos nossos "aparelhos" concretos
	notificacaoPorEmail := EmailNotificador{Destinatario: "chefe@empresa.com"}
	notificacaoPorSMS := &SMSNotificador{NumeroTelefone: "+55 67 99999-8888"}

	// Agora, usamos a mesma função EnviarAlerta para os dois tipos diferentes.
	// A função não se importa com o tipo, desde que ele cumpra o contrato Notificador.

	// Aqui, estamos "plugando" o aparelho de Email na função.
	EnviarAlerta(notificacaoPorEmail, "O servidor principal está offline!")

	// E aqui, estamos "plugando" o aparelho de SMS na mesma função.
	EnviarAlerta(notificacaoPorSMS, "Código de verificação: 45678")
}
