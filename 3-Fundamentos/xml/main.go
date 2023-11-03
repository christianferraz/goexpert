package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"strings"
)

type StaticMap struct {
	MAC      string `xml:"mac"`
	IPAddr   string `xml:"ipaddr"`
	Hostname string `xml:"hostname"`
}

type DHCPD struct {
	Options []struct {
		StaticMaps []StaticMap `xml:"staticmap"`
	} `xml:",any"`
}

type DHCPLease struct {
	HostName   string `xml:"HostName"`
	MACAddress string `xml:"MACAddress"`
	IPAddress  string `xml:"IPAddress"`
}

type DHCPServer struct {
	StaticLeases []DHCPLease `xml:"StaticLease>Lease"`
}

type Response struct {
	XMLName xml.Name     `xml:"Response"`
	Servers []DHCPServer `xml:"DHCPServer"`
}

func main() {

	xmlData1, err := os.ReadFile("lista1.xml")
	if err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
		return
	}
	outputFile, err := os.Create("resultado.txt")
	if err != nil {
		log.Fatal("Erro ao criar o arquivo de saída:", err)
	}
	defer outputFile.Close()

	// Redirecionar a saída padrão para o arquivo
	log.SetOutput(outputFile)

	var dhcpd DHCPD
	err = xml.Unmarshal(xmlData1, &dhcpd)
	if err != nil {
		fmt.Println("Erro ao decodificar XML:", err)
		return
	}

	var dataSlice []struct {
		MAC      string
		IPAddr   string
		Hostname string
	}

	for _, opt := range dhcpd.Options {
		for _, staticMap := range opt.StaticMaps {
			dataSlice = append(dataSlice, struct {
				MAC      string
				IPAddr   string
				Hostname string
			}{
				MAC:      staticMap.MAC,
				IPAddr:   staticMap.IPAddr,
				Hostname: staticMap.Hostname,
			})
		}
	}

	xmlData2, err := os.ReadFile("lista2.xml")
	if err != nil {
		log.Fatal(err)
	}

	var response Response
	err = xml.Unmarshal(xmlData2, &response)
	if err != nil {
		log.Fatal(err)
	}

	var staticLeases []DHCPLease
	for _, server := range response.Servers {
		staticLeases = append(staticLeases, server.StaticLeases...)
	}

	fmt.Println("Static Leases:")
	// for _, lease := range staticLeases {
	// 	fmt.Printf("HostName: %s\n", lease.HostName)
	// 	fmt.Printf("MACAddress: %s\n", lease.MACAddress)
	// 	fmt.Printf("IPAddress: %s\n", lease.IPAddress)
	// }
	// for _, data := range dataSlice {
	// 	fmt.Printf("MAC: %s, IPAddr: %s, Hostname: %s\n", data.MAC, data.IPAddr, data.Hostname)
	// }
	for _, data := range dataSlice {
		found := false

		for _, lease := range staticLeases {
			if strings.EqualFold(lease.MACAddress, data.MAC) &&
				strings.EqualFold(lease.IPAddress, data.IPAddr) {
				found = true
				break
			}
		}

		if found {
			// log.Printf("O elemento MAC: %s, IPAddr: %s, Hostname: %s está contido em staticLeases.\n", data.MACAddress, data.IPAddress, data.HostName)
		} else {
			log.Printf("O elemento MAC: %s, IPAddr: %s, Hostname: %s NÃO está contido em staticLeases.\n", data.MAC, data.IPAddr, data.Hostname)
		}
	}
}
