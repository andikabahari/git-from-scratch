package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path"
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
		content, err := os.ReadFile(filepath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
		}

		header := fmt.Sprintf("blob %d", len(content))
		store := []byte(header)
		store = append(store, '\x00')
		store = append(store, content...)

		h := sha1.New()
		h.Write(store)
		checksum := fmt.Sprintf("%x", h.Sum(nil))

		dirpath := path.Join(".git/objects", checksum[:2])
		err = os.Mkdir(dirpath, 0755)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error making directory: %s\n", err)
		}

		objpath := path.Join(dirpath, checksum[2:])
		f, err := os.Create(objpath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating file: %s\n", err)
		}
		defer f.Close()

		var b bytes.Buffer
		w := zlib.NewWriter(&b)
		w.Write(store)
		w.Close()
		f.Write(b.Bytes())

		fmt.Print(checksum)

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

		var names []string
		for len(b) > 0 {
			mode := b[:bytes.IndexByte(b, ' ')]
			b = b[len(mode)+1:]

			name := b[:bytes.IndexByte(b, '\x00')]
			b = b[len(name)+1:]

			sha := b[:20]
			b = b[len(sha):]

			names = append(names, string(name))
		}

		for _, name := range names {
			fmt.Println(name)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
