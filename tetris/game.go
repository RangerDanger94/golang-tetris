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
const gStartX int32 = 640/2 - (int32(gXLength)*gSize)/2
const gStartY int32 = 0

type coloredRect struct {
	color sdl.Color
	rect  sdl.Rect
}

type Game struct {
	activePiece  Tetromino
	nextPiece    Tetromino
	holdPiece    Tetromino
	board        Grid
	lockedPieces []coloredRect

	areFrames             int
	dasFrames             int
	lockDelay, lockFrames int
	clearFrames           int
	gravFrames            float64

	level   int
	gravity float64
}

// ProcessFrame runs the game logic for a frame
func (g *Game) ProcessFrame() {
	g.doGravity()
	g.checkLock()
}

// Draw renders the game using an sdl.Renderer
func (g *Game) Draw(r *sdl.Renderer) {
	g.board.Draw(r)
	g.activePiece.Draw(r)

	for _, v := range g.lockedPieces {
		r.SetDrawColor(v.color.R, v.color.G, v.color.B, v.color.A)
		r.FillRect(&v.rect)
	}
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

	fmt.Printf("Gravity: %v", g.gravity)
}

func (g *Game) checkLock() {
	testPiece := g.activePiece
	testPiece.Drop()

	if g.collision(testPiece) {
		g.lockFrames++
	} else {
		g.lockFrames = 0
	}

	if g.lockFrames >= g.lockDelay {
		locked := make([]coloredRect, len(g.activePiece.Blocks()))
		for i, v := range g.activePiece.Blocks() {
			locked[i].rect, locked[i].color = v, g.activePiece.Color()
		}

		g.lockedPieces = append(g.lockedPieces, locked...)
		g.nextTetromino()
	}
}

// Drop tries to lower the active piece
func (g *Game) Drop() {
	testPiece := g.activePiece
	testPiece.Drop()

	if !g.collision(testPiece) {
		g.activePiece.Drop()
	}
}

func (g Game) collisionRects() []sdl.Rect {
	c := make([]sdl.Rect, len(g.lockedPieces))
	for i, v := range g.lockedPieces {
		c[i] = v.rect
	}
	return c
}

// collision checks if a tetromino is colliding with the following
// 1. Locked pieces
// 2. Board edges
func (g *Game) collision(t Tetromino) bool {
	hit := false

	// Check tetromino outside the board
	for _, v := range t.Blocks() {
		if v.X < g.board.X() || v.X >= g.board.X()+g.board.PixelWidth() {
			hit = true
			break
		} else if v.Y < g.board.Y() || v.Y >= g.board.Y()+g.board.PixelHeight() {
			hit = true
			break
		}
	}

	// Check collision with any locked pieces
	for _, i := range t.Blocks() {
		for _, j := range g.collisionRects() {
			if i.HasIntersection(&j) {
				hit = true
				break
			}
		}

		if hit {
			break
		}
	}

	return hit
}

func (g *Game) Start() {
	ResetTGMRandomizer()
	g.level = 0
	g.areFrames = 30
	g.dasFrames = 14
	g.lockDelay, g.lockFrames = 30, 0
	g.clearFrames = 41

	g.board = NewGrid(gStartX, gStartY, gSize, gXLength, gYLength)
	g.lockedPieces = make([]coloredRect, g.board.Area())
	g.activePiece, g.nextPiece = NextTGMRandomizer(), NextTGMRandomizer()
	g.SpawnTetromino(&g.activePiece)
}

// gets the next tetromino from the randomizer
func (g *Game) nextTetromino() {
	g.activePiece = g.nextPiece
	g.nextPiece = NextTGMRandomizer()
	g.SpawnTetromino(&g.activePiece)
}

// SpawnTetromino on the grid
func (g *Game) SpawnTetromino(t *Tetromino) {
	t.Resize(g.board.cellSize)
	t.move(g.board.Cells()[3].X, g.board.Cells()[3].Y-g.board.CellSize())
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

var tgmBagHistory = []int32{Z, Z, S, S}
var tgmFirstPiece bool = true

// ResetTGMRandomizer seeds generator and resets bag history
func ResetTGMRandomizer() {
	rand.Seed(time.Now().UTC().UnixNano())
	tgmBagHistory = []int32{Z, Z, S, S}
	tgmFirstPiece = true
}

// NextTGMRandomizer gets the next Tetromino according to tgm randomization rules
func NextTGMRandomizer() Tetromino {
	tS := rand.Int31n(tetronimos)

	// The game never deals an S, Z or O as the first piece
	if tgmFirstPiece {
		for tS == S || tS == Z || tS == O {
			tS = rand.Int31n(tetronimos)
		}
	}

	// Attempt to get a tetronimo not in the bag history
	for _, t := range tgmBagHistory {
		if tS == t {
			tS = rand.Int31n(tetronimos)
		}
	}

	for i := len(tgmBagHistory) - 1; i > 0; i-- {
		tgmBagHistory[i] = tgmBagHistory[i-1]
	}
	tgmBagHistory = append([]int32{tS}, tgmBagHistory[1:4]...)
	tgmFirstPiece = false

	fmt.Printf("Termino: %v\n", tS)
	return generateTetronimo(tS)
}

// GetTGMGravityMap get tgm grav rules
func GetTGMGravityMap() map[int]float64 {
	return tgmGravity
}
