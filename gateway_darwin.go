
// +build darwin

/*
gateway_darwin.go
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
	"golang.org/x/net/route"
	"strconv"
	"strings"
)

var defaultRoute = [4]byte{0, 0, 0, 0}

func GetAllGateways(allNics []Nic) []string {
	rib, _ := route.FetchRIB(0, route.RIBTypeRoute, 0)
	messages, err := route.ParseRIB(route.RIBTypeRoute, rib)

	if err != nil {
		return nil
	}

	var quad []string
	var ip string
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
			//fmt.Println(gateway.IP)
			//fmt.Printf("%T %t %s\n", gateway.IP[0], gateway.IP[0], gateway.IP[0])

			for _,v := range gateway.IP {
				quad = append(quad, strconv.Itoa(int(v)))
			}
			//fmt.Println("q:",quad)
			ip = strings.Join(quad,".")
			break
		}
	}
	fmt.Println("ip:", ip)
	return []string{ip}
}

func GetGatewayForIP(ip string) (string, error) {
	return "",nil
}
