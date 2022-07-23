package vcl

import (
	"bytes"
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOpen(t *testing.T) {
	buf := newMockHeaders()
	r := bytes.NewReader(buf.Bytes())
	f, err := Open(r)

	assert.NoError(t, err)
	assert.NotNil(t, f)
}

func TestListSounds(t *testing.T) {
	buf := newMockHeaders()
	r := bytes.NewReader(buf.Bytes())
	f, _ := Open(r)

	sounds := ListSounds(f)
	assert.Len(t, sounds, 50)

	exp := 1
	for _, s := range sounds {
		assert.Equal(t, uint32(exp), s.Offset)
		assert.Equal(t, uint16(exp), s.Len)
		assert.Equal(t, uint16(exp), s.Freq)
		exp = exp + 1
	}
}

func TestListSoundsSkipsZeroLengthSounds(t *testing.T) {
	buf := newMockHeaders()
	r := bytes.NewReader(buf.Bytes())
	f, _ := Open(r)

	f.headers.Sound.Lengths[0] = 0
	f.headers.Sound.Lengths[3] = 0

	sounds := ListSounds(f)
	assert.Len(t, sounds, 48)
}

func TestReadPCM(t *testing.T) {
	var data = []byte{0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0}
	f := &File{
		headers: headers{},
		r:       bytes.NewReader(data),
	}

	sound := SoundInfo{4, 4, 0}
	pcmBuf, err := ReadPCM(sound, f)

	assert.NoError(t, err)
	assert.NotNil(t, pcmBuf)
	assert.Len(t, pcmBuf, 4)
	assert.Equal(t, []byte{1, 1, 1, 1}, pcmBuf)
}

func newMockHeaders() *bytes.Buffer {
	buf := new(bytes.Buffer)
	for i := 1; i <= 50; i++ {
		binary.Write(buf, binary.LittleEndian, uint32(i))
	}

	for i := 1; i <= 50; i++ {
		binary.Write(buf, binary.LittleEndian, uint16(i))
	}

	for i := 1; i <= 50; i++ {
		binary.Write(buf, binary.LittleEndian, uint16(i))
	}

	return buf
}
