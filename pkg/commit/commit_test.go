package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

const commit = `commit 213\x00tree e57181f20b062532907436169bb5823b6af2f099
author Thomas Richner <thomas.richner@oviva.com> 1653693519 +0200
committer Thomas Richner <thomas.richner@oviva.com> 1653693519 +0200

Initial commit
36abde0100000000`

func TestParseGitCommitObjectPrefix_simple(t *testing.T) {
	rawCommit := []byte(strings.Replace(commit, "\\x00", "\x00", -1))
	prefix, err := parseGitCommitObjectPrefix(rawCommit)
	assert.NoError(t, err)

	assert.Equal(t, []byte("commit 213"), prefix)
}

func TestParseGitCommitObjectPrefix_badType(t *testing.T) {
	rawCommit := []byte("tree 3\x00123")
	_, err := parseGitCommitObjectPrefix(rawCommit)
	assert.Error(t, err)
}

func TestParseGitCommitObjectPrefix_badPrefix(t *testing.T) {
	rawCommit := []byte("commit 123 this is crap")
	_, err := parseGitCommitObjectPrefix(rawCommit)
	assert.Error(t, err)
}

func TestParseGitCommitObjectPrefix_badLength(t *testing.T) {
	rawCommit := []byte("commit -23\x00abc")
	_, err := parseGitCommitObjectPrefix(rawCommit)
	assert.Error(t, err)
}
