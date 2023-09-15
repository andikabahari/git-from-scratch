package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

// Usage: your_git.sh <command> <arg1> <arg2> ...
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "init":
		for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			}
		}

		headFileContents := []byte("ref: refs/heads/master\n")
		if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		}

		fmt.Println("Initialized git directory")

	case "cat-file":
		sha := os.Args[3]
		f, err := os.Open(path.Join(".git/objects", sha[:2], sha[2:]))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file: %s\n", err)
		}
		defer f.Close()

		r, err := zlib.NewReader(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating zlib reader: %s\n", err)
		}
		defer r.Close()

		b, err := io.ReadAll(r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
		}

		startAt := bytes.IndexByte(b, '\x00') + 1
		fmt.Print(string(b[startAt:]))

	case "hash-object":
		filepath := os.Args[3]
		checksum, err := writeBlob(filepath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing blob object: %s\n", err)
		}
		fmt.Printf("%x\n", checksum)

	case "ls-tree":
		sha := os.Args[3]
		treepath := path.Join(".git/objects", sha[:2], sha[2:])
		f, err := os.Open(treepath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file: %s\n", err)
		}
		defer f.Close()

		r, err := zlib.NewReader(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating zlib reader: %s\n", err)
		}
		defer r.Close()

		b, err := io.ReadAll(r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
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

	case "write-tree":
		wd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current directory: %s\n", err)
		}

		checksum, err := writeTree(wd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing tree object: %s\n", err)
		}
		fmt.Printf("%x\n", checksum)

	case "commit-tree":
		treeSha := os.Args[2]
		commitSha := os.Args[4]
		msg := os.Args[6]
		checksum, err := writeCommit(treeSha, commitSha, msg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing commit object: %s\n", err)
		}
		fmt.Printf("%x\n", checksum)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}

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
