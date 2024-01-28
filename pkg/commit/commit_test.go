package commit

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const rawHeaderAndBodyObject = `tree e57181f20b062532907436169bb5823b6af2f099
author Thomas Richner <thomas.richner@oviva.com> 1653693519 +0200
committer Thomas Richner <thomas.richner@oviva.com> 1653693519 +0200

Initial commit
36abde0100000000`

func TestParseHeaders(t *testing.T) {
	raw := []byte(rawHeaderAndBodyObject)
	buf := bytes.NewBuffer(raw)
	headers, err := parseHeaders(buf)
	assert.NoError(t, err)

	assert.Len(t, headers, 3)
	assert.True(t, strings.HasPrefix(headers[0].Value, "tree "))
	assert.True(t, strings.HasPrefix(headers[1].Value, "author "))
	assert.True(t, strings.HasPrefix(headers[2].Value, "committer "))

	for _, h := range headers {
		fmt.Println(h.Value)
	}
}
