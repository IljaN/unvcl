package vcl

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestReadOffset(t *testing.T) {

	testBytes := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	testReader := bytes.NewReader(testBytes)

	rd, err := readOffset(testReader, 5, 5)
	assert.NoError(t, err)
	got, _ := ioutil.ReadAll(rd)

	assert.Len(t, got, 5)
	assert.Equal(t, testBytes[5:10], got)

	rd, err = readOffset(testReader, 0, 10)
	assert.NoError(t, err)
	got, _ = ioutil.ReadAll(rd)

	assert.Len(t, got, 10)
	assert.Equal(t, testBytes[0:10], got)
}
