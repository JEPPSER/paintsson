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

type line struct {
	from point
	to point
}

func (p point) distance(other point) float64 {
	return math.Hypot(float64(p.x - other.x), float64(p.y - other.y))
}

func drawLine(buffer *sdl.Texture, b brush, from point, to point) {
	if from.x > to.x {
		tempTo := point{from.x, from.y}
		tempFrom := point{to.x, to.y}
		to = tempTo
		from = tempFrom
	}

	var k float64
	if to.x - from.x == 0 {
		k = 10000;
		if to.y - from.y < 0 {
			k = k * -1
		}
	} else {
		k = float64(to.y - from.y) / float64(to.x - from.x);
	}
	
	var kx, ky float64

	// Calculating how much x and y will increase each iteration of drawing the line.
	if k < 0 {
		kx = 1 / (-1 + k) * -1
		ky = k / (-1 + k) * -1
	} else {
		kx = 1 / (1 + k)
		ky = k / (1 + k)
	}
	kx *= float64(b.rect.W) * 0.75;
	ky *= float64(b.rect.H) * 0.75;

	totLength := from.distance(to)
	length := 0.0
	increment := math.Sqrt(kx * kx + ky * ky)

	for i := 0; length < totLength + 0.5; i++ {
		x := float64(from.x) + kx * float64(i)
		y := float64(from.y) + ky * float64(i)
		draw(buffer, b, point{int32(x), int32(y)})
		length += increment
	}
}

func draw(buffer *sdl.Texture, b brush, p point) {
	pixels, _, err := buffer.Lock(nil)
	if err != nil { panic(err) }
	for x := p.x; x < p.x + b.rect.W; x++ {
		for y := p.y; y < p.y + b.rect.H; y++ {
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