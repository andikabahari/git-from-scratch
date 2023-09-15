package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path"
)

func lsTree(sha string) error {
	treepath := path.Join(".git/objects", sha[:2], sha[2:])
	f, err := os.Open(treepath)
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

	// skip header
	b = b[bytes.IndexByte(b, '\x00')+1:]

	for len(b) > 0 {
		mode := b[:bytes.IndexByte(b, ' ')]
		b = b[len(mode)+1:]

		name := b[:bytes.IndexByte(b, '\x00')]
		b = b[len(name)+1:]

		sha := b[:20]
		b = b[len(sha):]

		fmt.Println(string(name))
	}

	return nil
}
