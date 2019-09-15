package generator

import (
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
)

// Generate function generate ascii data
func Generate() (*SMW, error) {
	generated := new(SMW)
	generated.Arts = make(map[string]*Art)
	generated.Colors = make(map[byte]string)

	color, err := initColor()
	if err != nil {
		return generated, err
	}

	image, err := initImage()
	if err != nil {
		return generated, err
	}

	points, err := initPoints()
	if err != nil {
		return generated, err
	}

	for character, actions := range points {
		for action, point := range actions {
			point.X += (size - width) / 2
			point.Y += (size - height)
			bytes := pixel2Char(image, point, color)
			generated.Arts[character+"."+action] = &Art{
				Character: character,
				Action:    action,
				Width:     width,
				Height:    height,
				Data:      bytes,
			}
		}
	}

	for rgb, symbol := range color.mapping {
		if symbol == 'A' {
			rgb = ";;"
		}
		generated.Colors[symbol] = rgb
	}

	generated.statistics()

	bytes, err := json.Marshal(generated)
	if err != nil {
		return generated, err
	}

	err = ioutil.WriteFile(asciifile, bytes, 0644)
	if err != nil {
		return generated, err
	}

	return generated, nil
}

func pixel2Char(img image.Image, point image.Point, color *color) []byte {
	var chars []byte
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x+point.X, y+point.Y).RGBA()
			rgb := fmt.Sprintf("%d;%d;%d", r>>8, g>>8, b>>8)
			chars = append(chars, color.symbol(rgb))
		}
	}

	return chars
}
