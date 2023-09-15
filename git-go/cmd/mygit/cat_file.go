package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path"
)

func catFile(sha string) error {
	f, err := os.Open(path.Join(".git/objects", sha[:2], sha[2:]))
	if err != nil {
		return err
	}
	defer f.Close()

	r, err := zlib.NewReader(f)
	if err != nil {
		return err

	}
	defer r.Close()

	b, err := io.ReadAll(r)
	if err != nil {
		return err

	}

	startAt := bytes.IndexByte(b, '\x00') + 1
	fmt.Print(string(b[startAt:]))
	return nil
}
