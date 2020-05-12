package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"fmt"
	"strconv"
)

const DEBUG bool = true

func main() {
	fmt.Println("Starting...")

	window, renderer, err := initSDL()
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	defer renderer.Destroy()

	rect := sdl.Rect{X: 200, Y: 100, W: 20, H: 20}

	count := 0
	timer := sdl.GetTicks()

	for {
		// Poll events
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		mX, mY, button := sdl.GetMouseState()
		if button == 1 {
			rect.X = mX
			rect.Y = mY
		}

		// Clear screen
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		// Draw rect
		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.FillRect(&rect)

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

	window, err := sdl.CreateWindow("paintsson", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_OPENGL)

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)

	return window, renderer, err
}