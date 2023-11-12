package main

import (
	"encoding/xml"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type Range struct {
	From string `xml:"from"`
	To   string `xml:"to"`
}

type StaticMap struct {
	MAC      string `xml:"mac"`
	CID      string `xml:"cid"`
	IPAddr   string `xml:"ipaddr"`
	Hostname string `xml:"hostname"`
	Descr    string `xml:"descr"`
	// Adicione os campos adicionais conforme necessário
}

type VLAN struct {
	Network              string `xml:"-"`
	NrVlan               string `xml:"-"`
	Range                Range  `xml:"range"`
	Vlan                 int    `xml:"vlan"`
	Gateway              string `xml:"gateway"`
	DHCPLeaseInLocalTime string `xml:"dhcpleaseinlocaltime"`
	// Adicione outros campos conforme necessário
	StaticMap []StaticMap `xml:"staticmap"`
	// Adicione outros campos conforme necessário
}

type DHCPD struct {
	XMLName xml.Name `xml:"dhcpd"`
	VLANs   []VLAN   `xml:",any"`
}

func main() {
	xmlFile, err := os.Open("dhcp-pfsense.xml")
	var wg sync.WaitGroup
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)

	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			if strings.HasPrefix(se.Name.Local, "vlan") {
				var vlan VLAN
				vlan.NrVlan = strings.TrimPrefix(se.Name.Local, "vlan")
				if err := decoder.DecodeElement(&vlan, &se); err != nil {
					fmt.Println("Erro ao decodificar VLAN:", err)
				}
				wg.Add(1)
				go processVLAN(&wg, &vlan)
			}
		}
	}
	wg.Wait()
}

func processVLAN(wg *sync.WaitGroup, vlan *VLAN) {
	defer wg.Done()
	ip := net.ParseIP(vlan.Range.From)
	if ip == nil {
		fmt.Println("Endereço IP inválido")
		return
	}
	vlan.Network = ip.Mask(net.IPMask(net.ParseIP("255.255.254.0"))).String()

	// fmt.Printf("Número da VLAN: %v\n", vlan.NrVlan)
	// fmt.Printf("Gateway: %s\n", vlan.Network)
	// fmt.Printf("Range: %s - %s\n", vlan.Range.From, vlan.Range.To)
	// fmt.Printf("DHCPLeaseInLocalTime: %s\n", vlan.DHCPLeaseInLocalTime)

	for _, staticMap := range vlan.StaticMap {
		replaceColonWithHyphen(&staticMap.MAC)
		// fmt.Printf("MAC: %s, CID: %s, IPAddr: %s, Hostname: %s, Descr: %s\n", staticMap.MAC, staticMap.CID, staticMap.IPAddr, staticMap.Hostname, staticMap.Descr)
		// Processar outros campos estáticos, se necessário
		ExecutePowerShell(&vlan.Network, &staticMap.IPAddr, &staticMap.MAC, &staticMap.Hostname)
	}
}

func replaceColonWithHyphen(macAddress *string) {
	*macAddress = strings.ReplaceAll(*macAddress, ":", "-")
}

func ExecutePowerShell(scopeId, ipaddress, macaddress, hostname *string) {
	command := fmt.Sprintf("Add-DhcpServerv4Reservation -ScopeId %s -IPAddress %s -ClientId %s -Name %s", *scopeId, *ipaddress, *macaddress, *hostname)
	fmt.Println(command)
	cmd := exec.Command("powershell", "-Command", command)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Command finished with error: %v", err)
	}
}
