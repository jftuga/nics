/*
crossPlatform.go
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
	"bytes"
	"fmt"
	"strings"
	"time"
)

func sliceToString(slice []byte) string {
	n := bytes.IndexByte(slice, 0)
	return string(slice[:n])
}

func timeToString(input uint64) string {
	t := time.Unix(int64(input), 0)
	t = t.In(time.Local)
	raw := fmt.Sprintf("%s", t)
	slots := strings.Split(raw, " ")
	return fmt.Sprintf("%s %s", slots[0], slots[1])
}

func arrayContains(value string, array []string) bool {
	for _, a := range array {
		if a == value {
			return true
		}
	}
	return false
}

var noColor bool

const (
	ansiCyan  = "\033[36m"
	ansiReset = "\033[0m"
)

// ColorizeIP wraps an IP address string in ANSI cyan color codes.
func ColorizeIP(ip string) string {
	if noColor || ip == "" || ip == "N/A" {
		return ip
	}
	return ansiCyan + ip + ansiReset
}

// ColorizeCIDR colorizes just the IP portion of a CIDR notation address (e.g., 172.22.2.74/24).
func ColorizeCIDR(cidr string) string {
	parts := strings.SplitN(cidr, "/", 2)
	if len(parts) == 2 {
		return ColorizeIP(parts[0]) + "/" + parts[1]
	}
	return ColorizeIP(cidr)
}

// ColorizeIPList applies ColorizeCIDR to each element in a slice.
func ColorizeIPList(ips []string) []string {
	result := make([]string, len(ips))
	for i, ip := range ips {
		result[i] = ColorizeCIDR(ip)
	}
	return result
}
