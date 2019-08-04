/*
nics.go
-John Taylor
2019-08-03

Display information about Network Inferface Cards (NICs)

To compile:
go build -ldflags="-s -w" nics.go

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

	"github.com/olekukonko/tablewriter"
)

const version = "1.1.0"

func isBriefEntry(ifaceName, macAddr, mtu, flags string, ipv4List, ipv6List []string, debug bool) bool {
	if debug {
		fmt.Println("isBriefEntry:", ifaceName)
	}
	if strings.Contains(flags, "loopback") {
		if debug {
			fmt.Println("   not_brief: loopback flag")
		}
		return false
	}
	if strings.HasPrefix(macAddr, "00:00:00:00:00:00") {
		if debug {
			fmt.Println("   not_brief: NULL macAddr")
		}
		return false
	}
	if 0 == len(ipv4List) {
		if debug {
			fmt.Println("   not_brief: no IP addresses")
		}
		return false
	}
	for _, ipv4 := range ipv4List {
		if strings.HasPrefix(ipv4, "169.254.") {
			if debug {
				fmt.Println("   not_brief: self assigned:", ipv4)
			}
			return false
		}
	}
	if debug {
		fmt.Println("    is_brief: true")
	}
	return true
}

func extractIPAddrs(ifaceName string, allAddresses []net.Addr, brief bool) ([]string, []string) {
	var allIPv4 []string
	var allIPv6 []string

	for _, netAddr := range allAddresses {
		address := netAddr.String()
		if strings.Contains(address, ":") {
			allIPv6 = append(allIPv6, address)
		} else {
			allIPv4 = append(allIPv4, address)
		}
	}
	return allIPv4, allIPv6
}

func networkInterfaces(brief bool, debug bool) {
	adapters, err := net.Interfaces()
	if err != nil {
		fmt.Print(fmt.Errorf("%+v\n", err.Error()))
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	if brief {
		table.SetHeader([]string{"Name", "IPv4", "Mac Address", "MTU", "Flags"})
	} else {
		table.SetHeader([]string{"Name", "IPv4", "IPv6", "Mac Address", "MTU", "Flags"})
	}

	for _, iface := range adapters {
		allAddresses, err := iface.Addrs()
		if err != nil {
			fmt.Print(fmt.Errorf("%+v\n", err.Error()))
			return
		}

		allIPv4, allIPv6 := extractIPAddrs(iface.Name, allAddresses, brief)
		if debug {
			fmt.Println()
			fmt.Println("---------------------")
			fmt.Println(iface.Name, allAddresses)
			fmt.Println("ipv4:", allIPv4)
			fmt.Println("ipv6:", allIPv6)

		}

		ifaceName := strings.ToLower(iface.Name)
		macAddr := iface.HardwareAddr.String()
		mtu := strconv.Itoa(iface.MTU)
		flags := iface.Flags.String()

		if brief && isBriefEntry(ifaceName, macAddr, mtu, flags, allIPv4, allIPv6, debug) {
			table.Append([]string{iface.Name, strings.Join(allIPv4, " "), macAddr, mtu, flags})
			continue
		}

		if !brief {
			table.Append([]string{ifaceName, strings.Join(allIPv4, "\n"), strings.Join(allIPv6, "\n"), macAddr, mtu, strings.Replace(flags, "|", "\n", -1)})
			table.SetAutoWrapText(true)
			table.SetRowLine(true)
		}
	}
	table.Render()
}

func main() {
	argsAllDetails := flag.Bool("a", false, "show all details on ALL interfaces")
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

	networkInterfaces(!(*argsAllDetails), *argsDebug)
}
