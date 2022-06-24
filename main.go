package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/IljaN/unvcl/vcl"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
)

func main() {
	vclPath, extractPath := readArgsOrDie()

	file, err := os.Open(vclPath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	parsed, err := vcl.ParseFile(file)
	if err != nil {
		log.Fatal(err)
	}

	for i := range parsed.SoundTable {
		outPath := path.Join(extractPath, fmt.Sprintf("%s_%d.wav", fileNameWithoutExt(vclPath), i))
		if err = WriteWav(outPath, parsed.SoundTable[i].Samples); err != nil {
			log.Fatal(err)
		}
	}
}

func WriteWav(fileName string, samples []byte) error {
	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()

	e := wav.NewEncoder(out, 6000, 8, 1, 1)

	// Create new audio.IntBuffer.
	audioBuf, err := newAudioIntBuffer(bytes.NewReader(samples))
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

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func fileNameWithoutExt(fn string) string {
	fn = filepath.Base(fn)
	ext := filepath.Ext(fn)
	return fn[0 : len(fn)-len(ext)]
}

func readArgsOrDie() (vclPath, extractPath string) {
	if len(os.Args) < 3 {
		log.Fatal("unvcl VCL_FILE OUTPATH")
	}

	vclPath = os.Args[1]
	extractPath = os.Args[2]

	exists, _ := pathExists(extractPath)
	if !exists {
		fmt.Printf("Path %s does not exist", extractPath)
		os.Exit(1)
	}

	return vclPath, extractPath
}
