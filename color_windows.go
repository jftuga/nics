//go:build windows

package main

import (
	"os"
	"syscall"
)

const enableVirtualTerminalProcessing = 0x0004

func init() {
	handle := syscall.Handle(os.Stdout.Fd())
	var mode uint32
	if err := syscall.GetConsoleMode(handle, &mode); err == nil {
		syscall.SetConsoleMode(handle, mode|enableVirtualTerminalProcessing)
	}
}
