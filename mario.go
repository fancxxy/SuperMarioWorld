package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fancxxy/SuperMarioWorld/generator"
)

const (
	output = "  "
)

var (
	flagGenerate   bool
	flagCharacter  string
	flagAction     string
	flagBackground string
	smw            *generator.SuperMario
	characters     []string
	actions        []string
)

func main() {
	generate()
	validate()
	draw()
}

func generate() {
	if !flagGenerate {
		return
	}

	var err error
	characters, actions, err = generator.Generate()
	if err != nil {
		fmt.Printf("Generate pixel data failed: %v\n", err)
		return
	}

	fmt.Printf("Generate pixel data succeed\n")
	fmt.Printf("Generate characters: %s\n", strings.Join(characters, ", "))
	fmt.Printf("Generate actions: %s\n", strings.Join(actions, ", "))
}

func draw() {
	print(smw.Characters[flagCharacter][flagAction], smw.Colors, output)
}

func print(data []byte, colors map[byte]string, output string) {
	buf := bufio.NewWriter(os.Stdout)
	defer buf.Flush()

	for _, char := range data {
		if char != '\n' {
			buf.WriteString(fmt.Sprintf("\033[48;2;%sm%s", colors[char], output))
		} else {
			buf.WriteString("\033[m\n")
		}
	}
}

func validate() {
	if !flagGenerate && ((flagCharacter == "") || (flagAction == "") || !find(characters, flagCharacter) || !find(actions, flagAction)) {
		flag.CommandLine.Usage()
		os.Exit(1)
	}

	rgb := strings.Split(flagBackground, ",")
	if len(rgb) != 3 || rgb[0] < "0" || rgb[0] > "255" || rgb[1] < "0" || rgb[1] > "255" || rgb[2] < "0" || rgb[2] > "255" {
		flag.CommandLine.Usage()
		os.Exit(1)
	}
	smw.Colors['A'] = strings.Replace(flagBackground, ",", ";", -1)

}

func init() {
	bytes, _ := ioutil.ReadFile("assets/SMW.json")
	smw = generator.NewSuperMario()
	json.Unmarshal(bytes, &smw)

	characters, actions = generator.Report(smw)
	flag.Usage = func() {
		fmt.Print("USAGE:\n\n")
		fmt.Printf("  Support characters: %s\n", strings.Join(characters, ", "))
		fmt.Printf("  Support actions: %s\n", strings.Join(actions, ", "))
		fmt.Printf("  Background color format: r,g,b, default value is 255,255,255\n\n")
		fmt.Print("OPTOINS:\n\n")
		flag.CommandLine.SetOutput(os.Stdout)
		flag.CommandLine.PrintDefaults()
		flag.CommandLine.SetOutput(ioutil.Discard)
		fmt.Println()
	}

	flag.CommandLine.SetOutput(ioutil.Discard)
	flag.CommandLine.Init(os.Args[0], flag.ExitOnError)
	flag.CommandLine.BoolVar(&flagGenerate, "g", false, "generate pixel json data")
	flag.CommandLine.StringVar(&flagCharacter, "c", "", "character name")
	flag.CommandLine.StringVar(&flagAction, "a", "", "action name")
	flag.CommandLine.StringVar(&flagBackground, "b", "255,255,255", "background color")

	flag.CommandLine.Parse(os.Args[1:])

	flagCharacter = strings.Title(flagCharacter)
	flagAction = strings.Title(flagAction)
}

func find(list []string, item string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}
