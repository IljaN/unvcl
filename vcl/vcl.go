// Package vcl provides apis to extract audio-samples from VCL-Files. VCL is a file-format which is mainly used by some early DOS
// games developed by Epic MegaGames.
package vcl

import (
	"encoding/binary"
	"io"
)

// File represents a VCL File
//
// The text portion of the format is currently not implemented.
type File struct {
	headers
	r io.ReadSeeker
}

// The format consists of five tables: Sound Offsets table, Sound Lengths table, Sound Frequencies table, Text Offset table, and Text Length table.
// The order of occurrence is important; the game uses indexes to this table to determine which sound to play or what text to show. The sounds
// are stored in raw format without headers.
type headers struct {
	Sound struct {
		// OffsetTable: 50x uint32, marks the offset of every sound-file. Offset values are absolute to the file.
		Offsets [50]uint32
		// LengthTable: 50x uint16, denotes the length of the sound in uint8 relative to offset
		Lengths [50]uint16
		// FreqTable: 50x uint16, the sampling frequency at which the sample should be played  (not implemented)
		Frequencies [50]uint16
	}
}

// Open a VCL file for data extraction
func Open(vcl io.ReadSeeker) (*File, error) {
	f := &File{
		headers: headers{},
		r:       vcl,
	}

	_, err := f.r.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	if err = binary.Read(f.r, binary.LittleEndian, &f.headers); err != nil {
		return nil, err
	}
	return f, nil
}

// SoundInfo holds the parameters of a single sound-file inside the vcl file.
type SoundInfo struct {
	Offset uint32
	Len    uint16
	Freq   uint16
}

// ListSounds returns a list of sounds stored in File
func ListSounds(f *File) []SoundInfo {
	list := make([]SoundInfo, 0)
	sh := f.headers.Sound
	for i := range sh.Offsets {
		offs, length, freq := sh.Offsets[i], sh.Lengths[i], sh.Frequencies[i]
		if length == 0 {
			continue
		}

		list = append(list, SoundInfo{
			Offset: offs,
			Len:    length,
			Freq:   freq,
		})
	}

	return list
}

// ReadPCM reads the raw 8-bit unigned PCM-Data associated with a given sound
func ReadPCM(s SoundInfo, vcl *File) ([]byte, error) {
	if _, err := vcl.r.Seek(int64(s.Offset), io.SeekStart); err != nil {
		return nil, err
	}

	buf := make([]byte, s.Len)

	if err := binary.Read(vcl.r, binary.LittleEndian, &buf); err != nil {
		return nil, err
	}

	return buf, nil
}
