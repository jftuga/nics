/*
nics.go
-John Taylor
2019-08-03

Display information about Network Inferface Cards (NICs)

To compile:
go build -ldflags="-s -w"

MIT License; Copyright (c) 2019 John Taylor
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

*/

package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const version = "2.0.0"

type DHCPEntry struct {
	DhcpServer   string
	LeaseRenewed time.Time
	LeaseExpires time.Time
}

type GatewayEntry struct {
	iFaceIP   string
	ipAddress string
}

type Nic struct {
	iFace        net.Interface
	Addrs        []net.Addr
	IPv4         []string
	IPv6         []string
	DNS          []string
	DHCP         *DHCPEntry
	Gateways     []GatewayEntry
	isBriefEntry bool
}

func (n *Nic) Print() {
	fmt.Println("Name        :", n.iFace.Name)
	fmt.Println("Addrs       :", n.Addrs)
	fmt.Println("mtu         :", n.GetMTU())
	fmt.Println("flags       :", n.GetFlags())
	fmt.Println("IPv4 address:", n.IPv4)
	fmt.Println("IPv6 address:", n.IPv6)

}

func (n *Nic) GetMTU() string {
	return strconv.Itoa(n.iFace.MTU)
}

func (n *Nic) GetFlags() []string {
	return strings.Split(n.iFace.Flags.String(), "|")
}

// SetIPaddrs sets and returns the number of IPv4, IPv6 addresses
func (n *Nic) SetIPAddrs() (int, int) {
	addr, _ := n.iFace.Addrs()
	var ipPtr *([]string)
	for _, ip := range addr {
		ipAddress := ip.String()
		if strings.Contains(ipAddress, ":") {
			ipPtr = &n.IPv6
		} else {
			ipPtr = &n.IPv4
		}
		*ipPtr = append(*ipPtr, ipAddress)
	}

	return len(n.IPv4), len(n.IPv6)
}

func getGWForIP(n *Nic, IPv4or6 []string) int {
	ipMapGateway := make(map[string]string)

	for _, ipWithMask := range IPv4or6 {
		ip, _ := splitIPMask(ipWithMask)
		gateway, _ := GetGatewayForIP(ip)
		if len(gateway) > 0 {
			ipMapGateway[ip] = gateway
		}
	}
	if len(ipMapGateway) == 0 {
		return 0
	}

	if ipMapGateway != nil {
		for ip, gw := range ipMapGateway {
			fmt.Printf("i: %s  g: %s\n", ip, gw)
			var ge GatewayEntry
			ge.iFaceIP = ip
			ge.ipAddress = gw
			n.Gateways = append(n.Gateways, ge)
		}
	}
	return len(ipMapGateway)
}

func queryNetworkInterfaces() []Nic {
	adapters, err := net.Interfaces()
	if err != nil {
		fmt.Println(fmt.Errorf("%+v", err.Error()))
		return nil
	}

	var allNetworkInterfaces []Nic
	for _, iface := range adapters {
		var n Nic
		n.iFace = iface
		n.Addrs, _ = iface.Addrs()

		n.SetIPAddrs()
		allNetworkInterfaces = append(allNetworkInterfaces, n)
	}

	return allNetworkInterfaces
}

func main() {
	argsAllDetails := flag.Bool("a", false, "show all details on ALL interfaces, including DHCP")
	argsDebug := flag.Bool("d", false, "show debug information")
	argsVersion := flag.Bool("v", false, "show program version")
	flag.Usage = func() {
		pgmName := os.Args[0]
		if strings.HasPrefix(os.Args[0], "./") {
			pgmName = os.Args[0][2:]
		}
		fmt.Fprintf(os.Stderr, "\n%s: Display information about Network Inferface Cards (NICs)\n", pgmName)
		fmt.Fprintf(os.Stderr, "usage: %s [options]\n", pgmName)
		flag.PrintDefaults()
	}
	flag.Parse()

	if *argsVersion {
		fmt.Fprintf(os.Stderr, "version %s\n", version)
		os.Exit(1)
	}

	brief := !(*argsAllDetails)
	allNics := queryNetworkInterfaces()

	// output interface table

	if !brief {
		renderTable(allNics, brief)
		return
	}

	briefNics := createBriefNicList(allNics, *argsDebug)
	renderTable(briefNics, brief)

	// output gateway table

	allGateways := GetAllGateways(allNics)
	fmt.Println("allGateways", allGateways)
}
