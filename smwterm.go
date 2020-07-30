package main

import (
	"flag"
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
)

type arguments struct {
	Frames [][]charact
	Itv    int
}

type charact struct {
	Point image.Point
	Name  string
}

func main() {
	draw(parse())
}

func parse() arguments {
	flag.Usage = func() {
		_, file := filepath.Split(os.Args[0])
		fmt.Print("USAGE:\n\n")
		fmt.Printf("  %s [options] frame...\n\n", file)
		fmt.Printf("  Frame is x0,y0,character0,action0/x1,y1,character0,action1/...\n")
		fmt.Printf("  Support characters: %s\n", characters)
		fmt.Printf("  Support actions: %s\n", actions)
		fmt.Printf("  Format of background color is r,g,b\n")
		fmt.Printf("  characters' ascii code is generated by generator\n\n")
		fmt.Print("OPTOINS:\n\n")
		flag.CommandLine.SetOutput(os.Stdout)
		flag.CommandLine.PrintDefaults()
		flag.CommandLine.SetOutput(ioutil.Discard)
		fmt.Println()
	}
	flag.CommandLine.SetOutput(ioutil.Discard)
	flag.CommandLine.Init(os.Args[0], flag.ExitOnError)
	color := flag.CommandLine.String("b", "248,206,1", "background color")
	fps := flag.CommandLine.Int("f", 5, "frames per second")
	flag.CommandLine.Parse(os.Args[1:])
	frames := flag.Args()

	fs := make([][]charact, len(frames))
	for i, frame := range frames {
		frame = strings.Trim(frame, "\"")
		characts := strings.Split(frame, "/")
		fs[i] = make([]charact, len(characts))
		for j, charact := range characts {
			params := strings.Split(charact, ",")
			if len(params) != 4 {
				flag.CommandLine.Usage()
				os.Exit(1)
			}
			x, err := strconv.Atoi(params[0])
			if err != nil {
				flag.CommandLine.Usage()
				os.Exit(1)
			}
			y, err := strconv.Atoi(params[1])
			if err != nil {
				flag.CommandLine.Usage()
				os.Exit(1)
			}

			if _, ok := ASCII[params[2]+"."+params[3]]; !ok {
				flag.CommandLine.Usage()
				os.Exit(1)
			}

			fs[i][j].Name = params[2] + "." + params[3]
			fs[i][j].Point = image.Point{X: x, Y: y}
		}
	}

	rgb := strings.Split(*color, ",")
	if len(rgb) != 3 {
		flag.CommandLine.Usage()
		os.Exit(1)
	}
	if *color != ",," {
		for i := range rgb {
			if _, err := strconv.Atoi(rgb[i]); err != nil {
				flag.CommandLine.Usage()
				os.Exit(1)
			}
		}
	}

	palettes[char] = strings.Replace(*color, ",", ";", -1)

	if *fps <= 0 {
		*fps = 5
	}

	return arguments{Frames: fs, Itv: 1000 / (*fps)}
}

func draw(args arguments) {
	if len(args.Frames) == 0 {
		return
	}

	register()
	re := newRenderer()
	re.render(args)
}

func register() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func(c chan os.Signal) {
		<-c
		fmt.Print(cleanup())
		os.Exit(0)
	}(c)
}
