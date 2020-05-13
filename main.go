package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"fmt"
	"strconv"
)

const (
	DEBUG bool = true
	width int32 = 800
	height int32 = 600
)

var lines []line
var buffer *sdl.Texture
var renderer *sdl.Renderer
var window *sdl.Window
var b brush

func main() {
	fmt.Println("Starting...")

	window, renderer = initSDL()
	defer window.Destroy()
	defer renderer.Destroy()

	buffer, _ = renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_STREAMING, width, height)

	b = brush {
		rect: sdl.Rect{X: 200, Y: 100, W: 5, H: 5},
		color: sdl.Color{R: 255, G: 255, B: 0, A: 255},
	}

	sdl.ShowCursor(0)

	count := 0
	timer := sdl.GetTicks()

	var root point
	var down bool = false

	for {
		// Poll events
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		// Mouse input
		mX, mY, button := sdl.GetMouseState()
		b.rect.X = mX
		b.rect.Y = mY
		if button == 1 {
			p := point{x: mX, y: mY}

			if !down {
				root = p
				lines = append(lines, line{root, p})
			} else {
				lines = append(lines, line{root, p})
				root = p
			}
			
			down = true
		} else {
			down = false
		}

		render()

		// Debug stuff
		if DEBUG {
			count++
			if then := sdl.GetTicks(); then - timer > 1000 {
				window.SetTitle("paintsson  FPS: " + strconv.FormatInt(int64(count), 10))
				timer = then
				count = 0
			}
		}
	}
}

func render() {
	length := len(lines)
	for i := 0; i < length; i++ {
		drawLine(buffer, b, lines[0].from, lines[0].to)
		lines = lines[1:]
	}

	// Clear screen
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()

	// Draw buffer
	renderer.Copy(buffer, nil, nil)

	// Draw rect
	renderer.SetDrawColor(b.color.R, b.color.G, b.color.B, b.color.A)
	renderer.FillRect(&b.rect)

	renderer.Present()
}

func initSDL() (*sdl.Window, *sdl.Renderer) {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil { panic(err) }

	window, err := sdl.CreateWindow("paintsson", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_OPENGL)
	if err != nil { panic(err) }

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil { panic(err) }

	return window, renderer
}