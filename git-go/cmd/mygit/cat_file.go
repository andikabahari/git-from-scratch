package main

import (
	"bytes"
	"compress/zlib"
	"io"
	"os"
	"path"
)

func catFile(sha string) ([]byte, error) {
	f, err := os.Open(path.Join(".git/objects", sha[:2], sha[2:]))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r, err := zlib.NewReader(f)
	if err != nil {
		return nil, err

	}
	defer r.Close()

	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err

	}

	startAt := bytes.IndexByte(b, '\x00') + 1
	return b[startAt:], nil
}
