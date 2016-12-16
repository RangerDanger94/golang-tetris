package tetris

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const tetronimos int32 = 7

// Board presets
const gXLength int = 10
const gYLength int = 20
const gSize int32 = 20
const gStartX int32 = 600/2 - (int32(gXLength)*gSize)/2
const gStartY int32 = 0

const (
	s_LOCK = iota
	s_CLEAR
)

type Game struct {
	activePiece Tetromino
	nextPiece   Tetromino
	holdPiece   Tetromino
	board       Grid

	areDelay, areFrames   int
	dasFrames             int
	lockDelay, lockFrames int
	clearFrames           int
	gravFrames            float64

	level   int
	score   int
	gravity float64
	step    int
}

// ProcessFrame runs the game logic for a frame
func (g *Game) ProcessFrame() {
	g.doGravity()

	// Run delays
	switch g.step {
	case s_LOCK:
		if g.checkLock() {
			g.step++
		}
	case s_CLEAR:
		if g.checkClear() {
			g.step = s_LOCK
		}
	}
}

// Draw renders the game using an sdl.Renderer
func (g Game) Draw(r *sdl.Renderer) {
	g.board.Draw(r)
	g.activePiece.Draw(r)
}

func (g *Game) doGravity() {
	// Get current gravity according to game level
	for i, v := range tgmGravity {
		if i > g.level {
			break
		}
		g.gravity = v / 256
	}

	g.gravFrames += g.gravity

	for g.gravFrames >= 1 {
		g.Drop()
		g.gravFrames--
	}
}

func (g *Game) checkLock() bool {
	testPiece := g.activePiece
	testPiece.Drop()

	if g.collision(testPiece) {
		g.lockFrames++
	} else {
		g.lockFrames = 0
	}

	if g.lockFrames >= g.lockDelay {

		for _, v := range g.activePiece.blocks {
			for _, row := range g.board.cells {
				for col := range row {
					if v.Equals(&row[col].rect) {
						row[col].occupied = true
						row[col].color = g.activePiece.color
					}
				}
			}
		}

		g.nextTetromino()
		fmt.Println("Finished locking the active tetromino")
		return true
	}

	return false
}

// check all rows for successful line clear
func (g *Game) checkClear() bool {
	for row := range g.board.cells {
		if g.checkLineClear(row) {
			g.clearLine(row)
		}
	}

	return true
}

// return true if all cells in the row are occupied
func (g *Game) checkLineClear(row int) bool {
	for col := range g.board.cells[row] {
		if !g.board.cells[row][col].occupied {
			return false
		}
	}

	return true
}

// dro
func (g *Game) clearLine(row int) {
	g.board.cells[row] = g.board.createRow(row)

	// Drop all occupied spaces by 1
	for i := row; i > 0; i-- {
		for col := range g.board.cells[i] {
			g.board.cells[i][col].occupied = g.board.cells[i-1][col].occupied
			g.board.cells[i][col].color = g.board.cells[i-1][col].color
		}
	}
}

