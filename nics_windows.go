//go:build windows
// +build windows

/*
nics_windows.go
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
	"os"
	"syscall"
	"unsafe"

	"github.com/olekukonko/tablewriter"
)

const (
	MAX_HOSTNAME_LEN    = 128
	MAX_DOMAIN_NAME_LEN = 128
	MAX_SCOPE_ID_LEN    = 256
)

// https://docs.microsoft.com/en-us/windows/win32/api/iptypes/ns-iptypes-ip_addr_string
type ipAddressString struct {
	next    *ipAddressString
	address [16]byte
	mask    [16]byte
	context uint32
}

// https://docs.microsoft.com/en-us/windows/win32/api/iptypes/ns-iptypes-fixed_info_w2ksp1
type fixedInfo struct {
	hostName         [MAX_HOSTNAME_LEN + 4]byte
	domainName       [MAX_DOMAIN_NAME_LEN + 4]byte
	currentDNSServer *ipAddressString
	dnsServerList    ipAddressString
	nodeType         uint16
	scopeID          [MAX_SCOPE_ID_LEN + 4]byte
	enableRouting    uint16
	enableProxy      uint16
	enableDNS        uint16
}

const (
	MAX_ADAPTER_NAME_LENGTH        = 256
	MAX_ADAPTER_DESCRIPTION_LENGTH = 128
	MAX_ADAPTER_ADDRESS_LENGTH     = 8
)

// https://docs.microsoft.com/en-us/windows/win32/api/iptypes/ns-iptypes-ip_adapter_info
type ipAdapterInfo struct {
	next               *ipAdapterInfo
	comboIndex         uint32
	adapterName        [MAX_ADAPTER_NAME_LENGTH + 4]byte
	adapterDescription [MAX_ADAPTER_DESCRIPTION_LENGTH + 4]byte
	addressLength      uint16
	address            [MAX_ADAPTER_ADDRESS_LENGTH]byte
	index              uint32
	adapterType        uint32
	dhcpEnabled        uint16
	currentAddress     *ipAddressString
	ipAddressList      ipAddressString
	gatewayList        ipAddressString
	dhcpServer         ipAddressString
	haveWins           bool
	primaryWins        ipAddressString
	secondaryWins      ipAddressString
	leaseObtained      uint64
	leaseExpires       uint64
}

var (
	getNetworkParams = syscall.NewLazyDLL("Iphlpapi.dll").NewProc("GetNetworkParams")
	getAdaptersInfo  = syscall.NewLazyDLL("Iphlpapi.dll").NewProc("GetAdaptersInfo")
)

func getDNSEntries() ([]string, error) {
	err := getNetworkParams.Find()
	if err != nil {
		return nil, err
	}

	netInfo := [2]fixedInfo{}
	size := unsafe.Sizeof(netInfo)

	r0, _, err := getNetworkParams.Call(uintptr(unsafe.Pointer(&netInfo[0])), uintptr(unsafe.Pointer(&size)))
	if r0 != 0 {
		return nil, err
	}

	var dns1, dns2 string
	dns1 = sliceToString(netInfo[0].dnsServerList.address[:])

	nextDNS := netInfo[0].dnsServerList.next
	if nextDNS != nil {
		dns2 = sliceToString(nextDNS.address[:])
	}

	return []string{dns1, dns2}, nil
}

func getGatewaysAndDHCP(brief bool) (map[string]string, error) {
	err := getAdaptersInfo.Find()
	if err != nil {
		return nil, err
	}

	adapters := [16]ipAdapterInfo{}
	size := unsafe.Sizeof(adapters)

	result, _, err := getAdaptersInfo.Call(uintptr(unsafe.Pointer(&adapters[0])), uintptr(unsafe.Pointer(&size)))
	if result != 0 {
		return nil, err
	}

	ipMapDHCP := make(map[string][]string) // key:adapter IP; values:0 = dhcp server ip, 1 = leaseObtained, 2 = leaseExpires
	ipMapGateway := make(map[string]string)

	adapter := &adapters[0]
	for adapter != nil {
		/*
			fmt.Println("  ip:", sliceToString(adapter.ipAddressList.address[:]))
			fmt.Println("gate:", sliceToString(adapter.gatewayList.address[:]))
			fmt.Println("  primaryWins:", sliceToString(adapter.primaryWins.address[:]))
			fmt.Println("secondaryWins:", sliceToString(adapter.secondaryWins.address[:]))
			fmt.Println()
		*/
		ip := sliceToString(adapter.ipAddressList.address[:])
		gate := sliceToString(adapter.gatewayList.address[:])
		dhcpServer := sliceToString(adapter.dhcpServer.address[:])
		leaseObtained := adapter.leaseObtained
		leaseExpires := adapter.leaseExpires

		if ip != "0.0.0.0" && gate != "0.0.0.0" {
			ipMapGateway[ip] = gate
		}
		if len(dhcpServer) >= 4 {
			ipMapDHCP[ip] = []string{dhcpServer, timeToString(leaseObtained), timeToString(leaseExpires)}
		}

		adapter = adapter.next
	}

	if !brief && len(ipMapDHCP) > 0 {
		renderDHCPTable(ipMapDHCP)
		fmt.Println()
	}

	return ipMapGateway, nil
}

func renderDHCPTable(ipMapDHCP map[string][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"IP", "DHCP Server", "Lease Renewed", "Lease Expires"})
	for ip, dhcpInfo := range ipMapDHCP {
		table.Append([]string{ip, dhcpInfo[0], dhcpInfo[1], dhcpInfo[2]})
	}
	table.Render()
}

func gatewayAndDNS(allIPv4, allIPv6 []string, brief bool) {
	var err error
	var dns = []string{"N/A", "N/A"}
	dns, err = getDNSEntries()
	if err != nil {
		fmt.Println(err)
	}

	ipMapGateway := make(map[string]string)
	ipMapGateway, err = getGatewaysAndDHCP(brief)
	if err != nil {
		fmt.Println(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"Gateway", "DNS 1", "DNS 2"})
	for ip := range ipMapGateway {
		if arrayContains(ip, allIPv4) {
			table.Append([]string{ipMapGateway[ip], dns[0], dns[1]})
		}
		if arrayContains(ip, allIPv6) {
			table.Append([]string{ipMapGateway[ip], dns[0], dns[1]})
		}
		dns[0] = ""
		dns[1] = ""
	}
	table.Render()
}
