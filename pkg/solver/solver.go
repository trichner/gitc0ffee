package solver

import (
	"fmt"
	"github.com/trichner/gitc0ffee/pkg/digest"
	"math"
)

var ErrExhaustedSalts = fmt.Errorf("exhausted possible salts without finding a solution")

type CommitObject struct {
	Raw     []byte
	Payload []byte
	Hash    *digest.HexObjectDigest
}

type DigestPrefixSolver interface {
	//Solve finds a valid permutation of the ObjectTemplate for which the digest matches the given prefix
	Solve(obj *ObjectTemplate, prefix []byte) (*CommitObject, error)
}

func NewSingleThreadedSolver() DigestPrefixSolver {
	return &singleThreaded{
		saltStart: 0,
		saltEnd:   math.MaxUint64,
	}
}

type singleThreadedSolverFactory struct {
}

func (s *singleThreadedSolverFactory) NewSolver(startSalt, endSalt uint64) DigestPrefixSolver {
	return &singleThreaded{
		saltStart: startSalt,
		saltEnd:   endSalt,
	}
}

func NewConcurrentSolver() DigestPrefixSolver {
	return &concurrentSolver{
		solverFactory: &singleThreadedSolverFactory{},
	}
}
