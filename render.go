package main

import (
	"bufio"
	"fmt"
	"image"
	"os"
	"sync"
	"time"
)

type renderer struct {
	buffer  *bufio.Writer
	width   int
	height  int
	channel chan [][]byte
	chars   map[int]string
	line    []byte
}

func (r *renderer) frame(background [][]byte, charact []string, point image.Point) [][]byte {
	bgRect := image.Rect(0, 0, r.width, r.height).
		Intersect(image.Rect(point.X, point.Y, point.X+len(charact[0]), point.Y+len(charact)))

	// 判断是否有交集
	if bgRect.Dx() == 0 || bgRect.Dy() == 0 {
		return background
	}

	fgRect := image.Rect(
		bgRect.Min.X-point.X,
		bgRect.Min.Y-point.Y,
		bgRect.Max.X-point.X,
		bgRect.Max.Y-point.Y,
	)

	for i := 0; i < bgRect.Dy(); i++ {
		for j := 0; j < bgRect.Dx(); j++ {
			// 是背景色就跳过
			if charact[fgRect.Min.Y+i][fgRect.Min.X+j] != char {
				background[bgRect.Min.Y+i][bgRect.Min.X+j] = charact[fgRect.Min.Y+i][fgRect.Min.X+j]
			}
		}
	}

	// src := charact[fgRect.Min.Y+i][fgRect.Min.X:]
	// dst := background[bgRect.Min.Y+i][bgRect.Min.X:]
	// copy(dst, src)

	return background
}

func (r *renderer) print(frame [][]byte) {
	defer r.buffer.Flush()
	r.buffer.WriteString(reset())

	for i := range frame {
		for _, output := range r.split(frame[i]) {
			r.buffer.WriteString(output)
		}

		if i != len(frame)-1 {
			r.buffer.WriteString(linefeed())
		}
	}
}

func (r *renderer) render(args arguments) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		i := 0
		for {
			background := r.background()
			for _, frame := range args.Frames[i] {
				r.frame(background, ASCII[frame.Name], frame.Point)
			}

			select {
			case r.channel <- background:
			}

			i++
			if i == len(args.Frames) {
				i = 0
				// close(r.Ch)
				// return
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			frame, ok := <-r.channel
			if !ok {
				return
			}
			r.print(frame)
			time.Sleep(time.Duration(args.Itv) * time.Millisecond)
		}
	}()

	wg.Wait()
	fmt.Print(cleanup())
}

func (r *renderer) split(s []byte) []string {
	var ret []string
	var last int
	for i := 0; i < len(s)-1; i++ {
		if s[i] != s[i+1] {
			ret = append(ret, escape(palettes[s[i]], r.chars[i+1-last]))
			last = i + 1
		}
	}

	ret = append(ret, escape(palettes[s[last]], r.chars[len(s)-last]))
	return ret
}

func (r *renderer) background() [][]byte {
	ret := make([][]byte, r.height)
	for i := range ret {
		ret[i] = make([]byte, r.width)
		copy(ret[i], r.line)
	}

	return ret
}

func newRenderer() *renderer {
	buf := bufio.NewWriter(os.Stdout)
	buf.WriteString(prepare())
	buf.WriteString(title())

	width, height := term()

	// 预先创建好各种长度的字符串
	chars := make(map[int]string)
	for i := 1; i <= width; i++ {
		bytes := make([]byte, i)
		for j := 0; j < i; j++ {
			bytes[j] = ' '
		}
		chars[i] = string(bytes)
	}

	// 单行数据，填充'.'
	line := make([]byte, width)
	for i := 0; i < width; i++ {
		line[i] = char
	}

	return &renderer{
		buffer:  buf,
		width:   width,
		height:  height,
		channel: make(chan [][]byte, 2),
		chars:   chars,
		line:    line,
	}
}
