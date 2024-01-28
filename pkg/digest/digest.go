package digest

import (
	"encoding/hex"
	"fmt"
)

const hextable = "0123456789abcdef"

type (
	ObjectDigest    [20]byte
	HexObjectDigest [40]byte
)

func HexEncodeDigest(dst *HexObjectDigest, src *ObjectDigest) int {
	j := 0
	for i := 0; i < len(src); i++ {
		dst[j] = hextable[src[i]>>4]
		dst[j+1] = hextable[src[i]&0x0f]
		j += 2
	}
	return len(src) * 2
}

func HexDecodeDigest(dst *ObjectDigest, src *HexObjectDigest) error {
	_, err := hex.Decode(dst[:], src[:])
	if err != nil {
		return fmt.Errorf("cannot decode hex string %q", src)
	}
	return nil
}

func (o *ObjectDigest) String() string {
	return hex.EncodeToString(o[:])
}
