/*
outputTable.go
-John Taylor
2019-08-14

Display information about Network Inferface Cards (NICs)

MIT License; Copyright (c) 2019 John Taylor
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

*/

package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func renderTable(allNics []Nic, brief bool) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	if !brief {
		table.SetRowLine(true)
	}
	if brief {
		table.SetHeader([]string{"Name", "IPv4", "Mac Address", "MTU", "Flags"})
	} else {
		table.SetHeader([]string{"Name", "IPv4", "IPv6", "Mac Address", "MTU", "Flags"})
	}

	for _, n := range allNics {
		macAddr := n.iFace.HardwareAddr.String()
		mtu := strconv.Itoa(n.iFace.MTU)
		flags := n.iFace.Flags.String()
		if !brief {
			flags = strings.Replace(flags, "|", "\n", -1)
		}
		var ipv4, ipv6 []string
		for _, a := range n.Addrs {
			ip := a.String()

			if brief && strings.HasPrefix(ip, "169.254") {
				continue
			}

			if strings.Contains(ip, ":") {
				ipv6 = append(ipv6, ip)
			} else {
				ipv4 = append(ipv4, ip)
			}
		}

		if brief {
			table.Append([]string{n.iFace.Name, strings.Join(ipv4, "\n"), macAddr, mtu, flags})
		} else {
			table.Append([]string{n.iFace.Name, strings.Join(ipv4, "\n"), strings.Join(ipv6, "\n"), macAddr, mtu, flags})
		}
	}

	table.Render()
}
