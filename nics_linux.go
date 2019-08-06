/*
nics_linux.go
-John Taylor
2019-08-03

Display information about Network Inferface Cards (NICs)

MIT License; Copyright (c) 2019 John Taylor
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

*/

// +build linux

package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// adapted from: https://raw.githubusercontent.com/fiskeben/resolv/master/main.go

type Resolver struct {
	Domains     []string
	Nameservers []string
	Search      []string
	Sortlist    []string
}

// Config reads /etc/resolv.conf and returns it as a Resolver
func Config() (Resolver, error) {
	f, err := os.Open("/etc/resolv.conf")
	if err != nil {
		return Resolver{}, err
	}
	defer f.Close()
	return parse(f)
}

func parse(f io.Reader) (Resolver, error) {
	domains := make([]string, 0)
	nameservers := make([]string, 0)
	search := make([]string, 0)
	sortlist := make([]string, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, " ")
		if len(parts) < 2 {
			continue
		}

		kind := parts[0]
		rest := parts[1:]

		switch kind {
		case "domain":
			for _, d := range rest {
				d := strings.TrimSpace(d)
				if d != "" {
					domains = append(domains, d)
				}
			}
		case "nameserver":
			n := strings.Join(rest, "")
			n = strings.TrimSpace(n)
			nameservers = append(nameservers, n)
		case "search":
			for _, s := range rest {
				s := strings.TrimSpace(s)
				if s != "" {
					search = append(domains, s)
				}
			}
		case "sortlist":
			for _, s := range rest {
				s := strings.TrimSpace(s)
				if s != "" {
					sortlist = append(domains, s)
				}
			}
		}
	}

	return Resolver{
		Domains:     domains,
		Nameservers: nameservers,
		Search:      search,
		Sortlist:    sortlist,
	}, nil
}

func getGatewaysAndDHCP(brief bool) (map[string]string, error) {
	return nil, nil
}

// adapted from: https://stackoverflow.com/a/40695315/452281

/* /proc/net/route file:
Iface   Destination Gateway     Flags   RefCnt  Use Metric  Mask
eno1    00000000    C900A8C0    0003    0   0   100 00000000    0   00
eno1    0000A8C0    00000000    0001    0   0   100 00FFFFFF    0   00
*/

const (
	file  = "/proc/net/route"
	line  = 1    // line containing the gateway addr. (first line: 0)
	sep   = "\t" // field separator
	field = 2    // field containing hex gateway address (first field: 0)
)

func getGWFromRouteTable() string {
	file, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	gateway := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		// jump to line containing the agteway address
		for i := 0; i < line; i++ {
			scanner.Scan()
		}

		// get field containing gateway address
		tokens := strings.Split(scanner.Text(), sep)
		gatewayHex := "0x" + tokens[field]

		// cast hex address to uint32
		d, _ := strconv.ParseInt(gatewayHex, 0, 64)
		d32 := uint32(d)

		// make net.IP address from uint32
		ipd32 := make(net.IP, 4)
		binary.LittleEndian.PutUint32(ipd32, d32)
		//fmt.Printf("%T --> %[1]v\n", ipd32)

		// format net.IP to dotted ipV4 string
		gateway = net.IP(ipd32).String()
		//fmt.Printf("%T --> %[1]v\n", gateway)

		// exit scanner
		break
	}
	return gateway
}

func gatewayAndDNS(allIPv4, allIPv6 []string, brief bool) {
	conf, err := Config()
	if err != nil {
		fmt.Println(err)
		return
	}

	dns := make([]string, 5)
	for i := range conf.Nameservers {
		dns[i] = conf.Nameservers[i]
	}
	gateway := getGWFromRouteTable()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"Gateway", "DNS1", "DNS2"})
	table.Append([]string{gateway, dns[0], dns[1]})

	table.Render()
}
