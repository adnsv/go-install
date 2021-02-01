// +build windows

package dep

import (
	"runtime"

	"golang.org/x/sys/windows"
)

const EXE = ".exe"

// IsW32onW64 detects 32-bit code running on 64-bit windows host
func IsW32onW64() bool {
	if runtime.GOOS == "windows" && runtime.GOARCH == "386" {
		var w64 bool
		windows.IsWow64Process(windows.CurrentProcess(), &w64)
		return w64
	} else {
		return false
	}
}

func IsPosixSudo() (bool, error) {
	return false, nil
}
