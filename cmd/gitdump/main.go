package main

import (
	"bytes"
	"compress/zlib"
	"flag"
	"fmt"
	"io"
	"os"
)

const gitObjectBasePath = ".git/objects"

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) <= 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s <objectdigest>\n", os.Args[0])
		os.Exit(1)
	}
	digest := args[0]

	contents, err := readGitObject([]byte(digest))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(string(contents))
}

// https://matthew-brett.github.io/curious-git/reading_git_objects.html
func readGitObject(hexDigest []byte) ([]byte, error) {
	p := getGitObjectPath(hexDigest)

	data, err := os.ReadFile(p)
	if err != nil {
		return nil, fmt.Errorf("cannot read object %q: %w", hexDigest, err)
	}

	r, err := zlib.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("cannot decompress object %q: %w", hexDigest, err)
	}
	defer r.Close()

	contents, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("cannot decompress object %q: %w", hexDigest, err)
	}
	return contents, nil
}

func getGitObjectPath(digest []byte) string {
	return fmt.Sprintf("%s/%s/%s", gitObjectBasePath, digest[:2], digest[2:])
}
