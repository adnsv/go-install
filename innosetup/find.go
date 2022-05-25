package innosetup

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/adnsv/go-utils/filesystem"
)

var NotFound = errors.New("innosetup compiler (iscc) not found")

func Find() (exe string, err error) {
	exe, err = exec.LookPath("ISCC")
	if err == nil {
		return
	}

	path := filepath.Join(os.Getenv("ProgramFiles(x86)"), "Inno Setup 6")
	if filesystem.FileExists(filepath.Join(path, "ISCC.exe")) {
		return filepath.Join(path, "ISCC"), nil
	}

	path = filepath.Join(os.Getenv("ProgramFiles"), "Inno Setup 6")
	if filesystem.FileExists(filepath.Join(path, "ISCC.exe")) {
		return filepath.Join(path, "ISCC"), nil
	}

	return "iscc", NotFound
}
