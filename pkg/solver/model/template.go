package model

import (
	"crypto/sha1"

	"github.com/trichner/gitc0ffee/pkg/digest"
)

const hextable = "0123456789abcdef"

type ObjectTemplate struct {
	Bytes         []byte
	PayloadOffset int
	SaltOffset    int
}

func (t *ObjectTemplate) SetSalt(salt uint64) {
	hexEncodeUint64(t.Bytes[t.SaltOffset:], salt)
}

func (t *ObjectTemplate) Sum() *digest.ObjectDigest {
	d := digest.ObjectDigest(sha1.Sum(t.Bytes))
	return &d
}

func (t *ObjectTemplate) Payload() []byte {
	return t.Bytes[t.PayloadOffset:]
}

func (t *ObjectTemplate) Copy() *ObjectTemplate {
	newBytes := make([]byte, len(t.Bytes))

	copy(newBytes, t.Bytes)
	return &ObjectTemplate{
		Bytes:         newBytes,
		SaltOffset:    t.SaltOffset,
		PayloadOffset: t.PayloadOffset,
	}
}

func hexEncodeUint64(dst []byte, src uint64) {
	// dst must be at least 16 bytes long

	dst[15] = hextable[src&0x0f]
	dst[14] = hextable[(src>>4)&0x0f]

	dst[13] = hextable[(src>>8)&0x0f]
	dst[12] = hextable[(src>>12)&0x0f]

	dst[11] = hextable[(src>>16)&0x0f]
	dst[10] = hextable[(src>>20)&0x0f]

	dst[9] = hextable[(src>>24)&0x0f]
	dst[8] = hextable[(src>>28)&0x0f]

	dst[7] = hextable[(src>>32)&0x0f]
	dst[6] = hextable[(src>>36)&0x0f]

	dst[5] = hextable[(src>>40)&0x0f]
	dst[4] = hextable[(src>>44)&0x0f]

	dst[3] = hextable[(src>>48)&0x0f]
	dst[2] = hextable[(src>>52)&0x0f]

	dst[1] = hextable[(src>>56)&0x0f]
	dst[0] = hextable[(src>>60)&0x0f]
}
