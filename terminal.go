package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func term() (int, int) {
	width, height, err := terminal.GetSize(int(os.Stdout.Fd()))
	if width == 0 || height == 0 || err != nil {
		width, height = 80, 24
	}

	return width, height
}

func reset() string {
	return "\033[H"
}

func color(rgb string) string {
	return fmt.Sprintf("\033[48;2;%sm%s", rgb, "  ")
}

/*
	\033[?25l 隐藏光标
	\033[?25h 显示光标
	\033]50;SetProfile=smw\a 加载iterm2的配置
*/
func initialize() string {
	return "\033[H\033[2J\033[?25l\033[0m"
}

func cleanup() string {
	return "\033[H\033[2J\033[?25h\033[0m"
}

func linefeed() string {
	return "\033[m\n"
}

func title() string {
	return "\033];Super Mario World in terminal\007"
}
