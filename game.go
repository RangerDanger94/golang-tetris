package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"gitlab.com/rangerdanger/tetris/tetris"
)

const screenWidth int = 640
const screenHeight int = 480

const fps uint32 = 60
const delayTime uint32 = 1000.0 / fps

// Game Board
const gXLength int32 = 10
const gYLength int32 = 24
const gSize int32 = 20
const gStartX int32 = int32(screenWidth)/2 - (gXLength*gSize)/2
const gStartY int32 = 0 //int32(screenHeight)/2 - (gYLength*gSize)/2

func createBoard() []sdl.Rect {
	board := make([]sdl.Rect, gXLength*gYLength)
	var x, y int32 = gStartX, gStartY

	for i := 1; i <= int(gXLength*gYLength); i++ {
		board[i-1] = sdl.Rect{X: x, Y: y, W: gSize, H: gSize}
		x += gSize

		if i%int(gXLength) == 0 {
			y += gSize
			x = gStartX
		}
	}

	return board
}

// Assumes rects fill X and then Y [][][] ->
func createGround(b []sdl.Rect) []sdl.Rect {
	ground := make([]sdl.Rect, gXLength)
	for i, c := gXLength*gYLength-gXLength, 0; i < gXLength*gYLength; i++ {
		ground[c].X, ground[c].Y = b[i].X, b[i].Y+gSize
		ground[c].W, ground[c].H = gSize, gSize*5
		c++
	}

	return ground
}

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

	tetris.ResetTGMRandomizer()
	b := createBoard()
	ground := createGround(b)
	const lockDelay int = 31
	lockFrames := 0
	var lockedPieces []sdl.Rect = make([]sdl.Rect, gXLength*gYLength)
	activePiece, nextPiece := tetris.NextTGMRandomizer(), tetris.NextTGMRandomizer()
	activePiece.Resize(gSize)
	tetris.SpawnTetromino(b, &activePiece)

	// Main Loop
	running := true
	for running {
		frameStart := sdl.GetTicks()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.KeyDownEvent:
				switch t.Keysym.Sym {
				case sdl.K_LEFT:
					activePiece.ShiftLeft()
				case sdl.K_RIGHT:
					activePiece.ShiftRight()
				case sdl.K_UP:
					activePiece.RotateClockwise()
				case sdl.K_DOWN:
					activePiece.RotateCounterClockwise()
				case sdl.K_SPACE:
					activePiece.Drop()
				}
			case *sdl.QuitEvent:
				running = false
			}
		}

		// Gravity update
		gCounter += g
		if gCounter >= 1.0 {
			activePiece.Drop()
			gCounter = 0.0
		}

		renderer.SetDrawColor(0, 128, 255, 255)
		renderer.Clear()

		// Draw grid
		renderer.SetDrawColor(0x0, 0x0, 0x0, 0xFF)
		renderer.FillRects(b)

		// // Lock
		testPiece := activePiece
		testPiece.Drop()
		collision := false
		for _, t := range testPiece.Blocks() {
			for _, g := range ground {
				if t.HasIntersection(&g) {
					collision = true
					break
				}
			}

			if collision {
				break
			}
		}

		if collision {
			lockFrames++
		} else {
			lockFrames = 0
		}

		// Lock and get next piece
		if lockFrames >= lockDelay {
			lockedPieces = append(lockedPieces, activePiece.Blocks()...)
		
			for x := 1; x < int(gXLength); x++ {
				for y := x; y < int(gXLength*gYLength); y += int(gXLength) {
					//fmt.Printf("X:%v, Y:%v\n", x, y)
					c := false
					for j := 0; j < len(lockedPieces) - 1; j++ {
						if b[y].HasIntersection(&lockedPieces[j]) {
							ground[x] = lockedPieces[j]
							c = true
							break
						}
					}
					if c {
						break
					}
				}
			}

			for i, v := range ground {
				fmt.Printf("i: %v, Y:%v\n", i, v.Y)
			}

			activePiece = nextPiece
			nextPiece = tetris.NextTGMRandomizer()
			activePiece.Resize(gSize)
			tetris.SpawnTetromino(b, &activePiece)
		}

		// Draw tetrominos
		activePiece.Draw(renderer)
		nextPiece.Draw(renderer)
		renderer.FillRects(lockedPieces)

		renderer.Present()

		if frameTime := sdl.GetTicks() - frameStart; frameTime < delayTime {
			sdl.Delay(delayTime - frameTime)
		}
	}

	// Clean Up
	sdl.Quit()
}
