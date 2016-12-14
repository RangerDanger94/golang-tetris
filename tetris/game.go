package tetris

import (
	"fmt"
	"math/rand"
	"time"
)

const tetronimos int32 = 7

// Board presets
const gXLength int = 10
const gYLength int = 20
const gSize int32 = 20
const gStartX int32 = 640/2 - (int32(gXLength)*gSize)/2
const gStartY int32 = 0

type Game struct {
	activePiece Tetromino
	nextPiece   Tetromino
	holdPiece   Tetromino
	board       Grid

	areFrames   int
	dasFrames   int
	lockFrames  int
	clearFrames int
	gravFrames  float64

	level   int
	gravity float64
}

func (g *Game) ProcessFrame() {
	g.doGravity()
}

func (g *Game) doGravity() {
	// Get current gravity according to game level
	for i, v := range tgmGravity {
		if i > g.level {
			break
		}
		g.gravity = v
	}

	g.gravFrames += g.gravity

	for g.gravFrames >= 1 {
		g.activePiece.Drop()
		g.gravFrames--
	}
}

func (g Game) ActivePiece() Tetromino {
	return g.activePiece
}

func (g *Game) Drop() {
	g.activePiece.Drop()
}

func (g Game) Board() Grid {
	return g.board
}

func (g *Game) Start() {
	ResetTGMRandomizer()
	g.level = 0
	g.board = NewGrid(gStartX, gStartY, gSize, gXLength, gYLength)
	g.activePiece, g.nextPiece = NextTGMRandomizer(), NextTGMRandomizer()
	g.SpawnTetromino(&g.activePiece)
}

func (g *Game) NextTetromino() {
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

	valid := true

	for _, v := range testPiece.Blocks() {
		if v.X < g.board.X() || v.X >= g.board.X()+g.board.PixelWidth() {
			valid = false
		}
	}

	if valid {
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
		valid := true

		if clockwise {
			t.RotateClockwise()
		} else {
			t.RotateCounterClockwise()
		}

		for _, v := range t.Blocks() {
			if v.X < g.board.X() || v.X >= g.board.X()+g.board.PixelWidth() {
				valid = false
			}
		}

		return valid
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
