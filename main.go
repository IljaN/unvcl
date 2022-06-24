package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"io"
	"io/ioutil"
	"log"
	"os"
)

//  ffplay -f s8 -ar 6400 28.wav
func main() {

	file, err := ioutil.ReadFile("jill1.vcl")
	if err != nil {
		log.Fatal(err)
	}

	offsetTable := file[0:200]
	lengthTable := file[200:300]
	freqTable := file[300:400]

	var (
		u32min = 0
		u32max = 4
		u16min = 0
		u16max = 2
	)

	var offsets = make([]uint32, 50)
	var lengths = make([]uint16, 50)
	var freqs = make([]uint16, 50)

	for i := range offsets {
		offsets[i] = binary.LittleEndian.Uint32(offsetTable[u32min:u32max])
		lengths[i] = binary.LittleEndian.Uint16(lengthTable[u16min:u16max])
		freqs[i] = binary.LittleEndian.Uint16(freqTable[u16min:u16max])

		if lengths[i] != 0 {

			snd := file[offsets[i]:(offsets[i] + uint32(lengths[i]))]

			// Output file.
			out, err := os.Create(fmt.Sprintf("%d.wav", i))
			if err != nil {
				log.Fatal(err)
			}
			defer out.Close()

			e := wav.NewEncoder(out, 6000, 8, 1, 1)

			// Create new audio.IntBuffer.
			audioBuf, err := newAudioIntBuffer(bytes.NewReader(snd))
			if err != nil {
				log.Fatal(err)
			}
			// Write buffer to output file. This writes a RIFF header and the PCM chunks from the audio.IntBuffer.
			if err := e.Write(audioBuf); err != nil {
				log.Fatal(err)
			}
			if err := e.Close(); err != nil {
				log.Fatal(err)
			}

			out.Close()
		}

		if u32max == 200 {
			break
		}
		u32min = u32min + 4
		u32max = u32max + 4
		u16min = u16min + 2
		u16max = u16max + 2

	}

}

func newAudioIntBuffer(r io.Reader) (*audio.IntBuffer, error) {
	buf := audio.IntBuffer{
		Format: &audio.Format{
			NumChannels: 1,
			SampleRate:  6000,
		},
	}
	for {
		var sample uint8
		err := binary.Read(r, binary.LittleEndian, &sample)
		switch {
		case err == io.EOF:
			return &buf, nil
		case err != nil:
			return nil, err
		}
		buf.Data = append(buf.Data, int(sample))
	}
}
