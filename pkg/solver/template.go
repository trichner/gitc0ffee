package solver

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/trichner/gitc0ffee/pkg/commit"
	"github.com/trichner/gitc0ffee/pkg/solver/model"
)

const saltHeaderName = "coffeesalt"

func PrepareTemplate(commitObject *commit.Object) (*model.ObjectTemplate, error) {
	var buf bytes.Buffer

	saltHeaderPrefix := saltHeaderName + " "
	saltValue := hex.EncodeToString(bytes.Repeat([]byte{0}, 8))

	headers := commitObject.Headers()
	for _, h := range headers {
		if !strings.HasPrefix(h.Value, saltHeaderPrefix) {
			buf.WriteString(h.Value)
			buf.WriteString("\n")
		}
	}

	// append salt header, this may lead to duplicates though we don't really care
	buf.WriteString(saltHeaderPrefix)
	saltOffset := buf.Len()
	buf.WriteString(saltValue)
	buf.WriteString("\n\n")

	buf.Write(commitObject.Message)

	objectPayload := buf.Bytes()

	objectPrefix := fmt.Sprintf("commit %d\x00", len(objectPayload))
	payloadOffest := len(objectPrefix)

	var objectBuf bytes.Buffer
	objectBuf.WriteString(objectPrefix)
	prefixLength := objectBuf.Len()
	_, err := buf.WriteTo(&objectBuf)
	if err != nil {
		return nil, err
	}

	return &model.ObjectTemplate{
		Bytes:         objectBuf.Bytes(),
		SaltOffset:    prefixLength + saltOffset,
		PayloadOffset: payloadOffest,
	}, nil
}
