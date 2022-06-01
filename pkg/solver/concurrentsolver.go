package solver

import (
	"errors"
	"github.com/trichner/gitc0ffee/pkg/digest"
	"github.com/trichner/gitc0ffee/pkg/template"
	"math"
	"runtime"
)

const chunkSize = 4096

type solution struct {
	Digest *digest.HexObjectDigest
	Bytes  []byte
}

type concurrentSolver struct {
}

func (c *concurrentSolver) Solve(obj *template.ObjectTemplate, prefix []byte) (*digest.HexObjectDigest, []byte, error) {

	numWorkers := runtime.NumCPU()
	tasksChan := make(chan *singleThreaded)
	solutionChan := make(chan *solution)

	// create jobs
	go func() {
	}()

	// work on jobs
	for i := 0; i < numWorkers; i++ {
		t := obj.Copy()
		go func() {
			for job := range tasksChan {
				hexDigest, data, err := job.Solve(t, prefix)
				if errors.Is(err, ErrExhaustedSalts) {
					continue
				}
				if err != nil {
					panic(err)
				}
				solutionChan <- &solution{
					Digest: hexDigest,
					Bytes:  data,
				}
			}
		}()
	}

	for {
		start := uint64(0)
		end := uint64(chunkSize)
		for {
			task := &singleThreaded{
				saltStart: start,
				saltEnd:   end,
			}

			select {
			case tasksChan <- task:
			case s := <-solutionChan:
				close(tasksChan)
				// should we wait fo the others to finish?
				return s.Digest, s.Bytes, nil
			}

			if end == math.MaxUint64 {
				// we're done, exhausted all options
				break
			}

			start = end
			var ok bool
			end, ok = addU64(end, chunkSize)
			if !ok {
				end = math.MaxUint64
			}
		}
	}
}

func addU64(left, right uint64) (uint64, bool) {
	if left > math.MaxUint64-right {

	}
	return left + right, true
}
