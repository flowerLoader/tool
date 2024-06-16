//go:build windows
// +build windows

package main

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

func osInit() {
	// if !mousetrap.StartedByExplorer() {
	// 	return
	// }

	setCloseButtonEnabled(false)
	setConsoleFontSize(27)
}

var (
	kernel32                = windows.NewLazySystemDLL("kernel32.dll")
	getConsoleWindow        = kernel32.NewProc("GetConsoleWindow")
	getStdHandle            = kernel32.NewProc("GetStdHandle")
	setCurrentConsoleFontEx = kernel32.NewProc("SetCurrentConsoleFontEx")
	setConsoleCP            = kernel32.NewProc("SetConsoleCP")
	setConsoleOutputCP      = kernel32.NewProc("SetConsoleOutputCP")

	user32         = windows.NewLazySystemDLL("user32.dll")
	enableMenuItem = user32.NewProc("EnableMenuItem")
	getSystemMenu  = user32.NewProc("GetSystemMenu")
)

func setCloseButtonEnabled(enabled bool) {
	const (
		SC_CLOSE    = 0xF060
		MF_ENABLED  = 0x00000000
		MF_GRAYED   = 0x00000001
		MF_DISABLED = 0x00000002
	)

	// Get our hWnd, then get the system menu for the console window.
	hWnd, _, _ := getConsoleWindow.Call()
	hMenu, _, _ := getSystemMenu.Call(hWnd, 0)

	// Disable the close button.
	if enabled {
		enableMenuItem.Call(hMenu, SC_CLOSE, MF_ENABLED)
	} else {
		enableMenuItem.Call(hMenu, SC_CLOSE, MF_GRAYED)
	}
}

func setConsoleFontSize(fontSize int) {
	stdOutputHandle, _, _ := getStdHandle.Call(windows.STD_OUTPUT_HANDLE)
	type CONSOLE_FONT_INFOEX struct {
		cbSize     uint32
		nFont      uint32
		dwFontSize windows.Coord
		FontFamily uint32
		FontWeight uint32
		FaceName   [32]uint16
	}
	// var cp uint32 = 437 // IBM OEM
	var cp uint32 = 65001 // UTF-8
	setConsoleCP.Call(uintptr(cp))
	setConsoleOutputCP.Call(uintptr(cp))
	setCurrentConsoleFontEx.Call(
		stdOutputHandle,
		0,
		uintptr(unsafe.Pointer(&CONSOLE_FONT_INFOEX{
			cbSize:     uint32(unsafe.Sizeof(CONSOLE_FONT_INFOEX{})),
			nFont:      0,
			dwFontSize: windows.Coord{X: 0, Y: int16(fontSize)},
			FontFamily: 54,  // FF_MODERN << 4 | TMPF_VECTOR | TMPF_TRUETYPE
			FontWeight: 400, // 400 = normal, 700 = bold
			FaceName:   stringToUint16Array("Consolas"),
		})),
	)
}

func stringToUint16Array(s string) [32]uint16 {
	utf16, err := syscall.UTF16FromString(s)
	if err != nil {
		panic(err)
	}

	var arr [32]uint16
	copy(arr[:], utf16)
	return arr
}
