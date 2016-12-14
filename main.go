package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"gitlab.com/rangerdanger/tetris/tetris"
)

const screenWidth int = 640
const screenHeight int = 480

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

	level := 0
	gm := tetris.GetTGMGravityMap()
	g, gCounter := gm[level]/256, 0.0

	var foo tetris.Game
	foo.Start()

	//lockFrames := 0
	lockedPieces := make([]sdl.Rect, foo.Board().Area())

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

		// Gravity update
		gCounter += g
		if gCounter >= 1.0 {
			foo.Drop()
			gCounter = 0.0
		}

		renderer.SetDrawColor(0, 128, 255, 255)
		renderer.Clear()

		foo.Board().Draw(renderer)

		// Lock
		// testPiece := activePiece
		// testPiece.Drop()
		// collision := false
		// for _, t := range testPiece.Blocks() {
		// 	for _, g := range board.Ground() {
		// 		if t.HasIntersection(&g) {
		// 			collision = true
		// 			break
		// 		}
		// 	}

		// 	if collision {
		// 		break
		// 	}
		// }

		// if collision {
		// 	lockFrames++
		// } else {
		// 	lockFrames = 0
		// }

		// // Lock and get next piece
		// if lockFrames >= lockDelay {
		// 	lockedPieces = append(lockedPieces, activePiece.Blocks()...)

		// 	for x := 1; x < board.Width(); x++ {
		// 		for y := x; y < board.Area(); y += board.Width() {
		// 			c := false
		// 			for j := 0; j < len(lockedPieces)-1; j++ {
		// 				if board.Cells()[y].HasIntersection(&lockedPieces[j]) {
		// 					//ground[x] = lockedPieces[j]
		// 					c = true
		// 					break
		// 				}
		// 			}

		// 			if c {
		// 				break
		// 			}
		// 		}
		// 	}

		// 	game.NextTetromino()
		// }

		// Draw tetrominos
		foo.ActivePiece().Draw(renderer)
		//nextPiece.Draw(renderer)
		renderer.FillRects(lockedPieces)

		renderer.Present()

		if frameTime := sdl.GetTicks() - frameStart; frameTime < delayTime {
			sdl.Delay(delayTime - frameTime)
		}
	}

	// Clean Up
	sdl.Quit()
}
