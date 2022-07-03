package vcl

import (
	"bytes"
	"encoding/binary"
	"io"
)

const (
	OffsetTableStart = 0
	OffsetTableEnd   = 200
	LengthTableStart = 200
	LengthTableEnd   = 300
	FreqTableStart   = 300
	FreqTableEnd     = 400
)

// NumSounds is the max. number of sounds supported by the file-format for a single file
const NumSounds = 50

// File is a parsed VCL file. It holds Sounds when parsed.
type File struct {
	Sounds []Sound
}

// Sound is a parsed sound block in the vcl.File. Samples is PCM 8-Bit unsigned.
type Sound struct {
	Offset  uint32
	Len     uint16
	Freq    uint16
	Samples []byte
}

// ParseFile parses a VCL file and loads all sounds in to memory. See File
func ParseFile(r io.ReadSeeker) (*File, error) {
	vcl := &File{}

	tables, err := parseTables(r)
	if err != nil {
		return nil, err
	}

	for i := 0; i < NumSounds; i++ {
		offs, err := tables.soundOffsets.Next()
		if err != nil {
			return nil, err
		}
		length, err := tables.soundLengths.Next()
		if err != nil {
			return nil, err
		}
		freq, err := tables.soundFrequencies.Next()
		if err != nil {
			return nil, err
		}

		s := Sound{
			Offset: offs,
			Len:    length,
			Freq:   freq,
		}

		if s.Len != 0 {
			// read audio samples
			pcm, err := read(r, s.Offset, s.Len)
			if err != nil {
				return nil, err
			}

			s.Samples = pcm
			vcl.Sounds = append(vcl.Sounds, s)
		}
	}

	return vcl, nil
}

// headerTables bundles the sound-specific lookup-tables of the format
type headerTables struct {
	// 000-200 OffsetTable: 50x uint32, marks the offset of every sound-file. Offset values are absolute to the file.
	soundOffsets *table[uint32]
	// 200-300 LengthTable: 50x uint16, denotes the length of the sound in uint8 relative to offset
	soundLengths *table[uint16]
	// 300-400 FreqTable: 50x uint16, the sampling frequency at which the sample should be played  (not implemented)
	soundFrequencies *table[uint16]
}

// parseTables parses the headerTables which are required to lookup the sound files
func parseTables(r io.ReadSeeker) (*headerTables, error) {
	tables := &headerTables{}

	offs, err := parseTable[uint32](r, OffsetTableStart, OffsetTableEnd-OffsetTableStart)
	if err != nil {
		return nil, err
	}
	tables.soundOffsets = offs

	lengths, err := parseTable[uint16](r, LengthTableStart, LengthTableEnd-LengthTableStart)
	if err != nil {
		return nil, err
	}
	tables.soundLengths = lengths

	freqs, err := parseTable[uint16](r, FreqTableStart, FreqTableEnd-FreqTableStart)
	if err != nil {
		return nil, err
	}
	tables.soundFrequencies = freqs

	return tables, nil
}

// table is a single data table in the file-format
type table[T uint32 | uint16] struct {
	data []byte
	r    io.ReadSeeker
}

// parseTable reads len bytes at offset and parses them to a table[T]
func parseTable[T uint32 | uint16](r io.ReadSeeker, offset uint32, len uint16) (*table[T], error) {
	ofr, err := read(r, offset, len)
	if err != nil {
		return nil, err
	}

	tbl := &table[T]{
		data: ofr,
	}
	tbl.r = bytes.NewReader(tbl.data)

	return tbl, nil
}

func (tb *table[T]) Next() (T, error) {
	var x T = 0
	err := binary.Read(tb.r, binary.LittleEndian, &x)
	return x, err
}

// All seeks to the beginning and parses all values from the table
func (tb *table[T]) All() ([]T, error) {
	_, err := tb.r.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	res := make([]T, 0)

	for {
		var v T = 0
		if err = binary.Read(tb.r, binary.LittleEndian, &v); err != nil {
			if err == io.EOF {
				return res, nil
			}
		}

		res = append(res, v)
	}
}

// read seeks to offset and reads len
func read(r io.ReadSeeker, offset uint32, len uint16) ([]byte, error) {
	var data = make([]byte, len)
	var err error

	if _, err = r.Seek(int64(offset), 0); err != nil {
		return nil, err
	}
	if _, err = r.Read(data); err != nil {
		return nil, err
	}

	return data, nil
}
