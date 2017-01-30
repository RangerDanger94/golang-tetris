package tetris

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"gitlab.com/rangerdanger/sdlaudio"
)

// Board presets
const gXLength int = 10
const gYLength int = 20
const gSize int32 = 20
const gStartX int32 = 600/2 - (int32(gXLength)*gSize)/2
const gStartY int32 = 0

// Game steps
const (
	Menu = iota
	Locking
	Clearing
	ClearDelay
	Spawning
	Transition
	GameOver
)

// Input commands
const (
	_ = iota
	ShiftLeft
	ShiftRight
	RotateClockwise
	RotateCounterClockwise
	ManualDrop
	Start
)

type Game struct {
	start       time.Time
	activePiece Tetromino
	nextPiece   Tetromino
	holdPiece   Tetromino
	board       Grid
	command     int32
	lastCommand int32

	areDelay, areFrames     int
	dasDelay, dasFrames     int
	lockDelay, lockFrames   int
	clearDelay, clearFrames int
	activeFrames            int
	step                    int

	soft       bool
	softFrames int

	gravFrames float64
	gravity    float64

	level int
	score int
	combo int
	bravo int
}

// NewGame returns a new game struct
func NewGame() *Game {
	return new(Game)
}

// Step returns game step
func (g Game) Step() int {
	return g.step
}

// Init sets up the games variables
func (g *Game) Init() {
	ResetTGMRandomizer()
	g.loadMusic()
	sdlaudio.PlayMusic("menu", -1)

	g.level = 0
	g.score = 0
	g.combo = 1
	g.bravo = 1
	g.command = 0

	g.areDelay, g.areFrames = 30, 0
	g.dasDelay, g.dasFrames = 14, 0
	g.lockDelay, g.lockFrames = 30, 0
	g.clearDelay, g.clearFrames = 41, 0
	g.activeFrames = 0
	g.softFrames = 0
	g.soft = false
	g.step = Menu

}

// Start initalizes game
func (g *Game) Start() {
	g.board = NewGrid(gStartX, gStartY, gSize, gXLength, gYLength)
	g.step = Locking
	sdlaudio.PlayMusic("easy", -1)
	g.activePiece, g.nextPiece = NextTGMRandomizer(), NextTGMRandomizer()
	g.SpawnTetromino(&g.activePiece)
}

func (g *Game) RunTime() string {
	since := time.Since(g.start)
	return fmt.Sprintf("%02f:%02f:%02f", since.Minutes(), since.Seconds(), since.Seconds()/1000.0)
}

// BufferCommand sets active command
func (g *Game) BufferCommand(command int32) {
	g.command = command
}

// Increment the level counter
func (g *Game) nextLevelRequiresClear() bool {
	if g.level+1%100 == 0 || g.level == 998 {
		return true
	}

	return false
}

