package combiner

import (
	"encoding/json"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

type Spritesheet struct {
	sourceImages []image.Image
}

type SheetDefinition struct {
	Prefix   string   `json:"prefix"`
	Suffixes []string `json:"suffixes"`
}

type SheetDefinitions []SheetDefinition

func (s *Spritesheet) AddImage(filename string) (err error) {
	if s.sourceImages == nil {
		s.sourceImages = make([]image.Image, 0)
	}

	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return err
	}

	img, err := png.Decode(f)
	if err != nil {
		return err
	}

	s.sourceImages = append(s.sourceImages, img)
	return nil
}

func (s *Spritesheet) GetOutput(filename string, margin int) {
	maxX, maxY := 0, 0
	indexed := false
	var pal color.Palette

	for _, img := range s.sourceImages {
		bounds := img.Bounds()
		if bounds.Max.X > maxX {
			maxX = bounds.Max.X
		}
		maxY += bounds.Max.Y + margin

		if p, ok := img.ColorModel().(color.Palette); ok {
			pal = p
			indexed = true
		}
	}

	var output image.Image
	if indexed {
		output = getIndexedOutput(maxX, maxY, pal, s.sourceImages, margin)
	} else {
		output = getRGBAOutput(maxX, maxY, s.sourceImages, margin)
	}

	outputFile, err := os.Create(filename)
	if err != nil {
		log.Fatalf("oops: %v", err)
	}

	png.Encode(outputFile, output)
}

func getIndexedOutput(maxX int, maxY int, pal color.Palette, images []image.Image, margin int) *image.Paletted {
	output := image.NewPaletted(image.Rectangle{Max: image.Point{X: maxX, Y: maxY}}, pal)
	maxIndex := uint8(len(pal) - 1)

	// clear to index 255
	for x := 0; x < maxX; x++ {
		for y := 0; y < maxY; y++ {
			output.SetColorIndex(x, y, maxIndex)
		}
	}

	y := 0
	for _, img := range images {
		pImg := img.(*image.Paletted)

		for x := 0; x < img.Bounds().Max.X; x++ {
			for j := 0; j < img.Bounds().Max.Y; j++ {
				output.SetColorIndex(x, y+j, pImg.ColorIndexAt(x, j))
			}
		}

		y += img.Bounds().Max.Y + margin
	}

	return output
}

func getRGBAOutput(maxX int, maxY int, images []image.Image, margin int) *image.RGBA {
	bounds := image.Rectangle{Max: image.Point{X: maxX, Y: maxY}}
	output := image.NewRGBA(bounds)
	// clear to white
	draw.Draw(output, bounds, &image.Uniform{C: color.White}, image.Point{}, draw.Src)

	y := 0
	for _, img := range images {
		drawBounds := image.Rectangle{
			Min: image.Point{X: 0, Y: y},
			Max: image.Point{X: img.Bounds().Max.X, Y: y + img.Bounds().Max.Y},
		}
		draw.Draw(output, drawBounds, img, image.Point{}, draw.Src)

		y += img.Bounds().Max.Y + margin
	}

	return output
}

func GetImageMap(definitions SheetDefinitions, files []string) (result map[string][]string) {
	sort.Strings(files)

	result = make(map[string][]string)
	for _, file := range files {
		for _, definition := range definitions {
			for _, suffix := range definition.Suffixes {
				if strings.HasPrefix(file, definition.Prefix) && strings.HasSuffix(file, suffix+".png") {
					name := definition.Prefix + "_" + suffix
					if r, ok := result[name]; !ok {
						result[name] = make([]string, 0)
						result[name] = append(result[name], file)
					} else {
						result[name] = append(r, file)
					}
				}
			}
		}
	}

	return result
}

func (s *SheetDefinitions) FromJSON(r io.Reader) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, s)
}
