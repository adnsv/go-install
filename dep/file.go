package dep

import (
	"bytes"
	"compress/zlib"
	"io"
	"os"
	"path/filepath"
)

type File struct {
	Path string
	Perm os.FileMode
	Data []byte
}

func (f *File) Expand(dir string) error {
	in := bytes.NewReader(f.Data)
	r, err := zlib.NewReader(in)
	if err != nil {
		return err
	}
	defer r.Close()

	out, err := os.OpenFile(
		filepath.Join(dir, f.Path),
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Perm)
	if err != nil {
		return nil
	}
	defer out.Close()
	_, err = io.Copy(out, r)
	return err
}
