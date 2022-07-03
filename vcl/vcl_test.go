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

	gotBytes := sliceToBytes[uint16](expected)
	r := bytes.NewReader(gotBytes)

	got, err := parseTable[uint16](r, 0, 100)

	assert.NoError(t, err)
	exp, err := got.All()
	assert.NoError(t, err)
	assert.Equal(t, expected, exp)
}

func TestParseTables(t *testing.T) {
	offsetTable := make([]byte, 200)
	lengthTable := make([]byte, 100)
	freqTable := make([]byte, 100)

	for i := range offsetTable {
		offsetTable[i] = 1
	}

	for i := range lengthTable {
		lengthTable[i] = 2
	}

	for i := range freqTable {
		freqTable[i] = 3
	}

	tables := append(offsetTable, lengthTable...)
	tables = append(tables, freqTable...)

	r := bytes.NewReader(tables)

	ht, err := parseTables(r)
	assert.NoError(t, err)

	assert.Len(t, ht.soundOffsets.data, 200)

	for i := range ht.soundOffsets.data {
		assert.Equal(t, uint8(1), ht.soundOffsets.data[i])
	}

	assert.Len(t, ht.soundLengths.data, 100)

	for i := range ht.soundLengths.data {
		assert.Equal(t, uint8(2), ht.soundLengths.data[i])
	}

	assert.Len(t, ht.soundFrequencies.data, 100)

	for i := range ht.soundFrequencies.data {
		assert.Equal(t, uint8(3), ht.soundFrequencies.data[i])
	}
}

func sliceToBytes[T any](input []T) []byte {
	var gotBytes []byte

	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	for i := range input {
		_ = binary.Write(w, binary.LittleEndian, input[i])
		_ = w.Flush()
		buf := []byte{0, 0}
		b.Read(buf)
		b.Reset()
		gotBytes = append(gotBytes, buf...)
	}

	return gotBytes
}
