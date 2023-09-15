package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

func writeObject(objectType string, content []byte) ([20]byte, error) {
	header := fmt.Sprintf("%s %d", objectType, len(content))
	store := []byte(header)
	store = append(store, '\x00')
	store = append(store, content...)

	checksum := sha1.Sum(store)
	sumstr := fmt.Sprintf("%x", checksum)

	dirpath := path.Join(".git/objects", sumstr[:2])
	err := os.Mkdir(dirpath, 0755)
	if err != nil {
		return [20]byte{}, err
	}

	objpath := path.Join(dirpath, sumstr[2:])
	f, err := os.Create(objpath)
	if err != nil {
		return [20]byte{}, err
	}
	defer f.Close()

	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(store)
	w.Close()
	f.Write(b.Bytes())

	return checksum, nil
}

func writeBlob(filepath string) ([20]byte, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return [20]byte{}, err
	}

	return writeObject("blob", content)
}

func writeTree(rootpath string) ([20]byte, error) {
	entries, err := os.ReadDir(rootpath)
	if err != nil {
		return [20]byte{}, err
	}

	var b bytes.Buffer
	for _, entry := range entries {
		if entry.Name() == ".git" {
			continue
		}

		var mode int
		var checksum [20]byte
		var err error
		objpath := filepath.Join(rootpath, entry.Name())
		if entry.IsDir() {
			mode = 0o040000
			checksum, err = writeTree(objpath)
		} else {
			mode = 0o100644
			checksum, err = writeBlob(objpath)
		}
		if err != nil {
			return [20]byte{}, err
		}

		s := fmt.Sprintf("%o %s\x00%s", mode, entry.Name(), checksum)
		b.WriteString(s)
	}

	return writeObject("tree", b.Bytes())
}

func writeCommit(treeSha, commitSha, msg string) ([20]byte, error) {
	content := fmt.Sprintf("tree %s\n", treeSha)
	content += fmt.Sprintf("parent %s\n", commitSha)
	content += fmt.Sprintf("author %s\n", "user@example.com")
	content += fmt.Sprintf("committer %s\n\n", "user@example.com")
	content += fmt.Sprintf("%s\n", msg)
	return writeObject("commit", []byte(content))
}
