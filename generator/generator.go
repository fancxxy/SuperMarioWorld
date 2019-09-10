package generator

import (
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

type SuperMario struct {
	Characters map[string]map[string][]byte `json:"characters"`
	Colors     map[byte]string              `json:"colors"`
}

const (
	size      = 32
	ending    = '\n'
	imageFile = "assets/SMW.png"
	pointFile = "assets/SMW.point"
	jsonFile  = "assets/SMW.json"
)

var (
	symbols = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	index   = 0 // 最后一个使用的字符下标
	colors  = make(map[string]byte)
)

func Generate() ([]string, []string, error) {
	image, err := InitImage()
	if err != nil {
		return nil, nil, err
	}

	points, err := InitPoints()
	if err != nil {
		return nil, nil, err
	}

	return generate(image, points)
}

func InitImage() (image.Image, error) {
	file, err := os.Open(imageFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	image, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func InitPoints() (map[string]map[string]image.Point, error) {
	points := make(map[string]map[string]image.Point)
	bytes, err := ioutil.ReadFile(pointFile)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &points)
	if err != nil {
		return nil, err
	}

	return points, nil
}

func generate(img image.Image, characterActionPoints map[string]map[string]image.Point) ([]string, []string, error) {
	smw := NewSuperMario()
	for character, actionPoints := range characterActionPoints {
		actionBytes := make(map[string][]byte)
		for action, point := range actionPoints {
			bytes := ansi(img, point)
			actionBytes[action] = bytes
		}
		smw.Characters[character] = actionBytes
	}

	for color, symbol := range colors {
		smw.Colors[symbol] = color
	}

	bytes, err := json.Marshal(smw)
	if err != nil {
		return nil, nil, err
	}

	err = ioutil.WriteFile(jsonFile, bytes, 0644)
	if err != nil {
		return nil, nil, err
	}

	cs, as := Report(smw)

	return cs, as, nil
}

func Report(smw *SuperMario) ([]string, []string) {
	var characters, actions []string
	for character := range smw.Characters {
		characters = append(characters, character)
	}

	if len(characters) == 0 {
		return nil, nil
	}

	for action := range smw.Characters[characters[0]] {
		actions = append(actions, action)
	}
	sort.Strings(characters)
	sort.Strings(actions)

	return characters, actions
}

func ansi(img image.Image, point image.Point) []byte {
	var chars []byte
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			r, g, b, _ := img.At(x+point.X, y+point.Y).RGBA()
			color := fmt.Sprintf("%d;%d;%d", r>>8, g>>8, b>>8)
			char, ok := colors[color]
			if !ok {
				if index >= len(symbols) {
					log.Fatal("symbols are not enough")
				}
				char = symbols[index]
				index++
				colors[color] = char
			}
			chars = append(chars, char)
		}
		chars = append(chars, ending)
	}

	return chars
}

func NewSuperMario() *SuperMario {
	smw := new(SuperMario)
	smw.Characters = make(map[string]map[string][]byte)
	smw.Colors = make(map[byte]string)

	return smw
}
