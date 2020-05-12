package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Println("sdl: ", err)
		return
	}

	window, err := sdl.CreateWindow("paintsson", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_OPENGL)
	if err != nil {
		fmt.Println("window: ", err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println("renderer: ", err)
		return
	}
	defer renderer.Destroy()

	rect := sdl.Rect{X: 200, Y: 100, W: 10, H: 10}

	for {
		// Poll events
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		// Clear screen
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		// Draw rect
		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.FillRect(&rect)

		// Present
		renderer.Present()
	}
}