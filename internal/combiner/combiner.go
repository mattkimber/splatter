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
	mask         image.Image
}

type SheetDefinition struct {
	Mask     string   `json:"mask"`
	Prefix   string   `json:"prefix"`
	Suffixes []string `json:"suffixes"`
	Prefixes []string `json:"prefixes"`
}

type SheetDefinitions []SheetDefinition

func (s *Spritesheet) AddImage(filename string) (err error) {
	if s.sourceImages == nil {
		s.sourceImages = make([]image.Image, 0)
	}

	img, err := getImage(filename)
	if err != nil {
		return err
	}

	s.sourceImages = append(s.sourceImages, img)
	return nil
}

func (s *Spritesheet) AddMask(filename string) (err error) {
	img, err := getImage(filename)
	if err != nil {
		return err
	}

	s.mask = img
	return nil
}

func getImage(filename string) (img image.Image, err error) {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return
	}

	img, err = png.Decode(f)
	if err != nil {
		return
	}

	return
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
		output = getIndexedOutput(maxX, maxY, pal, s.sourceImages, s.mask, margin)
	} else {
		output = getRGBAOutput(maxX, maxY, s.sourceImages, s.mask, margin)
	}

	outputFile, err := os.Create(filename)
	if err != nil {
		log.Fatalf("oops: %v", err)
	}

	png.Encode(outputFile, output)
}

func getIndexedOutput(maxX int, maxY int, pal color.Palette, images []image.Image, mask image.Image, margin int) *image.Paletted {
	output := image.NewPaletted(image.Rectangle{Max: image.Point{X: maxX, Y: maxY}}, pal)
	maxIndex := uint8(len(pal) - 1)

	hasMask := false
	if mask != nil {
		hasMask = true
	}

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
				if !hasMask {
					output.SetColorIndex(x, y+j, pImg.ColorIndexAt(x, j))
				} else {
					_, _, _, a := mask.At(x, y+j).RGBA()
					if a == 0 {
						output.SetColorIndex(x, y+j, pImg.ColorIndexAt(x, j))
					} else {
						output.SetColorIndex(x, y+j, 0)
					}
				}

			}
		}

		y += img.Bounds().Max.Y + margin
	}

	return output
}

func getRGBAOutput(maxX int, maxY int, images []image.Image, mask image.Image, margin int) *image.RGBA {
	bounds := image.Rectangle{Max: image.Point{X: maxX, Y: maxY}}
	output := image.NewRGBA(bounds)

	hasMask := false
	if mask != nil {
		hasMask = true
	}

	// clear to white
	draw.Draw(output, bounds, &image.Uniform{C: color.White}, image.Point{}, draw.Src)

	y := 0
	for _, img := range images {
		for x := 0; x < img.Bounds().Max.X; x++ {
			for j := 0; j < img.Bounds().Max.Y; j++ {
				if !hasMask {
					output.Set(x, y+j, img.At(x, j))
				} else {
					_, _, _, a := mask.At(x, y+j).RGBA()
					if a == 0 {
						output.Set(x, y+j, img.At(x, j))
					} else {
						output.Set(x, y+j, color.Transparent)
					}
				}
			}
		}

		y += img.Bounds().Max.Y + margin
	}

	return output
}

type ImageSpec struct {
	Files []string
	Mask  string
}

type ImageSpecMap map[string]ImageSpec

func GetImageMap(definitions SheetDefinitions, files []string) (result ImageSpecMap) {
	sort.Strings(files)

	result = make(map[string]ImageSpec)

	for _, definition := range definitions {

		if len(definition.Prefixes) == 0 {
			definition.Prefixes = []string{definition.Prefix}
		}

		for _, prefix := range definition.Prefixes {
			for _, suffix := range definition.Suffixes {
				name := prefix + "_" + suffix

				is := ImageSpec{}
				is.Mask = definition.Mask
				is.Files = make([]string, 0)

				for _, file := range files {
					if strings.HasPrefix(file, prefix) && strings.HasSuffix(file, suffix+".png") {
						is.Files = append(is.Files, file)
					}
				}

				result[name] = is
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
