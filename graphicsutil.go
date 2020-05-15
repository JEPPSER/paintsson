package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

type brush struct {
	rect sdl.Rect
	color sdl.Color
	clearColor sdl.Color
	brushType int
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

func pastePixels(buffer *sdl.Texture, paste []byte, oldWidth int, oldHeight int) {
	pixels, _, err := buffer.Lock(nil)
	if err != nil { panic(err) }
	for x := 0; x < oldWidth; x++ {
		for y := 0; y < oldHeight; y++ {
			if int32(x) >= width || int32(y) >= height {
				continue
			}
			oldIndex := (oldWidth * y + x) * 4
			newIndex := (int(width) * y + x) * 4
			if int(newIndex + 3) > len(pixels) { continue }
			pixels[newIndex] = paste[oldIndex]
			pixels[newIndex + 1] = paste[oldIndex + 1]
			pixels[newIndex + 2] = paste[oldIndex + 2]
			pixels[newIndex + 3] = paste[oldIndex + 3]
		}
	}
	buffer.Unlock()
}

func clearBuffer(buffer *sdl.Texture, b *brush) {
	pixels, _, err := buffer.Lock(nil)
	if err != nil { panic(err) }
	for i := 0; i < len(pixels); i+=4 {
		pixels[i] = byte(b.clearColor.A)
		pixels[i + 1] = byte(b.clearColor.B)
		pixels[i + 2] = byte(b.clearColor.G)
		pixels[i + 3] = byte(b.clearColor.R)
	}
	buffer.Unlock()
}

func drawLine(buffer *sdl.Texture, b *brush, from point, to point) {
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
	kx *= float64(b.rect.W) * 0.1;
	ky *= float64(b.rect.H) * 0.1;

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

func drawMultiple(buffer *sdl.Texture, b *brush, list []point) {
	pixels, _, err := buffer.Lock(nil)
	if err != nil { panic(err) }
	for i := 0; i < len(list); i++ {
		if b.brushType == 0 {
			fillRect(pixels, b, list[i])
		} else if b.brushType == 1 {
			fillCircle(pixels, b, list[i])
		}
	}
	buffer.Unlock()
}

func draw(buffer *sdl.Texture, b *brush, p point) {
	pixels, _, err := buffer.Lock(nil)
	if err != nil { panic(err) }
	fillRect(pixels, b, p)
	buffer.Unlock()
}

func fillCircle(pixels []byte, b *brush, p point) {
	origo := float64(b.rect.H) / 2
	y := 0
	if b.rect.H > 1 { y = 1 }
	for ; y < int(b.rect.H); y++ {
		cY := origo - float64(y)
		cX := math.Sqrt(math.Pow(origo, 2) - math.Pow(cY, 2))
		fromX := int(p.x) + int(math.Round(cX * -1 + origo))
		toX := int(p.x) + int(math.Round(cX + origo))
		fillRow(pixels, b, fromX, toX, y + int(p.y))
	}
}

func drawCircle(renderer *sdl.Renderer, b *brush, p point) {
	origo := float64(b.rect.H) / 2
	y := 0
	if b.rect.H > 1 { y = 1 }
	for ; y < int(b.rect.H); y++ {
		cY := origo - float64(y)
		cX := math.Sqrt(math.Pow(origo, 2) - math.Pow(cY, 2))
		fromX := int(p.x) + int(math.Round(cX * -1 + origo))
		toX := int(p.x) + int(math.Round(cX + origo))
		
		for x := fromX; x <= toX; x++ {
			renderer.DrawPoint(int32(x), p.y + int32(y))
		}
	}
}

func fillRow(pixels []byte, b *brush, fromX int, toX int, y int) {
	for x := fromX; x <= toX; x++ {
		index := (int(width) * y + x) * 4
		if int(index + 3) > len(pixels) { continue }
		pixels[index] = byte(b.color.A)
		pixels[index + 1] = byte(b.color.B)
		pixels[index + 2] = byte(b.color.G)
		pixels[index + 3] = byte(b.color.R)
	}
}

func fillRect(pixels []byte, b *brush, p point) {
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