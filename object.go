package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"github.com/trichner/gitc0ffee/pkg/digest"
	"io"
	"os"
	"path/filepath"
)

const gitObjectBasePath = ".git/objects"

func writeGitObject(contents []byte) error {

	digestBytes := hashObject(contents)

	var hexDigest digest.HexObjectDigest
	digest.HexEncodeDigest(&hexDigest, digestBytes)

	f, err := OpenGitObjectFile(&hexDigest)
	if err != nil {
		return fmt.Errorf("cannot open git object %q: %w", digestBytes, err)
	}
	defer f.Close()

	writer := zlib.NewWriter(f)
	_, err = writer.Write(contents)
	if err != nil {
		return fmt.Errorf("cannot compress git object %q: %w", digestBytes, err)
	}
	defer writer.Close()

	return nil
}

func OpenGitObjectFile(digest *digest.HexObjectDigest) (io.WriteCloser, error) {

	path := getGitObjectPath(digest)

	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, os.FileMode(0755))
	if err != nil {
		return nil, fmt.Errorf("cannot create git object folder %q: %w", path, err)
	}
	f, err := os.OpenFile(path, os.O_TRUNC|os.O_CREATE|os.O_RDWR, os.FileMode(0666))
	return f, err
}

func hashObject(contents []byte) *digest.ObjectDigest {
	byteDigest := digest.ObjectDigest(sha1.Sum(contents))
	return &byteDigest
}

// https://matthew-brett.github.io/curious-git/reading_git_objects.html
func readGitObject(hexDigest *digest.HexObjectDigest) ([]byte, error) {

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

func getGitObjectPath(digest *digest.HexObjectDigest) string {
	return fmt.Sprintf("%s/%s/%s", gitObjectBasePath, digest[:2], digest[2:])
}
