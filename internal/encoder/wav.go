package encoder

import (
	"encoding/binary"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"io"
)

func WriteWav(out io.WriteSeeker, samples io.Reader, sampleRate int) error {
	e := wav.NewEncoder(out, sampleRate, 8, 1, 1)

	// Create new audio.IntBuffer.
	audioBuf, err := newAudioIntBuffer(samples, sampleRate)
	if err != nil {
		return err
	}
	// Write buffer to output file. This writes a RIFF header and the PCM chunks from the audio.IntBuffer.
	if err = e.Write(audioBuf); err != nil {
		return err
	}
	if err = e.Close(); err != nil {
		return err
	}

	return nil

}

func newAudioIntBuffer(r io.Reader, sampleRate int) (*audio.IntBuffer, error) {
	buf := audio.IntBuffer{
		Format: &audio.Format{
			NumChannels: 1,
			SampleRate:  sampleRate,
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
