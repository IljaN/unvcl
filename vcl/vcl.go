package vcl

import (
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
)

const (
	OffsetTableStart = 0
	OffsetTableEnd   = 200
	LengthTableStart = 200
	LengthTableEnd   = 300
	FreqTableStart   = 300
	FreqTableEnd     = 400
)

const NumSound = 50

// File is a parsed VCL file. It contains a SoundTable with all sounds.
type File struct {
	SoundTable []Sound
}

// Sound is a single sound block in the vcl.File. Samples is PCM 8-Bit unsigned.
type Sound struct {
	Offset  uint32
	Len     uint16
	Freq    uint16
	Samples []byte
}

// ParseFile parses a VCL file and loads all sounds in to memory. See File
func ParseFile(r io.ReadSeeker) (*File, error) {
	vcl := &File{}

	offsetTable, err := readOffset(r, OffsetTableStart, OffsetTableEnd-OffsetTableStart)
	if err != nil {
		return nil, err
	}
	lengthTable, err := readOffset(r, LengthTableStart, LengthTableEnd-LengthTableStart)
	if err != nil {
		return nil, err
	}
	freqTable, err := readOffset(r, FreqTableStart, FreqTableEnd-FreqTableStart)
	if err != nil {
		return nil, err
	}

	for i := 0; i < NumSound; i++ {
		s := Sound{}
		if err = binary.Read(offsetTable, binary.LittleEndian, &s.Offset); err != nil {
			return nil, err
		}

		if err = binary.Read(lengthTable, binary.LittleEndian, &s.Len); err != nil {
			return nil, err
		}

		if err = binary.Read(freqTable, binary.LittleEndian, &s.Freq); err != nil {
			return nil, err
		}

		if s.Len != 0 {
			sr, _ := readOffset(r, s.Offset, s.Len)
			smp, _ := ioutil.ReadAll(sr)
			s.Samples = smp
			vcl.SoundTable = append(vcl.SoundTable, s)
		}
	}

	return vcl, nil
}

func readOffset(r io.ReadSeeker, offset uint32, len uint16) (io.Reader, error) {
	var err error
	if _, err = r.Seek(int64(offset), 0); err != nil {
		return nil, err
	}
	data := make([]byte, len)
	if _, err = r.Read(data); err != nil {
		return nil, err
	}

	return bytes.NewReader(data), nil
}
