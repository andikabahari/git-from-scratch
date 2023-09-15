package main

import (
	"fmt"
	"os"
)

// Usage: your_git.sh <command> <arg1> <arg2> ...
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "init":
		err := initGit()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing git: %s\n", err)
			os.Exit(1)
		}
		fmt.Println("Initialized git directory")

	case "cat-file":
		sha := os.Args[3]
		b, err := catFile(sha)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file content: %s\n", err)
			os.Exit(1)
		}
		fmt.Print(string(b))

	case "hash-object":
		filepath := os.Args[3]
		checksum, err := writeBlob(filepath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing blob object: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("%x\n", checksum)

	case "ls-tree":
		sha := os.Args[3]
		names, err := lsTree(sha)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing tree: %s\n", err)
			os.Exit(1)
		}
		for _, b := range names {
			fmt.Println(string(b))
		}

	case "write-tree":
		wd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current directory: %s\n", err)
			os.Exit(1)
		}

		checksum, err := writeTree(wd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing tree object: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("%x\n", checksum)

	case "commit-tree":
		treeSha := os.Args[2]
		commitSha := os.Args[4]
		msg := os.Args[6]
		checksum, err := writeCommit(treeSha, commitSha, msg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing commit object: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("%x\n", checksum)

	case "clone":
		url := os.Args[2]
		dir := os.Args[3]
		err := clone(url, dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error cloning repository %s\n", err)
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
