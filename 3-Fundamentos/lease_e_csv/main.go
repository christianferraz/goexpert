package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"

	"golang.org/x/crypto/ssh"
)

func main() {
	conect_ssh()
	// Abra o arquivo "dhcpd.leases"
	inputFile, err := os.Open("saida_dhcpd.leases")
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo dhcpd.leases:", err)
		return
	}
	defer inputFile.Close()

	// Crie um arquivo CSV de saída
	outputFile, err := os.Create("output.csv")
	if err != nil {
		fmt.Println("Erro ao criar o arquivo de saída:", err)
		return
	}
	defer outputFile.Close()

	reader := bufio.NewReader(inputFile)
	writer := csv.NewWriter(outputFile)
	// Cabeçalhos das colunas
	headers := []string{"Lease", "Hardware Ethernet", "Client Hostname"}

	// Escreva os cabeçalhos no arquivo CSV
	err = writer.Write(headers)
	if err != nil {
		fmt.Println("Erro ao escrever cabeçalhos no arquivo CSV:", err)
		return
	}
	// Expressões regulares para extrair os valores desejados
	leaseRegex := regexp.MustCompile(`lease (\S+) {`)
	hardwareEthernetRegex := regexp.MustCompile(`hardware ethernet (\S+);`)
	clientHostnameRegex := regexp.MustCompile(`client-hostname "([^"]+)"`)

	// Mapa para rastrear valores únicos
	uniqueValues := make(map[string]struct{})

	// Variáveis para armazenar os valores
	lease := ""
	hardwareEthernet := ""
	clientHostname := ""

	for {
		// Leia uma linha do arquivo
		line, err := reader.ReadString('\n')
		if err != nil {
			break // Sai do loop no final do arquivo
		}

		// Encontre correspondências para cada expressão regular no bloco de texto
		leaseMatches := leaseRegex.FindStringSubmatch(line)
		hardwareEthernetMatches := hardwareEthernetRegex.FindStringSubmatch(line)
		clientHostnameMatches := clientHostnameRegex.FindStringSubmatch(line)

		// Extraia os valores correspondentes
		if len(leaseMatches) > 1 {
			lease = leaseMatches[1]
		}
		if len(hardwareEthernetMatches) > 1 {
			hardwareEthernet = hardwareEthernetMatches[1]
		}
		if len(clientHostnameMatches) > 1 {
			clientHostname = clientHostnameMatches[1]
		}

		// Se encontramos um bloco completo, escreva os valores no arquivo CSV e adicione o "lease" ao mapa de valores únicos
		if line == "}\n" {
			value := lease + "," + hardwareEthernet + "," + clientHostname
			_, exists := uniqueValues[value]
			if !exists {
				uniqueValues[value] = struct{}{}
				err := writer.Write([]string{lease, hardwareEthernet, clientHostname})
				if err != nil {
					fmt.Println("Erro ao escrever no arquivo CSV:", err)
				}
			}

			// Ressete as variáveis
			lease = ""
			hardwareEthernet = ""
			clientHostname = ""
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		fmt.Println("Erro ao finalizar o arquivo CSV:", err)
	}

	fmt.Println("Valores repetidos removidos e arquivo de saída criado.")
}

func conect_ssh() {
	// Configure as credenciais SSH
	sshConfig := &ssh.ClientConfig{
		User: "admin",
		Auth: []ssh.AuthMethod{
			ssh.Password("Miuchinh@050704"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Endereço e porta SSH do host
	sshAddress := "10.42.3.254:22"

	// Crie uma conexão SSH
	sshClient, err := ssh.Dial("tcp", sshAddress, sshConfig)
	if err != nil {
		log.Fatalf("Erro ao conectar via SSH: %v", err)
	}
	defer sshClient.Close()

	// Crie uma sessão SSH
	sshSession, err := sshClient.NewSession()
	if err != nil {
		log.Fatalf("Erro ao criar sessão SSH: %v", err)
	}
	defer sshSession.Close()
	// Redirecione a entrada e a saída para a sessão SSH
	stdin, err := sshSession.StdinPipe()
	if err != nil {
		log.Fatalf("Erro ao obter a entrada padrão da sessão SSH: %v", err)
	}
	stdout, err := sshSession.StdoutPipe()
	if err != nil {
		log.Fatalf("Erro ao obter a saída padrão da sessão SSH: %v", err)
	}

	err = sshSession.Start("/bin/sh")
	if err != nil {
		log.Fatalf("Erro ao iniciar a sessão SSH: %v", err)
	}
	// Envie os comandos 5 e 3
	commands := []string{"5", "3", "cat dhcpd.leases", "exit", "0", "0"}
	for _, command := range commands {
		_, err := io.WriteString(stdin, command+"\n")
		if err != nil {
			log.Fatalf("Erro ao enviar o comando: %v", err)
		}
	}
	// sshSession.Wait()

	// Salve o conteúdo do arquivo remoto em um arquivo local
	outputFile, err := os.Create("saida_dhcpd.leases")
	if err != nil {
		log.Fatalf("Erro ao criar o arquivo de saída: %v", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, stdout)
	err = sshSession.Wait()
	fmt.Println("passou do bin sh")
	if err != nil {
		log.Fatalf("Erro ao salvar a saída no arquivo: %v", err)
	}
	if err != nil {
		log.Fatalf("Erro ao aguardar o término da sessão SSH: %v", err)
	}
	fmt.Println("Conteúdo do arquivo /tmp/dhcpd.leases salvo localmente como dhcpd.leases")
}

func saveOutputToFile(filename string, output []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, bytes.NewReader(output))
	if err != nil {
		return err
	}

	return nil
}
