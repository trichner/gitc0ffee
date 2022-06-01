package main

import (
	"bytes"
	"fmt"
	"github.com/trichner/gitc0ffee/pkg/digest"
	"os"
	"os/exec"
)

func getCommitContents(digest *digest.HexObjectDigest) ([]byte, error) {

	out, err := runCommand("git", "cat-file", "-p", string(digest[:]))
	if err != nil {
		return nil, fmt.Errorf("cannot get contents of revision %q: %w", digest, err)
	}
	return out, nil
}

func getHeadDigest() (*digest.HexObjectDigest, error) {
	out, err := runCommand("git", "rev-parse", "HEAD")
	if err != nil {
		return nil, fmt.Errorf("cannot read HEAD rev: %w", err)
	}

	hexDigestBytes := bytes.TrimSpace(out)
	if err != nil {
		return nil, fmt.Errorf("cannot parse HEAD rev %q: %w", hexDigestBytes, err)
	}
	if len(hexDigestBytes) != 40 {
		return nil, fmt.Errorf("digest length not matching 40 != %d", len(hexDigestBytes))
	}
	var hexDigest digest.HexObjectDigest
	copy(hexDigest[:], hexDigestBytes)
	return &hexDigest, nil
}

func runCommand(prog string, args ...string) ([]byte, error) {

	cmd := exec.Command(prog, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
