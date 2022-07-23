package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileNameWithoutExt(t *testing.T) {
	res := fileNameWithoutExt("/foo/bar/baz/file.ext")
	assert.Equal(t, "file", res)

	res = fileNameWithoutExt("/foo/bar/baz/file")
	assert.Equal(t, "file", res)

	res = fileNameWithoutExt("file")
	assert.Equal(t, "file", res)
}

func TestPathExists(t *testing.T) {
	d := t.TempDir()
	exists, err := pathExists(d)
	assert.NoError(t, err)
	assert.Truef(t, exists, "Path must exist")
}
