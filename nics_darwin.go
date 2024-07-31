//go:build darwin
// +build darwin

/*
nics_darwin.go
-John Taylor
2019-08-03

Display information about Network Interface Cards (NICs)

MIT License; Copyright (c) 2019 John Taylor
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

*/

package main

import (
	"github.com/olekukonko/tablewriter"
	"golang.org/x/net/route"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func getGatewaysAndDHCP(brief bool) (map[string]string, error) {
	return nil, nil
}

// adopted from: https://stackoverflow.com/a/31221013/452281
func convert(b []byte) string {
	s := make([]string, len(b))
	for i := range b {
		s[i] = strconv.Itoa(int(b[i]))
	}
	return strings.Join(s, ".")
}

// extract the nameservers from this cmd-line output: scutil --dns
func parseOutput(output string) []string {
	var dns = []string{"N/A", "N/A"}
	allLines := strings.Split(output, "\n")
	i := 0
	resolverCount := 0
	for _, line := range allLines {
		//fmt.Println("rc:", resolverCount, "   i:", i, "   dns:", dns, "  line:", line)
		if strings.Contains(line, "resolver #") {
			resolverCount += 1
			if resolverCount >= 2 && i > 0 {
				break
			}
		}
		if strings.Contains(line, "nameserver[") {
			clean := strings.TrimSpace(line)
			slots := strings.Split(clean, " ")
			if len(slots) == 3 {
				dns[i] = slots[2]
				i += 1
				if i == 2 {
					break
				}
			}
		}
	}
	return dns
}

func getMacOSDNS() []string {
	var dns = []string{"N/A", "N/A"}

	dnsCmd := exec.Command("/usr/sbin/scutil", "--dns")
	output, err := dnsCmd.CombinedOutput()
	if err != nil {
		return dns
	}
	dns = parseOutput(string(output))
	return dns
}

// adopted from: https://gist.github.com/abimaelmartell/dcbbff464dc0778165b2dcc5092f90e6
func getMacOSDefaultGateway() string {
	var defaultRoute = [4]byte{0, 0, 0, 0}
	rib, _ := route.FetchRIB(0, route.RIBTypeRoute, 0)
	messages, err := route.ParseRIB(route.RIBTypeRoute, rib)

	if err != nil {
		return "N/A"
	}

	for _, message := range messages {
		route_message := message.(*route.RouteMessage)
		addresses := route_message.Addrs

		var destination, gateway *route.Inet4Addr
		ok := false

		if destination, ok = addresses[0].(*route.Inet4Addr); !ok {
			continue
		}

		if gateway, ok = addresses[1].(*route.Inet4Addr); !ok {
			continue
		}

		if destination == nil || gateway == nil {
			continue
		}

		if destination.IP == defaultRoute {
			return convert(gateway.IP[:])
		}
	}
	return "N/A"
}

func gatewayAndDNS(allIPv4, allIPv6 []string, brief bool) {

	gateway := getMacOSDefaultGateway()
	dns := getMacOSDNS()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"Gateway", "DNS 1", "DNS 2"})
	table.Append([]string{gateway, dns[0], dns[1]})
	table.Render()
}
