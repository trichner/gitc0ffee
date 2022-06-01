package solver

import (
	"fmt"
	"github.com/trichner/gitc0ffee/pkg/digest"
	"github.com/trichner/gitc0ffee/pkg/template"
	"math"
)

var ErrExhaustedSalts = fmt.Errorf("exhausted possible salts without finding a solution")

type DigestPrefixSolver interface {
	//Solve finds a valid permutation of the ObjectTemplate for which the digest matches the given prefix
	Solve(obj *template.ObjectTemplate, prefix []byte) (*digest.HexObjectDigest, []byte, error)
}

func NewSingleThreadedSolver() DigestPrefixSolver {
	return &singleThreaded{
		saltStart: 0,
		saltEnd:   math.MaxUint64,
	}
}

func NewConcurrentSolver() DigestPrefixSolver {
	return &concurrentSolver{}
}
