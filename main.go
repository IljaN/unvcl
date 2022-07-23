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

var Version string

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
		pcmSamples, err := vcl.ReadPCM(s, vclFile)
		if err != nil {
			log.Fatal(err)
		}

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
		fmt.Println("Usage: unvcl vcl_file out_path")
		fmt.Println("Extracts audio-files from a vcl-file")
		fmt.Println("\t vcl_file\tSource file from which to extract")
		fmt.Println("\t out_path\tTarget directory where audio-files should saved")
		fmt.Println("")
		fmt.Println("Version: " + Version)

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
