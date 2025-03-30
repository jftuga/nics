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
	"fmt"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/net/route"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// ParseDHCPInfo parses the input string and extracts DHCP information.
// It returns a map containing LeaseStartTime, LeaseExpirationTime, server_identifier,
// and lease_time (converted to decimal).
func ParseDHCPInfo(input string) (map[string]string, error) {
	result := make(map[string]string)

	// Extract LeaseStartTime
	leaseStartPattern := regexp.MustCompile(`LeaseStartTime : ([^\n]+)`)
	leaseStartMatch := leaseStartPattern.FindStringSubmatch(input)
	if len(leaseStartMatch) > 1 {
		result["LeaseStartTime"] = strings.TrimSpace(leaseStartMatch[1])
	}

	// Extract LeaseExpirationTime
	leaseExpPattern := regexp.MustCompile(`LeaseExpirationTime : ([^\n]+)`)
	leaseExpMatch := leaseExpPattern.FindStringSubmatch(input)
	if len(leaseExpMatch) > 1 {
		result["LeaseExpirationTime"] = strings.TrimSpace(leaseExpMatch[1])
	}

	// Extract server_identifier
	serverIdPattern := regexp.MustCompile(`server_identifier \(ip\): ([^\n]+)`)
	serverIdMatch := serverIdPattern.FindStringSubmatch(input)
	if len(serverIdMatch) > 1 {
		result["server_identifier"] = strings.TrimSpace(serverIdMatch[1])
	}

	// Extract lease_time and convert to decimal
	leaseTimePattern := regexp.MustCompile(`lease_time \(uint32\): (0x[0-9a-fA-F]+)`)
	leaseTimeMatch := leaseTimePattern.FindStringSubmatch(input)
	if len(leaseTimeMatch) > 1 {
		hexValue := strings.TrimSpace(leaseTimeMatch[1])
		// Convert hex to decimal
		hexValue = strings.TrimPrefix(hexValue, "0x")
		decimalValue, err := strconv.ParseInt(hexValue, 16, 64)
		if err != nil {
			return result, fmt.Errorf("failed to convert lease_time to decimal: %v", err)
		}
		result["lease_time"] = fmt.Sprintf("%d", decimalValue)
		formattedTime, err := FormatLeaseTime(result["lease_time"])
		if err != nil {
			return result, fmt.Errorf("failed to format lease_time: %v", err)
		}
		result["formatted_lease_time"] = formattedTime
	}
	return result, nil
}

// FormatLeaseTime converts lease time in seconds to a human-readable format
// showing days, hours, minutes, and seconds.
func FormatLeaseTime(secondsStr string) (string, error) {
	// Convert string to integer
	seconds, err := strconv.ParseInt(secondsStr, 10, 64)
	if err != nil {
		return "", fmt.Errorf("failed to parse lease time: %v", err)
	}

	// Calculate days, hours, minutes, seconds
	days := seconds / (24 * 60 * 60)
	seconds %= 24 * 60 * 60

	hours := seconds / (60 * 60)
	seconds %= 60 * 60

	minutes := seconds / 60
	seconds %= 60

	// Format the result
	return fmt.Sprintf("%d days, %d hrs, %d mins, %d secs", days, hours, minutes, seconds), nil
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
func parseScutilOutput(output string) []string {
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
	dns = parseScutilOutput(string(output))
	return dns
}

func getMacOSDhcp(allAdapters []string) {
	var allDhcpInfo []map[string]string
	for _, adapter := range allAdapters {
		dhcpCmd := exec.Command("/usr/sbin/ipconfig", "getsummary", adapter)
		output, err := dhcpCmd.CombinedOutput()
		if err != nil {
			continue
		}
		dhcpInfo, err := ParseDHCPInfo(string(output))
		if err != nil {
			continue
		}
		dhcpInfo["adapter"] = adapter
		allDhcpInfo = append(allDhcpInfo, dhcpInfo)
	}
	renderDHCPTable(allDhcpInfo)
}

func renderDHCPTable(allDhcpInfo []map[string]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"Name", "DHCP Server", "Lease Start", "Lease Expiration", "Lease Duration"})
	for _, adapter := range allDhcpInfo {
		shortLeaseDur := ShortenLeaseDuration(adapter["formatted_lease_time"])
		table.Append([]string{adapter["adapter"], adapter["server_identifier"], adapter["LeaseStartTime"], adapter["LeaseExpirationTime"], shortLeaseDur})
	}
	table.Render()
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

func gatewayAndDNS(allIPv4, allIPv6, allRenderedInterfaces []string, brief bool) {
	getMacOSDhcp(allRenderedInterfaces)

	gateway := getMacOSDefaultGateway()
	dns := getMacOSDNS()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"Gateway", "DNS 1", "DNS 2"})
	table.Append([]string{gateway, dns[0], dns[1]})
	table.Render()
}
