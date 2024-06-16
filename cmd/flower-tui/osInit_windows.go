//go:build windows
// +build windows

package main

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

func osInit() {
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	setConsoleMode := kernel32.NewProc("SetConsoleMode")
	stdout := windows.Handle(windows.Stdout)

	// Set all of the following:
	// ENABLE_VIRTUAL_TERMINAL_PROCESSING
	// DISABLE_NEWLINE_AUTO_RETURN
	// ENABLE_PROCESSED_OUTPUT
	// ENABLE_WRAP_AT_EOL_OUTPUT
	var mode uint32 = 0x0004 | 0x0008 | 0x0001 | 0x0002
	setConsoleMode.Call(uintptr(stdout), uintptr(unsafe.Pointer(&mode)))

	// Hide the cursor
	cursorInfo := struct {
		size    uint32
		visible int32
	}{
		size:    1,
		visible: 0,
	}
	setConsoleCursorInfo := kernel32.NewProc("SetConsoleCursorInfo")
	setConsoleCursorInfo.Call(uintptr(stdout), uintptr(unsafe.Pointer(&cursorInfo)))
}
