package solver

import (
	"math"

	"github.com/trichner/gitc0ffee/pkg/solver/model"
	"github.com/trichner/gitc0ffee/pkg/solver/native"
)

func NewSingleThreadedSolver() model.DigestPrefixSolver {
	return &singleThreaded{
		saltStart: 0,
		saltEnd:   math.MaxUint64,
	}
}

type singleThreadedSolverFactory struct{}

func (s *singleThreadedSolverFactory) NewSolver(startSalt, endSalt uint64) model.DigestPrefixSolver {
	return &singleThreaded{
		saltStart: startSalt,
		saltEnd:   endSalt,
	}
}

func NewConcurrentSolver() model.DigestPrefixSolver {
	return DecorateWithLogging(&concurrentSolver{
		solverFactory: &singleThreadedSolverFactory{},
	})
}

func NewNativeConcurrentSolver() model.DigestPrefixSolver {
	return DecorateWithLogging(&concurrentSolver{
		solverFactory: native.NewFactory(),
	})
}
