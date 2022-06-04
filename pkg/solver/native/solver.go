package native

/*
#include <stdio.h>
#include <stdint.h>
#include "solver.h"
*/
import "C"

import (
	"github.com/trichner/gitc0ffee/pkg/digest"
	"github.com/trichner/gitc0ffee/pkg/solver/model"
)

type nativeSolver struct {
	saltStart, saltEnd uint64
}

type carray struct {
	arr *C.uint8_t
	len C.size_t
}

type nativeSolverFactory struct {
}

func (n *nativeSolverFactory) NewSolver(startSalt, endSalt uint64) model.DigestPrefixSolver {
	return &nativeSolver{
		saltStart: startSalt,
		saltEnd:   endSalt,
	}
}

func NewFactory() model.SolverFactory {
	return &nativeSolverFactory{}
}

func (n *nativeSolver) Solve(obj *model.ObjectTemplate, prefix []byte) (*model.CommitObject, error) {

	rawBytes := toCArray(obj.Bytes)
	pfxBytes := toCArray(prefix)

	saltOffset := C.size_t(obj.SaltOffset)
	saltStart := C.uint64_t(n.saltStart)
	saltEnd := C.uint64_t(n.saltEnd)

	// hand over to native C
	csalt := C.solve(rawBytes.arr, rawBytes.len, pfxBytes.arr, pfxBytes.len, saltOffset, saltStart, saltEnd)
	if csalt == C.ERR_SALTS_EXHAUSTED {
		return nil, model.ErrExhaustedSalts
	}

	var hash digest.HexObjectDigest
	digest.HexEncodeDigest(&hash, obj.Sum())
	return &model.CommitObject{
		Raw:     obj.Bytes,
		Payload: obj.Payload(),
		Hash:    &hash,
	}, nil
}

func toCArray(s []byte) *carray {
	rawBytesPtr := (*C.uint8_t)(&s[0])
	rawBytesLen := (C.size_t)(len(s))
	return &carray{
		arr: rawBytesPtr,
		len: rawBytesLen,
	}
}
