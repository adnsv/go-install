// +build !windows

package dep

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const EXE = ""

func IsW32onW64() bool {
	return false
}

func IsPosixSudo() (bool, error) {
	stdout, err := exec.Command("ps", "-o", "user=", "-p", strconv.Itoa(os.Getpid())).Output()
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(stdout)) == "root", nil
}
