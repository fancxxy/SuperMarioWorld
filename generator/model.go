package generator

import (
	"encoding/json"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"sort"
)

// SMW reprenset generated ascii data
type SMW struct {
	Arts       map[string]*Art `json:"arts"`
	Colors     map[byte]string `json:"colors"`
	Characters []string        `json:"characters"`
	Actions    []string        `json:"actions"`
}

// Art represent an action of a character
type Art struct {
	Character string `json:"character"`
	Action    string `json:"action"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Data      []byte `json:"data"`
}

type position struct {
	Characters map[string]int `json:"characters"`
	Actions    map[string]int `json:"actions"`
}

type color struct {
	symbols  []byte
	indexing int
	mapping  map[string]byte
}

func (smw *SMW) statistics() {
	cset := make(map[string]bool)
	aset := make(map[string]bool)

	for _, art := range smw.Arts {
		cset[art.Character] = true
		aset[art.Action] = true
	}

	var characters, actions []string
	for character := range cset {
		characters = append(characters, character)
	}
	for action := range aset {
		actions = append(actions, action)
	}

	sort.Strings(characters)
	sort.Strings(actions)

	smw.Characters = characters
	smw.Actions = actions
}

func (c *color) symbol(rgb string) byte {
	if b, ok := c.mapping[rgb]; ok {
		return b
	}

	if c.indexing >= len(c.symbols) {
		return '.'
	}

	c.mapping[rgb] = c.symbols[c.indexing]
	c.indexing++
	return c.mapping[rgb]
}

func initColor() (*color, error) {
	c := new(color)
	c.indexing = 0
	c.mapping = make(map[string]byte)
	c.symbols = []byte(symbols)

	return c, nil
}

func initImage() (image.Image, error) {
	file, err := os.Open(imagefile)
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

func initPoints() (map[string]map[string]image.Point, error) {
	points := make(map[string]map[string]image.Point)

	bytes, err := ioutil.ReadFile(positionfile)
	if err != nil {
		return points, err
	}

	var pn position
	err = json.Unmarshal(bytes, &pn)
	if err != nil {
		return points, err
	}

	for character, Y := range pn.Characters {
		ps := make(map[string]image.Point)
		for action, X := range pn.Actions {
			ps[action] = image.Point{X: X, Y: Y}
		}
		points[character] = ps
	}

	return points, nil
}
