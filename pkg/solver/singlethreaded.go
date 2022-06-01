package solver

import (
	"github.com/trichner/gitc0ffee/pkg/digest"
	"github.com/trichner/gitc0ffee/pkg/template"
	"log"
	"time"
)

type singleThreaded struct {
	saltStart, saltEnd uint64
}

func (s *singleThreaded) Solve(obj *template.ObjectTemplate, prefix []byte) (*digest.HexObjectDigest, []byte, error) {

	tick := time.Now()
	for salt := s.saltStart; salt < s.saltEnd; salt++ {
		obj.SetSalt(salt)
		digestBytes := obj.Sum()

		//if digestBytes[0] == 0xc0 && digestBytes[1] == 0xff && digestBytes[2] == 0xee {
		if hasPrefix(digestBytes, prefix) {
			//if digest[0] == 0xc0 && digest[1] == 0xfe && digest[2] == 0xba && digest[3] == 0xbe {
			tock := time.Now()
			d := tock.Sub(tick)
			rate := float64(salt-s.saltStart) / d.Seconds() / 1000
			log.Printf("found in %4.f seconds at %4.f khash/s", d.Seconds(), rate)

			var hexDigest digest.HexObjectDigest
			digest.HexEncodeDigest(&hexDigest, digestBytes)
			return &hexDigest, obj.Bytes, nil
		}
		if salt&0xffffff == 0xffffff {
			tock := time.Now()
			d := tock.Sub(tick)
			rate := float64(salt-s.saltStart) / d.Seconds() / 1000
			log.Printf("brute forcing at %4.f khash/s", rate)
		}
	}

	return nil, nil, ErrExhaustedSalts
}

func hasPrefix(s *digest.ObjectDigest, prefix []byte) bool {
	var sum byte
	for i := 0; i < len(prefix); i++ {
		sum = sum | (s[i] ^ prefix[i])
	}
	return sum == 0
}
