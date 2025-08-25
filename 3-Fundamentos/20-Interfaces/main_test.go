package main

import "testing"

func TestMain(t *testing.T) {
	emailNotificador := EmailNotificador{Destinatario: "eu"}
	EnviarAlerta(emailNotificador, "O servidor está com 90% de uso de CPU")
	smsNotificador := &SMSNotificador{NumeroTelefone: "+55 67 99999-8888"}
	EnviarAlerta(smsNotificador, "Código de verificação: 12345")
}
