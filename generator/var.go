package main

import (
	"image"
	"image/png"
	"os"
)

var (
	width    = 28
	height   = 24
	chars    = ".abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	output   = "../smw.go"
	resource = "assets/smw.png"

	characters = map[string]int{
		"mario":    30,
		"luigi":    661,
		"toad":     1291,
		"toadette": 1921,
	}
	actions = map[string]int{
		"idle":       3,
		"up":         40,
		"crouch":     77,
		"walk":       149,
		"run":        217,
		"accelerate": 250,
		"skid":       287,
		"jump":       324,
		"fall":       357,
		"fly":        394,
		"front":      497,
		"left":       530,
		"back":       563,
	}
)

func initSMW() *smw {
	return &smw{
		Palette: make(map[byte]string),
		ASCII:   make(map[string][]string),
	}
}

func initColor() (*color, error) {
	c := new(color)
	c.index = 0
	c.mapping = make(map[string]byte)
	c.chars = []byte(chars)

	return c, nil
}

func initImage() (image.Image, error) {
	file, err := os.Open(resource)
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

func initPoints() (map[string]image.Point, error) {
	points := make(map[string]image.Point)

	for character, y := range characters {
		for action, x := range actions {
			points[character+"."+action] = image.Point{X: x, Y: y}
		}
	}

	return points, nil
}
