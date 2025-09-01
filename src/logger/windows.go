//go:build windows
// +build windows

package logger

import (
	"fmt"
	"os"
	"runtime"

	"golang.org/x/sys/windows"
)

// configureConsole sets up the Windows console to minimize blocking
func ConfigureConsole() {
	if runtime.GOOS != "windows" {
		return
	}
	// Disable QuickEdit mode
	handle, err := windows.GetStdHandle(windows.STD_INPUT_HANDLE)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get console handle: %v\n", err)
		return
	}
	var mode uint32
	if err := windows.GetConsoleMode(handle, &mode); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get console mode: %v\n", err)
		return
	}
	mode &^= windows.ENABLE_QUICK_EDIT_MODE
	if err := windows.SetConsoleMode(handle, mode); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to disable QuickEdit: %v\n", err)
	}
}
