package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fancxxy/smwterm/generator"
)

const (
	bg = 'A' // background symbol
)

type arguments struct {
	Arts  []string
	Itv   int
	Point image.Point
}

func main() {
	smw := load("assets/SMW.json")
	arguments := parse(smw)
	draw(smw, arguments)
}

func load(file string) *generator.SMW {
	bytes, _ := ioutil.ReadFile(file)
	var smw *generator.SMW
	json.Unmarshal(bytes, &smw)
	return smw
}

func parse(smw *generator.SMW) arguments {
	flag.Usage = func() {
		_, file := filepath.Split(os.Args[0])
		fmt.Print("USAGE:\n\n")
		fmt.Printf("  %s [options] art...\n\n", file)
		fmt.Printf("  Art is character.action\n")
		fmt.Printf("  Support characters: %s\n", strings.Join(smw.Characters, ", "))
		fmt.Printf("  Support actions: %s\n", strings.Join(smw.Actions, ", "))
		fmt.Printf("  Format of background color is r,g,b\n")
		fmt.Printf("  It's suggested to run -g to generate newest ascii data at first time\n\n")
		fmt.Print("OPTOINS:\n\n")
		flag.CommandLine.SetOutput(os.Stdout)
		flag.CommandLine.PrintDefaults()
		flag.CommandLine.SetOutput(ioutil.Discard)
		fmt.Println()
	}
	flag.CommandLine.SetOutput(ioutil.Discard)
	flag.CommandLine.Init(os.Args[0], flag.ExitOnError)
	generated := flag.CommandLine.Bool("g", false, "generate ascii data")
	background := flag.CommandLine.String("b", "248,206,1", "background color")
	fps := flag.CommandLine.Int("f", 30, "frames per second")
	point := flag.CommandLine.String("p", "0,0", "top left point in terminal")
	flag.CommandLine.Parse(os.Args[1:])
	arts := flag.Args()

	if *generated {
		smw = generate()
	}

	if !check(smw, arts) {
		flag.CommandLine.Usage()
		os.Exit(1)
	}

	rgb := strings.Split(*background, ",")
	if len(rgb) != 3 {
		flag.CommandLine.Usage()
		os.Exit(1)
	}
	if *background != ",," {
		for _, s := range rgb {
			if _, err := strconv.Atoi(s); err != nil {
				flag.CommandLine.Usage()
				os.Exit(1)
			}
		}
	}

	smw.Colors[bg] = strings.Replace(*background, ",", ";", -1)

	xy := strings.Split(*point, ",")
	if len(xy) != 2 {
		flag.CommandLine.Usage()
		os.Exit(1)
	}
	x, err := strconv.Atoi(xy[0])
	if err != nil {
		flag.CommandLine.Usage()
		os.Exit(1)
	}
	y, err := strconv.Atoi(xy[1])
	if err != nil {
		flag.CommandLine.Usage()
		os.Exit(1)
	}

	return arguments{Arts: arts, Itv: 1000 / (*fps), Point: image.Point{X: x, Y: y}}
}

func generate() *generator.SMW {
	smw, err := generator.Generate()
	if err != nil {
		fmt.Printf("Generate pixel data failed: %v\n", err)
		return smw
	}

	var colors string
	for _, c := range smw.Colors {
		colors += color(c)
	}

	fmt.Printf("Generate pixel data succeed\n")
	fmt.Printf("Generate characters: %s\n", strings.Join(smw.Characters, ", "))
	fmt.Printf("Generate actions: %s\n", strings.Join(smw.Actions, ", "))
	fmt.Printf("Generate colors: %s\n", colors)
	return smw
}

func draw(smw *generator.SMW, args arguments) {
	if len(args.Arts) == 0 {
		return
	}

	callback()
	buf := bufio.NewWriter(os.Stdout)
	buf.WriteString(initialize())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	renderer := render(ctx, smw, args)

	for {
		print(buf, <-renderer, smw.Colors)
		time.Sleep(time.Duration(args.Itv) * time.Millisecond)
	}
}

func render(ctx context.Context, smw *generator.SMW, args arguments) <-chan []byte {
	width, height := term()
	width /= 2
	raw := make([]byte, (width+1)*height)

	for j := 0; j < height; j++ {
		for i := 0; i < width+1; i++ {
			raw[j*(width+1)+i] = bg
		}
		raw[j*(width+1)+width] = '\n'
	}
	raw = raw[:(width+1)*height-1]

	channel := make(chan []byte)

	go func() {
		i := 0
		for {
			data := make([]byte, len(raw))
			copy(data, raw)

			assign(data, width, height, smw.Arts[args.Arts[i]], args.Point)
			select {
			case <-ctx.Done():
				return
			case channel <- data:
			}

			i++
			if i == len(args.Arts) {
				i = 0
			}
		}
	}()
	return channel
}

func assign(data []byte, width, height int, art *generator.Art, point image.Point) {
	for j := 0; j < art.Height; j++ {
		for i := 0; i < art.Width; i++ {
			if j+point.Y > 0 && i+point.X > 0 && j+point.Y < height &&
				i+point.X < width && art.Data[j*art.Width+i] != bg {
				data[(j+point.Y)*(width+1)+(i+point.X)] = art.Data[j*art.Width+i]
			}
		}
	}

}

func print(buf *bufio.Writer, data []byte, colors map[byte]string) {
	defer buf.Flush()
	buf.WriteString(reset())
	for _, char := range data {
		switch char {
		case '\n':
			buf.WriteString(linefeed())
		default:
			buf.WriteString(color(colors[char]))
		}
	}
}

func check(smw *generator.SMW, arts []string) bool {
	for _, art := range arts {
		splits := strings.Split(art, ".")
		if len(splits) != 2 {
			return false
		}

		character, action := splits[0], splits[1]
		if character == "" || !find(smw.Characters, character) {
			return false
		}

		if action == "" || !find(smw.Actions, action) {
			return false
		}
	}

	return true
}

func find(list []string, item string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}

func callback() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func(c chan os.Signal) {
		<-c
		fmt.Print(cleanup())
		os.Exit(0)
	}(c)
}