func remove(s []cell, i int) []cell {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

// Drop tries to lower the active piece
func (g *Game) Drop() {
	testPiece := g.activePiece
	testPiece.Drop()

	if !g.collision(testPiece) {
		g.activePiece.Drop()
	}
}

// collisionRects returns the sdl.Rect elements from lockedPieces
func (g Game) terminoIntersection(t Tetromino, r sdl.Rect) bool {
	for _, v := range t.blocks {
		if v.HasIntersection(&r) {
			return true
		}
	}

	return false
}

// collision checks if a tetromino is colliding with the following
// 1. Locked pieces
// 2. Board edges
func (g *Game) collision(t Tetromino) bool {
	// Check tetromino outside the board
	for _, v := range t.blocks {
		if v.X < g.board.X() || v.X >= g.board.X()+g.board.PixelWidth() {
			return true
		} else if v.Y >= g.board.Y()+g.board.PixelHeight() {
			return true
		}
	}

	// Check collision with any occupied spaces
	for _, row := range g.board.cells {
		for _, col := range row {
			if col.occupied && g.terminoIntersection(t, col.rect) {
				return true
			}
		}
	}

	return false
}

// Start initalizes game
func (g *Game) Start() {
	for i := I; i <= Z; i++ {
		fmt.Println(i)
	}
	ResetTGMRandomizer()
	g.level = 0
	g.areDelay, g.areFrames = 30, 0
	g.dasFrames = 14
	g.lockDelay, g.lockFrames = 30, 0
	g.clearFrames = 41
	g.step = s_LOCK

	g.board = NewGrid(gStartX, gStartY, gSize, gXLength, gYLength)
	g.activePiece, g.nextPiece = NextTGMRandomizer(), NextTGMRandomizer()
	g.SpawnTetromino(&g.activePiece)
}

// gets the next tetromino from the randomizer
func (g *Game) nextTetromino() {
	// g.areFrames++

	// if g.areFrames >= g.areDelay {
	g.areFrames = 0
	g.activePiece, g.nextPiece = g.nextPiece, NextTGMRandomizer()
	g.SpawnTetromino(&g.activePiece)
	//}
}

// SpawnTetromino on the grid
func (g *Game) SpawnTetromino(t *Tetromino) {
	t.Resize(g.board.cellSize)
	t.move(g.board.cells[0][3].rect.X, g.board.cells[0][3].rect.Y-g.board.CellSize())
	g.level++
}

// BufferShift buffers the shift command
func (g *Game) BufferShift(right bool) {
	testPiece := g.activePiece
	if right {
		testPiece.ShiftRight()
	} else {
		testPiece.ShiftLeft()
	}

	if !g.collision(testPiece) {
		if right {
			g.activePiece.ShiftRight()
		} else {
			g.activePiece.ShiftLeft()
		}
	}
}

// BufferRotate buffers the rotation command
func (g *Game) BufferRotate(clockwise bool) {
	rotate := func() {
		if clockwise {
			g.activePiece.RotateClockwise()
		} else {
			g.activePiece.RotateCounterClockwise()
		}
	}

	testRotation := func(t Tetromino) bool {
		if clockwise {
			t.RotateClockwise()
		} else {
			t.RotateCounterClockwise()
		}

		return !g.collision(t)
	}

	if testRotation(g.activePiece) {
		rotate()
	} else {
		testPiece := g.activePiece
		testPiece.ShiftRight()
		if testRotation(testPiece) {
			g.activePiece.ShiftRight()
			rotate()
		} else {
			testPiece = g.activePiece
			testPiece.ShiftLeft()
			if testRotation(testPiece) {
				g.activePiece.ShiftLeft()
				rotate()
			}
		}
	}

}

var tgmBagHistory []int32
var tgmFirstPiece bool

// ResetTGMRandomizer seeds generator and resets bag history
func ResetTGMRandomizer() {
	rand.Seed(time.Now().UTC().UnixNano())
	tgmBagHistory = []int32{Z, Z, Z, Z}
	tgmFirstPiece = true
}

// NextTGMRandomizer gets the next Tetromino according to tgm randomization rules
func NextTGMRandomizer() Tetromino {
	tS := rand.Int31n(tetronimos - 1)

	// The game never deals an S, Z or O as the first piece
	if tgmFirstPiece {
		for tS == S || tS == Z || tS == O {
			tS = rand.Int31n(tetronimos - 1)
		}
	}

	// Attempt to get a tetronimo not in the bag history
	for _, t := range tgmBagHistory {
		if tS == t {
			tS = rand.Int31n(tetronimos - 1)
		}
	}

	for i := len(tgmBagHistory) - 1; i > 0; i-- {
		tgmBagHistory[i] = tgmBagHistory[i-1]
	}
	tgmBagHistory = append([]int32{tS}, tgmBagHistory[1:4]...)
	tgmFirstPiece = false

	return generateTetronimo(tS)
}

var tgmGravity = map[int]float64{
	0:   4.0,
	30:  6.0,
	35:  8.0,
	40:  10.0,
	50:  12.0,
	60:  16.0,
	70:  32.0,
	80:  48.0,
	90:  64.0,
	100: 80.0,
	120: 96.0,
	140: 112.0,
	160: 128.0,
	170: 144.0,
	200: 4.0,
	220: 32.0,
	230: 64.0,
	233: 96.0,
	236: 128.0,
	239: 160.0,
	243: 192.0,
	247: 224.0,
	251: 256.0, // 1G
	300: 512.0,
	330: 768.0,
	360: 1024.0,
	400: 1280.0,
	420: 1024.0,
	450: 768.0,
	500: 5120.0, // 20G
}

// GetTGMGravityMap get tgm grav rules
func GetTGMGravityMap() map[int]float64 {
	return tgmGravity
}
