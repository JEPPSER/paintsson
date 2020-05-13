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

func clearBuffer(buffer *sdl.Texture, color sdl.Color) {
	pixels, _, err := buffer.Lock(nil)
	if err != nil { panic(err) }
	for i := 0; i < len(pixels); i+=4 {
		pixels[i] = byte(color.A)
		pixels[i + 1] = byte(color.B)
		pixels[i + 2] = byte(color.G)
		pixels[i + 3] = byte(color.R)
	}
	buffer.Unlock()
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

	var list []point

	for i := 0; length < totLength + 0.5; i++ {
		x := float64(from.x) + kx * float64(i)
		y := float64(from.y) + ky * float64(i)
		list = append(list, point{int32(x), int32(y)})
		length += increment
	}

	drawMultiple(buffer, b, list)
}

func drawMultiple(buffer *sdl.Texture, b brush, list []point) {
	pixels, _, err := buffer.Lock(nil)
	if err != nil { panic(err) }
	for i := 0; i < len(list); i++ {
		fillRect(pixels, b, list[i])
	}
	buffer.Unlock()
}

func draw(buffer *sdl.Texture, b brush, p point) {
	pixels, _, err := buffer.Lock(nil)
	if err != nil { panic(err) }
	fillRect(pixels, b, p)
	buffer.Unlock()
}

func fillRect(pixels []byte, b brush, p point) {
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
}