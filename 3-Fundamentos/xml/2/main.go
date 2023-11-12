package main

import (
	"encoding/xml"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
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
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer xmlFile.Close()

	var dhcp DHCPD
	if err := xml.NewDecoder(xmlFile).Decode(&dhcp); err != nil {
		fmt.Println("Erro ao decodificar o XML:", err)
		return
	}

	for _, vlan := range dhcp.VLANs {
		ip := net.ParseIP(vlan.Range.From)
		if ip == nil {
			fmt.Println("Endereço IP inválido")
			return
		}
		vlan.Network = ip.Mask(net.IPMask(net.ParseIP("255.255.254.0"))).String()
		fmt.Printf("Número da VLAN: %v\n", vlan.Vlan)
		fmt.Printf("Gateway: %s\n", vlan.Network)
		fmt.Printf("Range: %s - %s\n", vlan.Range.From, vlan.Range.To)
		fmt.Printf("DHCPLeaseInLocalTime: %s\n", vlan.DHCPLeaseInLocalTime)
		// Processar os StaticMaps e outros campos conforme necessário
		for _, staticMap := range vlan.StaticMap {
			replaceColonWithHyphen(&staticMap.MAC)
			fmt.Printf("MAC: %s, CID: %s, IPAddr: %s, Hostname: %s, Descr: %s\n", staticMap.MAC, staticMap.CID, staticMap.IPAddr, staticMap.Hostname, staticMap.Descr)
			// Processar outros campos estáticos, se necessário
			ExecutePowerShell(&vlan.Network, &staticMap.IPAddr, &staticMap.MAC, &staticMap.Hostname)
		}
		fmt.Println("-------------------------------------------------")
	}

}

func replaceColonWithHyphen(macAddress *string) {
	*macAddress = strings.ReplaceAll(*macAddress, ":", "-")
}

func ExecutePowerShell(scopeId, ipaddress, macaddress, hostname *string) {

	command := fmt.Sprintf("Add-DhcpServerv4Reservation -ScopeId %s -IPAddress %s -ClientId %s -Name %s", *scopeId, *ipaddress, *macaddress, *hostname)
	cmd := exec.Command("powershell", "-Command", command)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Command finished with error: %v", err)
	}
}
