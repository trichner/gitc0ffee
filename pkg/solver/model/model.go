package model

import (
	"fmt"

	"github.com/trichner/gitc0ffee/pkg/digest"
)

var ErrExhaustedSalts = fmt.Errorf("exhausted possible salts without finding a solution")

type CommitObject struct {
	Raw     []byte
	Payload []byte
	Hash    *digest.HexObjectDigest
}

type DigestPrefixSolver interface {
	// Solve finds a valid permutation of the ObjectTemplate for which the digest matches the given prefix
	Solve(obj *ObjectTemplate, prefix []byte) (*CommitObject, error)
}

type SolverFactory interface {
	NewSolver(startSalt, endSalt uint64) DigestPrefixSolver
}
