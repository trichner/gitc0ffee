package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func getCommitContents(digest *hexObjectDigest) (string, error) {

	out, err := runCommand("git", "cat-file", "-p", string(digest[:]))
	if err != nil {
		return "", fmt.Errorf("cannot get contents of revision %q: %w", digest, err)
	}
	return out, nil
}

func getHeadDigest() (*hexObjectDigest, error) {
	out, err := runCommand("git", "rev-parse", "HEAD")
	if err != nil {
		return nil, fmt.Errorf("cannot read HEAD rev: %w", err)
	}

	digest := strings.TrimSpace(out)
	if err != nil {
		return nil, fmt.Errorf("cannot parse HEAD rev %q: %w", digest, err)
	}
	if len(digest) != 40 {
		return nil, fmt.Errorf("digest length not matching 40 != %d", len(digest))
	}
	var hexDigest hexObjectDigest
	copy(hexDigest[:], digest)
	return &hexDigest, nil
}

func runCommand(prog string, args ...string) (string, error) {

	cmd := exec.Command(prog, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		return "", err
	}
	return out.String(), nil
}
