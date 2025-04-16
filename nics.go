/*
nics.go
-John Taylor
2019-08-03

Display information about Network Interface Cards (NICs)

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

const version = "1.6.2"

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

func networkInterfaces(brief bool, debug bool, singleInterface string) ([]string, []string, []string) {
	adapters, err := net.Interfaces()
	if err != nil {
		fmt.Print(fmt.Errorf("%+v\n", err.Error()))
		return nil, nil, nil
	}

	foundSingleInterface := false
	if len(singleInterface) > 0 {
		brief = false
		singleInterface = strings.ToLower(singleInterface)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	if brief {
		table.SetHeader([]string{"Name", "IP", "Mac Address", "MTU", "Flags"})
	} else {
		table.SetHeader([]string{"Name", "IPv4", "IPv6", "Mac Address", "MTU", "Flags"})
	}

	var v4Addresses []string
	var v6Addresses []string
	var allRenderedInterfaces []string
	for _, iface := range adapters {
		//fmt.Printf("%T %v\n", iface, iface)
		allAddresses, err := iface.Addrs()
		if err != nil {
			fmt.Print(fmt.Errorf("%+v\n", err.Error()))
			return nil, nil, nil
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
		if len(singleInterface) > 0 && ifaceName != singleInterface {
			continue
		} else if len(singleInterface) > 0 && ifaceName == singleInterface {
			foundSingleInterface = true
		}
		macAddr := iface.HardwareAddr.String()
		mtu := strconv.Itoa(iface.MTU)
		flags := iface.Flags.String()

		if brief && isBriefEntry(ifaceName, macAddr, mtu, flags, allIPv4, allIPv6, debug) {
			joined := strings.Join(allIPv4, "\n") // + "\n" + strings.Join(allIPv6, "\n")
			table.Append([]string{iface.Name, joined, macAddr, mtu, flags})
			allRenderedInterfaces = append(allRenderedInterfaces, iface.Name)
			for _, ipWithMask := range allIPv4 {
				ip := strings.Split(ipWithMask, "/")
				v4Addresses = append(v4Addresses, ip[0])
			}
			continue
		}

		if !brief {
			table.SetAutoWrapText(true)
			table.SetRowLine(true)
			table.Append([]string{ifaceName, strings.Join(allIPv4, "\n"), strings.Join(allIPv6, "\n"), macAddr, mtu, strings.Replace(flags, "|", "\n", -1)})
			allRenderedInterfaces = append(allRenderedInterfaces, iface.Name)
			for _, ipWithMask := range allIPv4 {
				ip := strings.Split(ipWithMask, "/")
				v4Addresses = append(v4Addresses, ip[0])
			}
			for _, ipWithMask := range allIPv6 {
				ip := strings.Split(ipWithMask, "/")
				v6Addresses = append(v6Addresses, ip[0])
			}
		}
	}
	if len(singleInterface) > 0 && !foundSingleInterface {
		_, _ = fmt.Fprintf(os.Stderr, "\ninterface not found: %v\n", singleInterface)
	} else {
		table.Render()
	}

	return v4Addresses, v6Addresses, allRenderedInterfaces
}

// FormatWithCorrectPlurals ensures time units use singular form when the value is 1.
// For example:
// - "1 days" becomes "1 day"
// - "2 mins" remains "2 mins"
// - "0 hrs, 1 mins" becomes "0 hrs, 1 min"
//
// The function handles days, hrs, mins, and secs units.
func FormatWithCorrectPlurals(duration string) string {
	// Replace "1 days" with "1 day"
	duration = strings.Replace(duration, "1 days", "1 day", 1)

	// Replace "1 hrs" with "1 hr"
	duration = strings.Replace(duration, "1 hrs", "1 hr", 1)

	// Replace "1 mins" with "1 min"
	duration = strings.Replace(duration, "1 mins", "1 min", 1)

	// Replace "1 secs" with "1 sec"
	duration = strings.Replace(duration, "1 secs", "1 sec", 1)

	return duration
}

// ShortenLeaseDuration trims trailing zero units from a lease duration string.
//
// For example:
// - "1 days, 0 hrs, 0 mins, 0 secs" becomes "1 day"
// - "1 days, 0 hrs, 5 mins, 0 secs" becomes "1 day, 0 hrs, 5 mins"
// - "1 days, 0 hrs, 0 mins, 1 secs" becomes "1 day, 0 hrs, 0 mins, 1 sec"
//
// The function preserves any non-zero units and removes only trailing zero units.
// It also ensures correct singular/plural forms.
func ShortenLeaseDuration(leaseDuration string) string {
	// Split the duration into its components
	components := strings.Split(leaseDuration, ", ")

	// Find the last non-zero component
	lastNonZeroIndex := len(components) - 1
	for i := len(components) - 1; i > 0; i-- {
		if !strings.HasPrefix(components[i], "0 ") {
			break
		}
		lastNonZeroIndex = i - 1
	}

	// If all components after the first one are zero, return just the first component
	if lastNonZeroIndex == 0 {
		return FormatWithCorrectPlurals(components[0])
	}

	// Join the components up to and including the last non-zero one
	shortened := strings.Join(components[:lastNonZeroIndex+1], ", ")

	// Apply pluralization formatting
	return FormatWithCorrectPlurals(shortened)
}

func main() {
	argsAllDetails := flag.Bool("a", false, "show all details on ALL interfaces, includes DHCP info on Windows")
	argsDebug := flag.Bool("d", false, "show debug information")
	argsVersion := flag.Bool("v", false, "show program version")
	argsSingleInterface := flag.String("i", "", "interface name")

	flag.Usage = func() {
		pgmName := os.Args[0]
		if strings.HasPrefix(os.Args[0], "./") {
			pgmName = os.Args[0][2:]
		}
		fmt.Fprintf(os.Stderr, "\n%s: Display information about Network Interface Cards (NICs)\n", pgmName)
		fmt.Fprintf(os.Stderr, "usage: %s [options]\n", pgmName)
		flag.PrintDefaults()
	}
	flag.Parse()

	if *argsVersion {
		fmt.Fprintf(os.Stderr, "version %s\n", version)
		fmt.Fprintf(os.Stderr, "https://github.com/jftuga/nics\n")
		os.Exit(0)
	}

	allIPv4, allIPv6, allRenderedInterfaces := networkInterfaces(!(*argsAllDetails), *argsDebug, *argsSingleInterface)
	gatewayAndDNS(allIPv4, allIPv6, allRenderedInterfaces, !(*argsAllDetails))
}
