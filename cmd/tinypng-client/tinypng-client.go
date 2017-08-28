package main

import (
	"flag"
	"fmt"
	"github.com/peterhellberg/tinypng"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path"
)

var (
	outputDir = flag.String("o", "out", "output diretory")
	inputFile = flag.String("f", "", "input file")
	inputDir  = flag.String("i", "", "input diretory")
	apiKey    = flag.String("key", "", "API Key")
	verbose   = flag.Bool("v", false, "Verbose")

	inputFiles []string
)

func main() {
	flag.Parse()

	switch {
	case *inputFile != "":
		inputFiles = []string{*inputFile}
	case *inputDir != "":
		readFiles, _ := ioutil.ReadDir(*inputDir)
		inputFiles = make([]string, len(readFiles))
		for no, f := range readFiles {
			inputFiles[no] = *inputDir + "/" + f.Name()
		}
	default:
		fmt.Println("No such input file or input diretory.")
		os.Exit(0)
	}

	for _, inputFilename := range inputFiles {

		// First check if the input file actually exist
		if !fileExists(inputFilename) {
			fmt.Println(inputFilename, ": Input file does not exist.")
			continue
		}

		// Verify that the input file is a PNG or JPEG file
		if !validFileType(inputFilename) {
			fmt.Println(inputFilename, ": Input file is not a valid PNG or JPEG file.")
			continue
		}

		// Then make sure that the output file doesn’t exist
		_, outputFileName := path.Split(inputFilename)
		if fileExists(*outputDir + "/" + outputFileName) {
			fmt.Println(outputFileName, ": Output file already exist.")
			continue
		}

		res, err := tinypng.ShrinkFn(*apiKey, inputFilename)

		if err != nil {
			fatal(res.Error, ":", res.Message)
		}

		// Print if *verbose is true
		if *verbose {
			res.Print()
		}

		// Download the compressed PNG
		os.MkdirAll(*outputDir, 0755)
		res.SaveAs(*outputDir + "/" + outputFileName)

		//Done!
		fmt.Println(outputFileName, ": Done!")

	}

}

// IO

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// Valid file type

func validFileType(fn string) bool {
	return validPNGFile(fn) || validJPEGFile(fn)
}

// PNG

func validPNGFile(fn string) bool {
	pngImage, err := os.Open(fn)

	check(err)

	defer pngImage.Close()

	_, err = png.DecodeConfig(pngImage)

	if err != nil {
		return false
	}

	return true
}

// JPEG

func validJPEGFile(fn string) bool {
	jpegImage, err := os.Open(fn)

	check(err)

	defer jpegImage.Close()

	_, err = jpeg.DecodeConfig(jpegImage)

	if err != nil {
		return false
	}

	return true
}

// Fatal

func check(err error) {
	if err != nil {
		fatal("Error:", err)
	}
}

func fatal(v ...interface{}) {
	fmt.Println(v...)

	os.Exit(1)
}
