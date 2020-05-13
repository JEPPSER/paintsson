package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

type brush struct {
	rect sdl.Rect
	color sdl.Color
}

type point struct {
	x int32
	y int32
}

func (p point) distance(other point) float64 {
	return math.Hypot(float64(p.x - other.x), float64(p.y - other.y))
}

func drawLine(buffer *sdl.Texture, b brush, from point, to point) {
	
}

func draw(buffer *sdl.Texture, b brush) {
	pixels, _, _ := buffer.Lock(nil)
	for x := b.rect.X; x < b.rect.X + b.rect.W; x++ {
		for y := b.rect.Y; y < b.rect.Y + b.rect.H; y++ {
			index := (width * y + x) * 4
			if int(index + 3) > len(pixels) { continue }
			pixels[index] = byte(b.color.A)
			pixels[index + 1] = byte(b.color.B)
			pixels[index + 2] = byte(b.color.G)
			pixels[index + 3] = byte(b.color.R)
		}
	}
	buffer.Unlock()
}