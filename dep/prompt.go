package dep

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/adnsv/go-utils/filesystem"
	"github.com/blang/semver/v4"
	"github.com/manifoldco/promptui"
)

type TargetPrompt struct {
	DefaultDir   string
	MainExe      string
	Lookup       bool
	CheckVerArgs []string
	AllowCustom  bool
	AvoidGoBin   bool
}

func execute(cmd string, args ...string) string {
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func insubdir(parent, child string) bool {
	if parent == "" {
		return false
	}
	rp, _ := filepath.Rel(parent, child)
	return !strings.HasPrefix(rp, "..")
}

var ErrPromptCancelled = errors.New("Prompt cancelled")

func (tp *TargetPrompt) Run() (string, error) {
	prev := ""
	if tp.Lookup && tp.MainExe != "" {
		lookup, err := exec.LookPath(tp.MainExe)
		if err == nil {
			p := filepath.Dir(lookup)
			avoids := []string{}
			if tp.AvoidGoBin {
				if s := os.Getenv("GOPATH"); s != "" {
					avoids = append(avoids, os.ExpandEnv(s))
				} else {
					s := os.ExpandEnv("${HOME}/go")
					if filesystem.DirExists(s) {
						avoids = append(avoids, s)
					}
				}
				if s := os.Getenv("GOROOT"); s != "" {
					avoids = append(avoids, os.ExpandEnv(s))
				}
			}
			avoid := false
			for _, d := range avoids {
				if insubdir(d, p) {
					avoid = true
					break
				}
			}
			if !avoid {
				prev = p
			}
		}
	}

	type entry struct {
		id        int
		path      string
		exefound  bool
		exevererr error
		exesemver *semver.Version
		comment   string
	}

	entries := []*entry{}

	if prev != "" {
		entries = append(entries, &entry{id: 0, path: prev, comment: "active"})
	}
	if tp.DefaultDir != "" && tp.DefaultDir != prev {
		entries = append(entries, &entry{id: 1, path: tp.DefaultDir, comment: "recommended"})
	}
	if tp.CheckVerArgs != nil {
		for _, e := range entries {
			e.exefound = false
			if filesystem.FileExists(e.path, tp.MainExe) {
				e.exefound = true
				out, err := exec.Command(filepath.Join(e.path, tp.MainExe), tp.CheckVerArgs...).Output()
				if err != nil {
					e.exevererr = err
				} else {
					v, err := semver.ParseTolerant(strings.TrimSpace(string(out)))
					if err != nil {
						e.exevererr = err
					} else {
						e.exesemver = &v
					}
				}
			}
		}
	}
	if tp.AllowCustom {
		entries = append(entries, &entry{id: -2, comment: "Select custom location ..."})
	}
	if tp.AllowCustom {
		entries = append(entries, &entry{id: -1, comment: "Cancel"})
	}

	choices := []string{}
	for _, e := range entries {
		s := e.comment
		if e.id >= 0 {
			s = e.path
			if e.exesemver != nil {
				s += fmt.Sprintf(" (%s)", e.exesemver.String())
			}
			if e.comment != "" {
				s += fmt.Sprintf(" [%s]", e.comment)
			}
		}
		choices = append(choices, "• "+s)
	}
	prompt := promptui.Select{
		Label:    "where to? ↑↓",
		Items:    choices,
		HideHelp: true,
	}
	i, _, err := prompt.Run()
	if err != nil {
		return "", err
	}
	e := entries[i]
	if e.id == -1 {
		// cancel
		return "", ErrPromptCancelled
	} else if e.id == -2 {
		// custom
		dp := promptui.Prompt{
			Label:       "install to",
			HideEntered: true,
			Validate: func(input string) error {
				if input == "" {
					return errors.New("empty string is not allowed")
				}
				return nil
			},
		}
		dir, err := dp.Run()
		if err != nil {
			return "", err
		}
		dir = os.ExpandEnv(dir)
		vdir, err := filepath.Abs(dir)
		if err != nil {
			return "", err
		}
		return filepath.FromSlash(vdir), err
	} else if e.id == -1 {
		// cancel
		return "", ErrPromptCancelled
	}
	vdir, err := filepath.Abs(e.path)
	return filepath.FromSlash(vdir), err
}
