package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"gitlab.com/rangerdanger/tetris/tetris"
)

const screenWidth int = 600
const screenHeight int = 400

const fps uint32 = 60
const delayTime uint32 = 1000.0 / fps

const lockDelay int = 31

func main() {
	sdl.Init(sdl.INIT_EVERYTHING)

	// Create Window
	window, err := sdl.CreateWindow(
		"Tetris xTreme 2016",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		screenWidth, screenHeight, sdl.WINDOW_SHOWN)

	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	// Create Renderer
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()
	var foo tetris.Game
	foo.Start()

	// Main Loop
	running := true
	for running {
		frameStart := sdl.GetTicks()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.KeyDownEvent:
				switch t.Keysym.Sym {
				case sdl.K_LEFT:
					foo.BufferShift(false)
				case sdl.K_RIGHT:
					foo.BufferShift(true)
				case sdl.K_UP:
					foo.BufferRotate(true)
				case sdl.K_DOWN:
					foo.BufferRotate(false)
				case sdl.K_SPACE:
					foo.Drop()
				}
			case *sdl.QuitEvent:
				running = false
			}
		}

		foo.ProcessFrame()

		renderer.SetDrawColor(0, 128, 255, 255)
		renderer.Clear()

		// Draw game
		foo.Draw(renderer)

		renderer.Present()

		if frameTime := sdl.GetTicks() - frameStart; frameTime < delayTime {
			sdl.Delay(delayTime - frameTime)
		}
	}

	// Clean Up
	sdl.Quit()
}
