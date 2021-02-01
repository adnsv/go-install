package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adnsv/go-utils/binpack"
	"github.com/adnsv/go-utils/fs"
	"github.com/adnsv/slog"
)

type File struct {
	Ident string
	Path  string
	Perm  os.FileMode
}

func ToIdent(s string) string {
	ret := ""
	for i, c := range strings.ToLower(s) {
		switch {
		case c >= 0 && c <= 9:
			if i == 0 {
				ret += "_"
			}
			ret += string(c)
		case c >= 'a' && c <= 'z':
			ret += string(c)
		default:
			ret += "_"
		}
	}
	return ret
}

func Pack(binfile string, unpackfile string, sourcedir string, ff []*File) error {
	bins := []*binpack.Source{}
	for _, f := range ff {
		if f.Ident == "" {
			f.Ident = ToIdent(filepath.Base(f.Path))
		}
		slog.Infof("compressing %s ... ", f.Path)
		s, err := binpack.FromFile(filepath.Join(sourcedir, f.Path), f.Ident, true)
		if err != nil {
			return err
		}
		slog.Printf("%s -> %s",
			fs.ByteSizeStr(uint64(s.OriginalSize)),
			fs.ByteSizeStr(uint64(len(s.Content))))
		bins = append(bins, s)
	}
	out, err := os.OpenFile(binfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	slog.Info("generating binary resources")
	binpack.MakeGolang(out, "main", bins)
	out.Close()

	slog.Info("generating unpacker")
	out, err = os.OpenFile(unpackfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	fmt.Fprintln(out, "package main")
	fmt.Fprintln(out)
	fmt.Fprintln(out, `import "github.com/adnsv/go-install/dep"`)
	fmt.Fprintln(out)
	fmt.Fprintln(out, "// DO NOT EDIT: Generated file")
	fmt.Fprintln(out)
	if len(ff) > 0 {
		fmt.Fprintf(out, "var files = []dep.File{\n")
		for _, f := range ff {
			fmt.Fprintf(out, "\t{Data: %s, Path: %q, Perm: 0%o},\n",
				f.Ident, f.Path, f.Perm)
		}
		fmt.Fprintf(out, "}\n")
	}

	out.Close()

	return nil
}
