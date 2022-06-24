package vcl

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseTable(t *testing.T) {
	expected := []uint16{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var gotBytes []byte

	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	for i := range expected {
		_ = binary.Write(w, binary.LittleEndian, expected[i])
		_ = w.Flush()
		buf := []byte{0, 0}
		b.Read(buf)
		b.Reset()
		gotBytes = append(gotBytes, buf...)

	}

	r := bytes.NewReader(gotBytes)

	got, err := parseTable[uint16](r, 0, 100)

	assert.NoError(t, err)
	exp, err := got.All()
	assert.NoError(t, err)
	assert.Equal(t, expected, exp)

}