// LoadMusic passes back the map of audio assets to be loaded
func (g *Game) loadMusic() error {
	for k, v := range tgmAudio {
		if err := sdlaudio.LoadMusic(k, v); err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

var lastStep int

// ProcessFrame runs the game logic for a frame
func (g *Game) ProcessFrame() {
	switch g.step {
	case Menu:
		if g.command == Start {
			g.step = Transition
			lastStep = Menu
		}
	case Transition:
		var track string
		if lastStep == Menu {
			track = "start"
		} else if lastStep == GameOver {
			track = "gameOver"
		}

		if b, err := sdlaudio.PlayMusicOS(track); b {
			if lastStep == Menu {
				g.Start()
				g.step = Locking
			} else if lastStep == GameOver {
				sdlaudio.PlayMusic("menu", -1)
				g.step = Menu
			}

		} else if err != nil {
			panic(err)
		}
	case GameOver:
		g.step = Transition
		lastStep = GameOver
	}

	if g.step != Menu && g.step != Transition {
		g.doGravity()
		g.soft = false

		if g.command == ShiftLeft || g.command == ShiftRight {
			if !g.dasLocked() {
				g.dasFrames++
			}

			if g.dasFrames == 1 || g.dasFrames >= g.dasDelay {
				if g.command == ShiftLeft {
					g.tryShift(false)
				} else {
					g.tryShift(true)
				}
			}

		} else {
			if !g.dasLocked() {
				g.dasFrames = 0
			}

			if g.command == RotateClockwise && g.lastCommand != RotateClockwise {
				g.tryRotate(true)
			} else if g.command == RotateCounterClockwise && g.lastCommand != RotateCounterClockwise {
				g.tryRotate(false)
			} else if g.command == ManualDrop {
				g.softFrames++
				g.soft = true
				g.tryDrop()
			}
		}

		// State machine
		switch g.step {
		case Locking:
			g.activeFrames++
			if g.checkLock() {

				// Check we aren't out of bounds
				if g.activePiece.Above(g.board.Y()) {
					g.step = GameOver
				} else {
					g.activeFrames = 0
					g.activePiece, g.nextPiece = g.nextPiece, NextTGMRandomizer()
					g.step = Clearing
				}
			}
		case Clearing:
			if g.checkClear() {
				g.step = ClearDelay
			} else {
				g.step = Spawning
			}
		case ClearDelay:
			g.clearFrames++
			if g.clearFrames >= g.clearDelay {
				g.clearFrames = 0
				g.step = Spawning
			}
		case Spawning:
			g.areFrames++
			if g.areFrames >= g.areDelay {
				g.areFrames = 0
				g.softFrames = 0
				g.SpawnTetromino(&g.activePiece)
				g.step = Locking
			}
		}
	}

	g.lastCommand = g.command
}

// The players DAS charge is unmodified during line clear delay,
// the first 4 frames of ARE, the last frame of ARE and the frame
// on which a piece spawns
func (g Game) dasLocked() bool {
	if g.activeFrames != 1 &&
		g.clearFrames == 0 &&
		g.areFrames > 4 &&
		g.areFrames != g.areDelay {
		return true
	}
	return false
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
		g.tryDrop()
		g.gravFrames--
	}
}

// Attempts to drop piece if valid
func (g *Game) tryDrop() {
	testPiece := g.activePiece
	testPiece.Drop()

	if !g.collision(testPiece) {
		g.activePiece.Drop()
	}
}

func (g *Game) checkLock() bool {
	testPiece := g.activePiece
	testPiece.Drop()

	if g.collision(testPiece) {
		g.lockFrames++

		if g.soft { // manual locking
			g.lockFrames = g.lockDelay
		}
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

		return true
	}

	return false
}

// check all rows for successful line clear
func (g *Game) checkClear() bool {
	cleared := 0
	for row := range g.board.cells {
		if g.checkLineClear(row) {
			g.clearLine(row)
			cleared++
		}
	}

	// Check Bravo & Combo, update Score
	if cleared > 0 {
		if g.board.Unoccupied() {
			g.bravo = 4
		} else {
			g.bravo = 1
		}
		g.combo += (2 * cleared) - 2
		g.score += (roof(g.level+cleared, 4) + g.softFrames) * cleared * ((2 * cleared) - 1) * g.combo * g.bravo
	} else {
		g.combo = 1
	}

	g.level += cleared // level up

	return cleared > 0
}

// Return the rounded up value of a division
func roof(foo, bar int) int {
	if foo%bar == 0 {
		return int(foo / bar)
	}

	return int(foo/bar) + 1
}

// Return true if all cells in the row are occupied
func (g *Game) checkLineClear(row int) bool {
	for col := range g.board.cells[row] {
		if !g.board.cells[row][col].occupied {
			return false
		}
	}

	return true
}

// Drop lines above cleared line
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

// CollisionRects returns the sdl.Rect elements from lockedPieces
func (g Game) terminoIntersection(t Tetromino, r sdl.Rect) bool {
	for _, v := range t.blocks {
		if v.HasIntersection(&r) {
			return true
		}
	}

	return false
}

// Collision checks if a tetromino is colliding with the following
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

// SpawnTetromino on the grid
func (g *Game) SpawnTetromino(t *Tetromino) {
	t.Resize(g.board.cellSize)
	t.move(g.board.spawnX, g.board.spawnY)

	if g.collision(*t) {
		t.move(g.board.altSpawnX, g.board.altSpawnY)
	}

	g.activeFrames = 1
	if !g.nextLevelRequiresClear() {
		g.level++
	}
}

// TryShift will check and perform valid shift
func (g *Game) tryShift(right bool) {
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

// TryRotate will check and perform valid rotations
func (g *Game) tryRotate(clockwise bool) {
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
	tS := rand.Int31n(tetrominos)

	// The game never deals an S, Z or O as the first piece
	if tgmFirstPiece {
		for tS == S || tS == Z || tS == O {
			tS = rand.Int31n(tetrominos)
		}
	}

	// Attempt to get a tetronimo not in the bag history
	for _, t := range tgmBagHistory {
		if tS == t {
			tS = rand.Int31n(tetrominos)
		}
	}

	for i := len(tgmBagHistory) - 1; i > 0; i-- {
		tgmBagHistory[i] = tgmBagHistory[i-1]
	}
	tgmBagHistory = append([]int32{tS}, tgmBagHistory[1:4]...)
	tgmFirstPiece = false

	return generateTetronimo(tS)
}

var tgmAudio = map[string]string{
	"start":    "assets/03_insert_coin.mp3",
	"easy":     "assets/04_hardening_drops.mp3",
	"hard":     "assets/05_hardening_drops_hard.mp3",
	"menu":     "assets/07_menu.mp3",
	"gameOver": "assets/08_game_over.mp3",
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

var tgmGrading = map[string]int{
	"9":  0,
	"8":  400,
	"7":  800,
	"6":  1400,
	"5":  2000,
	"4":  3500,
	"3":  5500,
	"2":  8000,
	"1":  12000,
	"S1": 16000,
	"S2": 22000,
	"S3": 30000,
	"S4": 40000,
	"S5": 52000,
	"S6": 66000,
	"S7": 82000,
	"S8": 100000,
	"S9": 120000,
}

// GetTGMGravityMap get tgm grav rules
func GetTGMGravityMap() map[int]float64 {
	return tgmGravity
}
