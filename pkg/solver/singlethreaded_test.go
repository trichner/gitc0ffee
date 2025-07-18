package solver

import (
	"crypto/sha1"
	"errors"
	"testing"

	"github.com/trichner/gitc0ffee/pkg/solver/util"

	"github.com/trichner/gitc0ffee/pkg/assert"

	"github.com/trichner/gitc0ffee/pkg/commit"
	"github.com/trichner/gitc0ffee/pkg/solver/model"
)

const rawHeaderAndBodyObject = `tree e57181f20b062532907436169bb5823b6af2f099
author Thomas Richner <thomas.richner@oviva.com> 1653693519 +0200
committer Thomas Richner <thomas.richner@oviva.com> 1653693519 +0200

Initial commit
36abde0100000000`

func BenchmarkSingleThreaded_Solve(b *testing.B) {
	c, err := commit.ParseGitCommitObject([]byte(rawHeaderAndBodyObject))
	assert.NoError(b, err)

	tpl, err := util.PrepareTemplate(c)
	assert.NoError(b, err)

	s := &singleThreaded{
		saltStart: 0,
		saltEnd:   uint64(b.N),
	}

	b.ResetTimer()
	_, err = s.Solve(tpl, []byte{0xc0, 0xff, 0xee, 0xba, 0xdc, 0x0d})
	if err != nil && !errors.Is(err, model.ErrExhaustedSalts) {
		b.Fatal(err)
	}
}

func BenchmarkSha1_Sum(b *testing.B) {
	data, err := getBenchBytes()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sha1.Sum(data)
	}
}

func getBenchBytes() ([]byte, error) {
	c, err := commit.ParseGitCommitObject([]byte(rawHeaderAndBodyObject))
	if err != nil {
		return nil, err
	}

	tpl, err := util.PrepareTemplate(c)
	if err != nil {
		return nil, err
	}

	return tpl.Bytes, nil
}
