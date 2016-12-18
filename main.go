package main

import (
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/sdl"

	"gitlab.com/rangerdanger/sdlaudio"
	"gitlab.com/rangerdanger/tetris/tetris"
)

const screenWidth int = 600
const screenHeight int = 400

const fps uint32 = 60
const delayTime uint32 = 1000.0 / fps

const lockDelay int = 31

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Println(err)
		return
	}

	if err := sdlaudio.Init(); err != nil {
		log.Println(err)
		return
	}

	// Create Window
	window, err := sdl.CreateWindow(
		"Tetris The Grand Master - Clone",
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
	foo := tetris.NewGame()
	foo.Init()

	// Main Loop
	running := true
	for running {
		frameStart := sdl.GetTicks()

		//foo.BufferCommand(0)
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.KeyDownEvent:
				switch t.Keysym.Sym {
				case sdl.K_LEFT:
					foo.BufferCommand(tetris.ShiftLeft)
				case sdl.K_RIGHT:
					foo.BufferCommand(tetris.ShiftRight)
				case sdl.K_UP:
					foo.BufferCommand(tetris.RotateClockwise)
				case sdl.K_DOWN:
					foo.BufferCommand(tetris.RotateCounterClockwise)
				case sdl.K_SPACE:
					foo.BufferCommand(tetris.ManualDrop)
				case sdl.K_RETURN:
					foo.BufferCommand(tetris.Start)
					fmt.Println("STARTING")
				}
			case *sdl.KeyUpEvent:
				foo.BufferCommand(0)
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
	sdlaudio.Quit()
}
