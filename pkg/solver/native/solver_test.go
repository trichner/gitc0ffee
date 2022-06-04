package native

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/trichner/gitc0ffee/pkg/commit"
	"github.com/trichner/gitc0ffee/pkg/solver"
	"github.com/trichner/gitc0ffee/pkg/solver/model"
	"testing"
)

const rawHeaderAndBodyObject = `tree e57181f20b062532907436169bb5823b6af2f099
author Thomas Richner <thomas.richner@oviva.com> 1653693519 +0200
committer Thomas Richner <thomas.richner@oviva.com> 1653693519 +0200

Initial commit
36abde0100000000`

func TestNativeSolver_Solve(t *testing.T) {

	tpl, err := getTemplate()
	assert.NoError(t, err)

	s := &nativeSolver{}

	pfx := []byte{0x88, 0x70}

	//when
	obj, err := s.Solve(tpl, pfx)

	//then
	assert.NoError(t, err)
	assert.Equal(t, []byte("8870"), obj.Hash[:4])
}

func getTemplate() (*model.ObjectTemplate, error) {

	c, err := commit.ParseGitCommitObject([]byte(rawHeaderAndBodyObject))
	if err != nil {
		return nil, err
	}

	tpl, err := solver.PrepareTemplate(c)
	if err != nil {
		return nil, err
	}

	return tpl, nil
}

func BenchmarkSingleThreaded_Solve(b *testing.B) {

	c, err := commit.ParseGitCommitObject([]byte(rawHeaderAndBodyObject))
	assert.NoError(b, err)

	tpl, err := solver.PrepareTemplate(c)
	assert.NoError(b, err)

	s := NewFactory().NewSolver(0, uint64(b.N))

	b.ResetTimer()
	_, err = s.Solve(tpl, []byte{0xc0, 0xff, 0xee, 0xba, 0xdc, 0x0d})
	if err != nil && !errors.Is(err, model.ErrExhaustedSalts) {
		b.Fatal(err)
	}
}
