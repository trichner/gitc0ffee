package git

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/trichner/gitc0ffee/pkg/digest"
)

func GetCommitContents(digest *digest.HexObjectDigest) ([]byte, error) {
	out, err := runCommand("git", "cat-file", "-p", string(digest[:]))
	if err != nil {
		return nil, fmt.Errorf("cannot get contents of revision %q: %w", digest, err)
	}
	return out, nil
}

func GetHeadDigest() (*digest.HexObjectDigest, error) {
	out, err := runCommand("git", "rev-parse", "HEAD")
	if err != nil {
		return nil, fmt.Errorf("cannot read HEAD rev: %w", err)
	}

	hexDigestBytes := bytes.TrimSpace(out)
	if len(hexDigestBytes) != 40 {
		return nil, fmt.Errorf("digest length not matching 40 != %d", len(hexDigestBytes))
	}
	var hexDigest digest.HexObjectDigest
	copy(hexDigest[:], hexDigestBytes)
	return &hexDigest, nil
}

func UpdateReference(ref string, hash string) error {
	_, err := runCommand("git", "update-ref", ref, hash)
	if err != nil {
		return fmt.Errorf("cannot update reference %q to object %q: %w", ref, hash, err)
	}

	return nil
}

func GetCurrentBranch() (string, error) {
	out, err := runCommand("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", fmt.Errorf("cannot read current branch: %w", err)
	}

	currentBranch := strings.TrimSpace(string(out))
	if currentBranch == "" {
		return "", fmt.Errorf("current branch undefined")
	}

	return currentBranch, nil
}

func WriteObject(t string, contents []byte) (*digest.HexObjectDigest, error) {
	buf := bytes.NewBuffer(contents)
	out, err := runCommandWithStdin(buf, "git", "hash-object", "-w", "-t", t, "--stdin")
	if err != nil {
		return nil, err
	}
	out = bytes.TrimSpace(out)
	return (*digest.HexObjectDigest)(out), nil
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

func runCommandWithStdin(stdin io.Reader, prog string, args ...string) ([]byte, error) {
	cmd := exec.Command(prog, args...)
	var out bytes.Buffer
	cmd.Stdin = stdin
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
