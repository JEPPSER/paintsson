package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"fmt"
	"strconv"
	"io/ioutil"
	"strings"
)

var width int32
var height int32
var DEBUG bool

// Colors
var colors map[string]sdl.Color

type textfield struct {
	command string
	surface *sdl.Surface
	texture *sdl.Texture
	font *ttf.Font
}

func main() {
	fmt.Println("Starting...")

	width = 1000
	height = 700
	DEBUG = true

	window, renderer := initSDL()
	defer window.Destroy()
	defer renderer.Destroy()

	buffer, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_STREAMING, width, height)
	if err != nil { panic(err) }
	defer buffer.Destroy()

	var textSurface *sdl.Surface
	var textTexture *sdl.Texture
	font, err := ttf.OpenFont("fonts/CONSOLAB.ttf", 26)
	if err != nil { panic(err) }
	
	textfield := &textfield {
		command: "",
		surface: textSurface,
		texture: textTexture,
		font: font,
	}

	initColors()
	clearBuffer(buffer, colors["chalkboard"])
	updateTextfield(renderer, textfield)

	b := &brush {
		rect: sdl.Rect{X: 200, Y: 100, W: 5, H: 5},
		color: colors["white"],
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
			case *sdl.WindowEvent:
				e := event.(*sdl.WindowEvent)
				if e.Event == sdl.WINDOWEVENT_SIZE_CHANGED {
					oldWidth := width;
					oldHeight := height;
					width = e.Data1
					height = e.Data2

					oldPixels, _, _ := buffer.Lock(nil)
					buffer.Unlock()

					buffer, _ = renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_STREAMING, width, height)
					pastePixels(buffer, oldPixels, int(oldWidth), int(oldHeight))
				}
			case *sdl.QuitEvent:
				return
			case *sdl.KeyboardEvent:
				keyboardPressed(event.(*sdl.KeyboardEvent), renderer, buffer, b, textfield)
			}
		}

		var l *line

		// Mouse input
		mX, mY, button := sdl.GetMouseState()
		b.rect.X = mX
		b.rect.Y = mY
		if button == 1 {
			p := point{x: mX, y: mY}

			if !down {
				root = p
				l = &line{root, p}
			} else if p.distance(root) > 4 {
				l = &line{root, p}
				root = p
			}
			
			down = true
		} else {
			down = false
		}

		// Rendering
		renderBuffer(renderer, buffer, b, l)
		renderTextfield(renderer, textfield, window)
		renderer.Present()

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

func keyboardPressed(e *sdl.KeyboardEvent, renderer *sdl.Renderer, buffer *sdl.Texture, b *brush, tf *textfield) {
	if e.Type == sdl.KEYDOWN {
		k := e.Keysym.Sym

		// Printable characters
		if 		k <= 122 && k >= 97 ||
				k <= 57 && k >= 48 ||
				k == 32 ||
				k == 44 ||
				k == 45 {
			tf.command += string(k)
		}

		// Execute command
		if k == 13 {
			parseCommand(tf.command, buffer, b)
			tf.command = ""
		}

		// Backspace
		if k == 8 && len(tf.command) > 0 {
			tf.command = tf.command[:len(tf.command) - 1]
		}

		updateTextfield(renderer, tf)
	}
}

func updateTextfield(renderer *sdl.Renderer, tf *textfield) {
	var err error
	tf.surface, err = tf.font.RenderUTF8Solid(">" + tf.command, colors["black"])
	tf.texture, err = renderer.CreateTextureFromSurface(tf.surface)
	if err != nil { panic(err) }
}

func renderBuffer(renderer *sdl.Renderer, buffer *sdl.Texture, b *brush, l *line) {
	if l != nil {
		drawLine(buffer, b, l.from, l.to)
	}

	// Draw buffer
	renderer.Copy(buffer, nil, nil)

	// Draw cursor
	renderer.SetDrawColor(b.color.R, b.color.G, b.color.B, b.color.A)
	if b.brushType == 0 {
		renderer.FillRect(&b.rect)
	} else if b.brushType == 1 {
		drawCircle(renderer, b, point{b.rect.X, b.rect.Y})
	}
}

func renderTextfield(renderer *sdl.Renderer, tf *textfield, window *sdl.Window) {
	renderer.SetDrawColor(220, 220, 220, 255)
	renderer.FillRect(&sdl.Rect{X: 0, Y: height - 30, W: width, H: 30})

	renderer.SetDrawColor(170, 170, 170, 255)
	renderer.FillRect(&sdl.Rect{X: 0, Y: height - 33, W: width, H: 3})

	w, h, _ := tf.font.SizeUTF8(">" + tf.command)
	renderer.Copy(tf.texture, nil, &sdl.Rect{X: 5, Y: height - 27, W: int32(w), H: int32(h)})
}

func initSDL() (*sdl.Window, *sdl.Renderer) {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil { panic(err) }

	window, err := sdl.CreateWindow("paintsson", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_OPENGL)
	if err != nil { panic(err) }
	window.SetResizable(true)

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil { panic(err) }

	if err := ttf.Init(); err != nil {
		panic(err)
	}

	return window, renderer
}

func initColors() {
	colors = make(map[string]sdl.Color)
	dat, err := ioutil.ReadFile("colors.cfg")
	if err != nil { panic(err) }
	str := strings.ReplaceAll(string(dat), "\r", "")
	parseColors(str)
}