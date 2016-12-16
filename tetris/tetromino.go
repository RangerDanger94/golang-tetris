package tetris

import (
	"fmt"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

// Tetromino color scheme
var (
	Red    = sdl.Color{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}
	Blue   = sdl.Color{R: 0x00, G: 0x00, B: 0xFF, A: 0xFF}
	Orange = sdl.Color{R: 0xEF, G: 0x79, B: 0x21, A: 0xFF}
	Yellow = sdl.Color{R: 0xF7, G: 0xD3, B: 0x08, A: 0xFF}
	Aqua   = sdl.Color{R: 0x31, G: 0xC7, B: 0xEF, A: 0xFF}
	Purple = sdl.Color{R: 0xAD, G: 0x4D, B: 0x9C, A: 0xFF}
	Green  = sdl.Color{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF}
	Black  = sdl.Color{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
)

const tetrominos int32 = 7

// Tetronimo shapes
const (
	I = iota
	J
	L
	O
	S
	T
	Z
)

func (t *Tetromino) setBounds(sX int32, sY int32) {
	t.bounds = make([]sdl.Rect, t.boundaryArea*t.boundaryArea)
	var x, y int32 = sX, sY
	for i := 1; i <= t.boundaryArea*t.boundaryArea; i++ {
		t.bounds[i-1] = sdl.Rect{X: x, Y: y, W: t.size, H: t.size}
		x += t.size

		if i%t.boundaryArea == 0 {
			y += t.size
			x = sX
		}
	}
}

// Tetromino - tetris block
type Tetromino struct {
	shape        int32
	size         int32
	color        sdl.Color
	orientation  int
	orientations int
	boundaryArea int
	bounds       []sdl.Rect
	blocks       [4]sdl.Rect
}

// Type e.g. I, J, L
func (t Tetromino) Type() int32 {
	return t.shape
}

// Shape - getter for shape
func (t Tetromino) Shape() int32 {
	return t.shape
}

// ITetromino
// [][][][]
func ITetromino() Tetromino {
	var t Tetromino
	t.shape = I
	t.size = 10
	t.orientation = 1
	t.orientations = 2
	t.color = Red
	t.boundaryArea = 4
	t.setBounds(0, 0)
	t.setOrientation(t.orientation)
	return t
}

// JTetromino
// [][][]
//     []
func JTetromino() Tetromino {
	var t Tetromino
	t.shape = J
	t.size = 10
	t.orientation = 1
	t.orientations = 4
	t.color = Blue
	t.boundaryArea = 3
	t.setBounds(0, 0)
	t.setOrientation(t.orientation)
	return t
}

// LTetromino
// [][][]
// []
func LTetromino() Tetromino {
	var t Tetromino
	t.shape = L
	t.size = 10
	t.orientation = 1
	t.orientations = 4
	t.color = Orange
	t.boundaryArea = 3
	t.setBounds(0, 0)
	t.setOrientation(t.orientation)
	return t
}

// OTetromino
// [][]
// [][]
func OTetromino() Tetromino {
	var t Tetromino
	t.shape = O
	t.size = 10
	t.orientation = 1
	t.orientations = 1
	t.color = Yellow
	t.boundaryArea = 4
	t.setBounds(0, 0)
	t.setOrientation(t.orientation)
	return t
}

// TTetromino
// [][][]
//   []
func TTetromino() Tetromino {
	var t Tetromino
	t.shape = T
	t.size = 10
	t.orientation = 1
	t.orientations = 4
	t.color = Aqua
	t.boundaryArea = 3
	t.setBounds(0, 0)
	t.setOrientation(t.orientation)
	return t
}

// STetromino
//   [][]
// [][]
func STetromino() Tetromino {
	var t Tetromino
	t.shape = S
	t.size = 10
	t.orientation = 1
	t.orientations = 2
	t.color = Purple
	t.boundaryArea = 3
	t.setBounds(0, 0)
	t.setOrientation(t.orientation)
	return t
}

// 	ZTetromino
//  [][]
// 	  [][]
func ZTetromino() Tetromino {
	var t Tetromino
	t.shape = Z
	t.size = 10
	t.orientation = 1
	t.orientations = 2
	t.color = Green
	t.boundaryArea = 3
	t.setBounds(0, 0)
	t.setOrientation(t.orientation)
	return t
}

// GenerateTetronimo creates a tetronimo of the given type
func generateTetronimo(s int32) Tetromino {
	switch s {
	case I:
		return ITetromino()
	case J:
		return JTetromino()
	case L:
		return LTetromino()
	case O:
		return OTetromino()
	case T:
		return TTetromino()
	case S:
		return STetromino()
	case Z:
		return ZTetromino()
	default:
		return ITetromino()
	}
}

// [0][1][2]	[0 ][1 ][2 ][3 ]
// [3][4][5]	[4 ][5 ][6 ][7 ]
// [6][7][8]	[8 ][9 ][10][11]
//				[12][13][14][15]
// Apply rotations according to TGM rotation rules
func (t *Tetromino) setOrientation(o int) {
	switch t.shape {
	case I:
		switch o {
		case 1:
			t.blocks = [4]sdl.Rect{t.bounds[4], t.bounds[5], t.bounds[6], t.bounds[7]}
		case 2:
			t.blocks = [4]sdl.Rect{t.bounds[2], t.bounds[6], t.bounds[10], t.bounds[14]}
		}
	case J:
		switch o {
		case 1:
			t.blocks = [4]sdl.Rect{t.bounds[3], t.bounds[4], t.bounds[5], t.bounds[8]}
		case 2:
			t.blocks = [4]sdl.Rect{t.bounds[1], t.bounds[4], t.bounds[7], t.bounds[6]}
		case 3:
			t.blocks = [4]sdl.Rect{t.bounds[3], t.bounds[6], t.bounds[7], t.bounds[8]}
		case 4:
			t.blocks = [4]sdl.Rect{t.bounds[1], t.bounds[4], t.bounds[7], t.bounds[2]}
		}
	case L:
		switch o {
		case 1:
			t.blocks = [4]sdl.Rect{t.bounds[3], t.bounds[4], t.bounds[5], t.bounds[6]}
		case 2:
			t.blocks = [4]sdl.Rect{t.bounds[1], t.bounds[4], t.bounds[7], t.bounds[0]}
		case 3:
			t.blocks = [4]sdl.Rect{t.bounds[5], t.bounds[6], t.bounds[7], t.bounds[8]}
		case 4:
			t.blocks = [4]sdl.Rect{t.bounds[1], t.bounds[4], t.bounds[7], t.bounds[8]}
		}
	case O:
		t.blocks = [4]sdl.Rect{t.bounds[5], t.bounds[6], t.bounds[9], t.bounds[10]}
	case T:
		switch o {
		case 1:
			t.blocks = [4]sdl.Rect{t.bounds[3], t.bounds[4], t.bounds[5], t.bounds[7]}
		case 2:
			t.blocks = [4]sdl.Rect{t.bounds[1], t.bounds[4], t.bounds[7], t.bounds[3]}
		case 3:
			t.blocks = [4]sdl.Rect{t.bounds[4], t.bounds[6], t.bounds[7], t.bounds[8]}
		case 4:
			t.blocks = [4]sdl.Rect{t.bounds[1], t.bounds[4], t.bounds[7], t.bounds[5]}
		}
	case S:
		switch o {
		case 1:
			t.blocks = [4]sdl.Rect{t.bounds[6], t.bounds[7], t.bounds[4], t.bounds[5]}
		case 2:
			t.blocks = [4]sdl.Rect{t.bounds[0], t.bounds[3], t.bounds[4], t.bounds[7]}
		}

	case Z:
		switch o {
		case 1:
			t.blocks = [4]sdl.Rect{t.bounds[3], t.bounds[4], t.bounds[7], t.bounds[8]}
		case 2:
			t.blocks = [4]sdl.Rect{t.bounds[4], t.bounds[7], t.bounds[5], t.bounds[2]}
		}

	}
}

// Blocks returns blocks that make up actual Tetromino
func (t Tetromino) Blocks() []sdl.Rect {
	return t.blocks[:]
}

// Color returns tetromino color
func (t Tetromino) Color() sdl.Color {
	return t.color
}

// Resize scales the tetromino and its bounding box
func (t *Tetromino) Resize(d int32) {
	for i := range t.bounds {
		if t.bounds[i].X != 0 {
			t.bounds[i].X = (d / t.bounds[i].W) * t.bounds[i].X
		}

		if t.bounds[i].Y != 0 {
			t.bounds[i].Y = (d / t.bounds[i].W) * t.bounds[i].Y
		}

		t.bounds[i].W = d
		t.bounds[i].H = d
	}

	t.size = d
}

func (t *Tetromino) move(x int32, y int32) {
	t.setBounds(x, y)
	t.setOrientation(t.orientation)
}

// ShiftRight shifts tetromino to the right 1 grid space
func (t *Tetromino) ShiftRight() {
	t.move(t.bounds[0].X+20, t.bounds[0].Y)
}

// ShiftLeft shifts tetromino to the left 1 grid space
func (t *Tetromino) ShiftLeft() {
	t.move(t.bounds[0].X-20, t.bounds[0].Y)
}

// Drop drops Tetromino 1 grid space
func (t *Tetromino) Drop() {
	t.move(t.bounds[0].X, t.bounds[0].Y+20)
}

// Draw uses passed in renderer to draw tetromino
func (t Tetromino) Draw(r *sdl.Renderer) {
	// r.SetDrawColor(0x0, 0x0, 0x0, t.color.A)
	// r.FillRects(t.bounds)

	r.SetDrawColor(t.color.R, t.color.G, t.color.B, t.color.A)
	r.FillRects(t.Blocks())
}

// RotateClockwise rotates tetromino clockwise
func (t *Tetromino) RotateClockwise() {
	t.orientation++
	if t.orientation > t.orientations {
		t.orientation = 1
	}

	t.setOrientation(t.orientation)
}

// RotateCounterClockwise rotates tetromino counter-clockwise
func (t *Tetromino) RotateCounterClockwise() {
	t.orientation--
	if t.orientation < 1 {
		t.orientation = t.orientations
	}

	t.setOrientation(t.orientation)
}

// 2D matrix rotations
func (t *Tetromino) rotate(d float64) {
	for i, pos := range t.blocks {
		var oX, oY int32 = (pos.W * 3) / 2, (pos.W * 3) / 2
		fmt.Printf("Origin X&Y:\t%v\n", oX)
		dX, dY := oX-pos.X, oY-pos.Y

		t.blocks[i].X = int32(cosDegrees(d)*float64(dX)) + int32(-sinDegrees(d)*float64(dY)) + oX
		t.blocks[i].Y = int32(sinDegrees(d)*float64(dX)) + int32(cosDegrees(d)*float64(dY)) + oY
	}
}

func cosDegrees(d float64) float64 {
	return math.Cos(d * math.Pi / 180.0)
}

func sinDegrees(d float64) float64 {
	return math.Sin(d * math.Pi / 180.0)
}
