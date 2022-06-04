package solver

import (
	"github.com/trichner/gitc0ffee/pkg/solver/model"
	"log"
	"time"
)

type timedSolver struct {
	solver model.DigestPrefixSolver
}

func DecorateWithLogging(s model.DigestPrefixSolver) model.DigestPrefixSolver {
	return &timedSolver{solver: s}
}

func (t *timedSolver) Solve(obj *model.ObjectTemplate, prefix []byte) (*model.CommitObject, error) {

	tick := time.Now()
	o, err := t.solver.Solve(obj, prefix)
	tock := time.Now()

	d := tock.Sub(tick)
	log.Printf("found solution in %.2fs", d.Seconds())
	return o, err
}
