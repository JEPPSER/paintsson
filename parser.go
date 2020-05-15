package main

import (
	"strings"
	"strconv"
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

func parseCommand(str string, buffer *sdl.Texture, b *brush) {
	parts := strings.Split(str, " ")
	if len(parts) == 1 {
		
		if parts[0] == "clear" {
			clearBuffer(buffer, b)
		}
	} else if len(parts) == 2 {

		if parts[0] == "size" {
			size, err := strconv.ParseInt(parts[1], 10, 32)
			if err != nil {
				fmt.Println(err)
				return
			}
			b.rect.W = int32(size)
			b.rect.H = int32(size)
		} else if parts[0] == "color" {
			b.color = colors[parts[1]]
		} else if parts[0] == "clear" {
			b.clearColor = colors[parts[1]]
			clearBuffer(buffer, b)
		} else if parts[0] == "brush" {
			t, err := strconv.ParseInt(parts[1], 10, 32)

			if err != nil { return }
			if t < 0 || t > 1 { return }

			b.brushType = int(t)
		}
	}
}

func parseColors(str string) {
	lines := strings.Split(str, "\n")
	for i := 0; i < len(lines); i++ {
		parts := strings.Split(lines[i], " ")
		if len(parts) != 2 { continue }

		values := strings.Split(parts[1], ",")
		if len(values) != 3 { continue }

		r, err := strconv.ParseUint(values[0], 10, 8)
		if err != nil { continue }
		g, err := strconv.ParseUint(values[1], 10, 8)
		if err != nil { continue }
		b, err := strconv.ParseUint(values[2], 10, 8)
		if err != nil { continue }

		colors[parts[0]] = sdl.Color{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(255)}
	}
}