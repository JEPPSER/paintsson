package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
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
var command string
var font *ttf.Font

// Colors
var colors map[string]sdl.Color

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

	initColors()

	var err error
	font, err = ttf.OpenFont("fonts/CONSOLAB.ttf", 26)
	if err != nil { panic(err) }

	clearBuffer(buffer, colors["black"])

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
			case *sdl.KeyboardEvent:
				e := event.(*sdl.KeyboardEvent)
				keyboardPressed(e)
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

func keyboardPressed(e *sdl.KeyboardEvent) {
	if e.Type == sdl.KEYDOWN {
		k := e.Keysym.Sym

		// Printable characters
		if 		k <= 122 && k >= 97 ||
				k <= 57 && k >= 48 ||
				k == 32 ||
				k == 44 ||
				k == 45 {
			command += string(k)
		}

		// Execute command
		if k == 13 {
			parse(command)
			command = ""
		}

		// Backspace
		if k == 8 && len(command) > 0 {
			command = command[:len(command) - 1]
		}
	}
}

func render() {
	length := len(lines)
	for i := 0; i < length; i++ {
		drawLine(buffer, b, lines[0].from, lines[0].to)
		lines = lines[1:]
	}

	// Draw buffer
	renderer.Copy(buffer, nil, nil)

	// Draw cursor
	renderer.SetDrawColor(b.color.R, b.color.G, b.color.B, b.color.A)
	renderer.FillRect(&b.rect)

	// Text field
	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.FillRect(&sdl.Rect{X: 0, Y: height - 30, W: width, H: 30})
	surface, _ := font.RenderUTF8Solid(command, colors["black"])
	tex, _ := renderer.CreateTextureFromSurface(surface)
	w, h, _ := font.SizeUTF8(command)
	renderer.Copy(tex, nil, &sdl.Rect{X: 5, Y: height - 27, W: int32(w), H: int32(h)})

	renderer.Present()
}

func initSDL() (*sdl.Window, *sdl.Renderer) {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil { panic(err) }

	window, err := sdl.CreateWindow("paintsson", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_OPENGL)
	if err != nil { panic(err) }

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil { panic(err) }

	if err := ttf.Init(); err != nil {
		panic(err)
	}

	return window, renderer
}

func initColors() {
	colors = make(map[string]sdl.Color)
	colors["black"] = sdl.Color{R: 0, G: 0, B: 0, A: 255}
	colors["white"] = sdl.Color{R: 255, G: 255, B: 255, A: 255}
	colors["blue"] = sdl.Color{R: 0, G: 0, B: 255, A: 255}
	colors["green"] = sdl.Color{R: 0, G: 255, B: 0, A: 255}
	colors["red"] = sdl.Color{R: 255, G: 0, B: 0, A: 255}
	colors["yellow"] = sdl.Color{R: 255, G: 255, B: 0, A: 255}
}