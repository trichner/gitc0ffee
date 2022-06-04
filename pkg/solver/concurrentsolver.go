package solver

import (
	"errors"
	"github.com/trichner/gitc0ffee/pkg/solver/model"
	"math"
	"runtime"
)

const chunkSize = 4096

type concurrentSolver struct {
	solverFactory model.SolverFactory
}

func (c *concurrentSolver) Solve(obj *model.ObjectTemplate, prefix []byte) (*model.CommitObject, error) {

	numWorkers := runtime.NumCPU()
	tasksChan := make(chan model.DigestPrefixSolver)
	solutionChan := make(chan *model.CommitObject)

	// workers
	for i := 0; i < numWorkers; i++ {
		t := obj.Copy()
		go func() {
			for job := range tasksChan {
				res, err := job.Solve(t, prefix)
				if errors.Is(err, model.ErrExhaustedSalts) {
					continue
				}
				if err != nil {
					panic(err)
				}
				solutionChan <- res
			}
		}()
	}

	// creating jobs
	for {
		start := uint64(0)
		end := uint64(chunkSize)
		for {
			task := c.solverFactory.NewSolver(start, end)

			select {
			case tasksChan <- task:
			case solution := <-solutionChan:
				close(tasksChan)
				//FIXME: if we exhaust all options and don't find a solution we're stuck
				return solution, nil
			}

			if end == math.MaxUint64 {
				// we're done, exhausted all options
				break
			}

			start = end
			var ok bool
			end, ok = safeAddU64(end, chunkSize)
			if !ok {
				end = math.MaxUint64
			}
		}
	}
}

func safeAddU64(left, right uint64) (uint64, bool) {
	if left > math.MaxUint64-right {

	}
	return left + right, true
}
