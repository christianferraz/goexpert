package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func gerarCertificadoPDF(nome string) ([]byte, error) {
	// Lê o template SVG
	svgTemplate, err := os.ReadFile("certificado_template.svg")
	if err != nil {
		return nil, err
	}

	// Substitui placeholders
	svgStr := string(svgTemplate)
	svgStr = strings.Replace(svgStr, "NOME_AQUI", nome, 1)
	svgStr = strings.Replace(svgStr, "TIPO_AQUI", "Participante", 1)

	// Converte SVG para PDF usando rsvg-convert
	cmd := exec.Command("rsvg-convert", "-f", "pdf", "-")
	cmd.Stdin = bytes.NewBufferString(svgStr)

	var out bytes.Buffer
	cmd.Stdout = &out

	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func main() {
	pdfData, err := gerarCertificadoPDF("Maria Oliveira")
	if err != nil {
		log.Fatalf("Erro ao gerar PDF: %v", err)
	}

	err = os.WriteFile("certificado.pdf", pdfData, 0644)
	if err != nil {
		log.Fatalf("Erro ao salvar PDF: %v", err)
	}

	fmt.Println("Certificado PDF gerado com sucesso!")
}
