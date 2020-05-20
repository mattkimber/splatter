package main

import (
	"combiner"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Flags struct {
	InputDir, OutputDir string
	DefinitionFile      string
	Margin				int
}

var flags Flags

func init() {
	// Long format
	flag.StringVar(&flags.InputDir, "input_dir", "", "directory to scan for image files (default: current directory)")
	flag.StringVar(&flags.OutputDir, "output_dir", "", "directory to output files to (default: current directory)")
	flag.StringVar(&flags.DefinitionFile, "definition", "def.json", "JSON file containing the spritesheet definitions (default: def.json)")
	flag.IntVar(&flags.Margin, "margin", 0, "Vertical margin to leave between rows in the output")

	// Short format
	flag.StringVar(&flags.InputDir, "i", "", "shorthand for -input_dir")
	flag.StringVar(&flags.OutputDir, "o", "", "shorthand for -output_dir")
	flag.StringVar(&flags.DefinitionFile, "d", "def.json", "shorthand for -definition")
	flag.IntVar(&flags.Margin, "m", 0, "shorthand for -margin")

}

func main() {
	flag.Parse()

	ensureTrailingSlash(&flags.OutputDir)
	ensureTrailingSlash(&flags.InputDir)

	info, err := ioutil.ReadDir(flags.InputDir)
	if err != nil {
		log.Panicf("Couldn't stat directory: %v", err)
	}

	filenames := make([]string, 0)
	for _, fi := range info {
		filenames = append(filenames, fi.Name())
	}

	defFile, err := os.Open(flags.DefinitionFile)
	defer defFile.Close()
	if err != nil {
		log.Panicf("Could open definition file: %v", err)
	}

	definitions := combiner.SheetDefinitions{}
	if err := definitions.FromJSON(defFile); err != nil {
		log.Panicf("Couldn't read JSON definitions: %v", err)
	}

	m := combiner.GetImageMap(definitions, filenames)

	for k, v := range m {
		s := combiner.Spritesheet{}
		for _, img := range v {
			s.AddImage(flags.InputDir + img)
		}

		outputFilename := flags.OutputDir + k + ".png"
		log.Printf("Writing %s", outputFilename)
		s.GetOutput(outputFilename, flags.Margin)
	}

}

func ensureTrailingSlash(s *string) {
	if *s == "" {
		return
	}

	if !strings.HasSuffix(*s, "/") {
		*s = *s + "/"
	}
}
