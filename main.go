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

func main() {
	fmt.Println("Starting...")

	window, renderer, err := initSDL()
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	defer renderer.Destroy()

	buffer, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_STREAMING, width, height)
	if err != nil {
		panic(err)
	}
	
	b := brush {
		rect: sdl.Rect{X: 200, Y: 100, W: 10, H: 10},
		color: sdl.Color{R: 255, G: 255, B: 0, A: 255},
	}

	sdl.ShowCursor(0)

	count := 0
	timer := sdl.GetTicks()

	var root point

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
			if root.distance(p) > 2 {
				root = p
				draw(buffer, b)
			}
		}

		// Clear screen
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		// Draw buffer
		renderer.Copy(buffer, nil, nil)

		// Draw rect
		renderer.SetDrawColor(b.color.R, b.color.G, b.color.B, b.color.A)
		renderer.FillRect(&b.rect)

		// Debug stuff
		if DEBUG {
			count++
			if then := sdl.GetTicks(); then - timer > 1000 {
				window.SetTitle("paintsson  FPS: " + strconv.FormatInt(int64(count), 10))
				timer = then
				count = 0
			}
		}

		renderer.Present()
	}
}

func initSDL() (*sdl.Window, *sdl.Renderer, error) {
	err := sdl.Init(sdl.INIT_EVERYTHING)

	window, err := sdl.CreateWindow("paintsson", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_OPENGL)

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)

	return window, renderer, err
}