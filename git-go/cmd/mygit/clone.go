package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

func clone(url, dir string) error {
	err := os.Mkdir(dir, 0755)
	if err != nil {
		return err
	}

	err = os.Chdir(dir)
	if err != nil {
		return err
	}

	err = initGit()
	if err != nil {
		return err
	}

	refs, err := findRefs(url)
	if err != nil {
		return err
	}

	b, err := wantRefs(url, refs)
	if err != nil {
		return err
	}

	p, err := parsePack(b)
	if err != nil {
		return err
	}

	return unpack(dir, p)
}

type ref struct {
	Sha  []byte
	Name []byte
}

func findRefs(url string) ([]ref, error) {
	url += "/info/refs?service=git-upload-pack"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// skip twice
	b = b[bytes.IndexByte(b, '\n')+1:]
	b = b[bytes.IndexByte(b, '\n')+1:]

	// remove 0000
	end := []byte("0000")
	b = b[:len(b)-len(end)]

	refs := make([]ref, 0)
	for len(b) > 0 {
		// skip first 4 bytes
		b = b[4:]

		sha := b[:bytes.IndexByte(b, ' ')]
		b = b[len(sha)+1:]

		name := b[:bytes.IndexByte(b, '\n')]
		b = b[len(name)+1:]

		refs = append(refs, ref{
			Sha:  sha,
			Name: name,
		})
	}

	return refs, nil
}

func wantRefs(url string, refs []ref) ([]byte, error) {
	url += "/git-upload-pack"

	var body bytes.Buffer
	for _, ref := range refs {
		want := fmt.Sprintf("0032want %s\n", string(ref.Sha))
		body.WriteString(want)
	}
	body.WriteString("0000")
	body.WriteString("0009done\n")

	resp, err := http.Post(url, "application/x-git-upload-pack-request", &body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	startAt := bytes.IndexByte(b, '\n') + 1
	return b[startAt:], nil
}
