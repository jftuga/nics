// +build windows

/*
gateway_windows.go
-John Taylor
2019-08-12

Display information about Network Inferface Cards (NICs)

MIT License; Copyright (c) 2019 John Taylor
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

*/

package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	MAX_ADAPTER_NAME_LENGTH        = 256
	MAX_ADAPTER_DESCRIPTION_LENGTH = 128
	MAX_ADAPTER_ADDRESS_LENGTH     = 8
)

// https://docs.microsoft.com/en-us/windows/win32/api/iptypes/ns-iptypes-ip_addr_string
type ipAddressString struct {
	next    *ipAddressString
	address [16]byte
	mask    [16]byte
	context uint32
}

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
	getAdaptersInfo = syscall.NewLazyDLL("Iphlpapi.dll").NewProc("GetAdaptersInfo")
)

func GetGatewayForIP(ip string) (string, error) {
	err := getAdaptersInfo.Find()
	if err != nil {
		return "", err
	}

	adapters := [16]ipAdapterInfo{}
	size := unsafe.Sizeof(adapters)

	result, _, err := getAdaptersInfo.Call(uintptr(unsafe.Pointer(&adapters[0])), uintptr(unsafe.Pointer(&size)))
	if result != 0 {
		return "", err
	}

	adapter := &adapters[0]
	for adapter != nil {
		currentIP := sliceToString(adapter.ipAddressList.address[:])
		currentGateway := sliceToString(adapter.gatewayList.address[:])

		if ip == currentIP {
			return currentGateway, nil
		}

		adapter = adapter.next
	}

	return "", nil
}

func GetAllGateways(allNics []Nic) []string {
	var allGateways []string

	for _, n := range allNics {
		for _, a := range n.Addrs {
			ip := a.String()
			fmt.Println("ip:", ip)
			gw, err := GetGatewayForIP(ip)
			if err != nil {
				continue
			}
			if len(gw) > 0 {
				allGateways = append(allGateways, gw)
			}
		}
	}
	return allGateways
}
