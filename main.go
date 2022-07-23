package main

import (
	"bytes"
	"fmt"
	"github.com/IljaN/unvcl/internal/encoder"
	"github.com/IljaN/unvcl/vcl"
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

	vclFile, err := vcl.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	for i, s := range vcl.ListSounds(vclFile) {
		outPath := path.Join(extractPath, fmt.Sprintf("%s_%d.wav", fileNameWithoutExt(vclPath), i))
		pcmSamples, _ := vcl.ReadPCM(s, vclFile)

		outFile, err := os.Create(outPath)
		if err != nil {
			log.Fatal(err)
		}

		if err = encoder.WriteWav(outFile, bytes.NewBuffer(pcmSamples), int(s.Freq)); err != nil {
			log.Fatal(err)
		}

		outFile.Close()
	}
}

func readArgsOrDie() (vclPath, extractPath string) {
	if len(os.Args) < 3 {
		fmt.Println("usage: unvcl vcl_file out_path")
		os.Exit(0)
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
