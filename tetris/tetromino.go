package tetris

import (
	"fmt"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

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

// Bounds - generate bounding box
func bounds(n int, d int32) []sdl.Rect {
	box := make([]sdl.Rect, n*n)
	var x, y int32 = 0, 0
	for i := 1; i <= n*n; i++ {
		box[i-1] = sdl.Rect{X: x, Y: y, W: d, H: d}
		x += d

		if i%n == 0 {
			y += d
			x = 0
		}
	}
	return box
}

// Tetromino - tetris block
type Tetromino struct {
	color        sdl.Color
	blocks       [4]sdl.Rect
	shape        int32
	orientations int
	orientation  int
	bounds       []sdl.Rect
}

func (t *Tetromino) MoveRight() {
	for i := range t.bounds {
		t.bounds[i].X += t.bounds[i].W
	}

	t.blocks = t.setOrientation(t.orientation)
}

func (t *Tetromino) MoveLeft() {
	for i := range t.bounds {
		t.bounds[i].X -= t.bounds[i].W
	}

	t.blocks = t.setOrientation(t.orientation)
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
	t.orientation = 1
	t.orientations = 2
	t.color = sdl.Color{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}
	t.bounds = bounds(4, 20)
	t.blocks = t.setOrientation(t.orientation)
	return t
}

// JTetromino
// [][][]
//     []
func JTetromino() Tetromino {
	var t Tetromino
	t.shape = J
	t.orientation = 1
	t.orientations = 4
	t.color = sdl.Color{R: 0x00, G: 0x00, B: 0xFF, A: 0xFF}
	t.bounds = bounds(3, 20)
	t.blocks = t.setOrientation(t.orientation)
	return t
}

// LTetromino
// [][][]
// []
func LTetromino() Tetromino {
	var t Tetromino
	t.shape = L
	t.orientation = 1
	t.orientations = 4
	t.color = sdl.Color{R: 0xEF, G: 0x79, B: 0x21, A: 0xFF}
	t.bounds = bounds(3, 20)
	t.blocks = t.setOrientation(t.orientation)
	return t
}

// OTetromino
// [][]
// [][]
func OTetromino() Tetromino {
	var t Tetromino
	t.shape = O
	t.orientation = 1
	t.orientations = 1
	t.color = sdl.Color{R: 0xF7, G: 0xD3, B: 0x08, A: 0xFF}
	t.bounds = bounds(4, 20)
	t.blocks = t.setOrientation(t.orientation)
	return t
}

// TTetromino
// [][][]
//   []
func TTetromino() Tetromino {
	var t Tetromino
	t.shape = T
	t.orientation = 1
	t.orientations = 4
	t.color = sdl.Color{R: 0x31, G: 0xC7, B: 0xEF, A: 0xFF}
	t.bounds = bounds(3, 20)
	t.blocks = t.setOrientation(t.orientation)
	return t
}

// STetromino
//   [][]
// [][]
func STetromino() Tetromino {
	var t Tetromino
	t.shape = S
	t.orientation = 1
	t.orientations = 2
	t.color = sdl.Color{R: 0xAD, G: 0x4D, B: 0x9C, A: 0xFF}
	t.bounds = bounds(3, 20)
	t.blocks = t.setOrientation(t.orientation)
	return t
}

// 	ZTetromino
//  [][]
// 	  [][]
func ZTetromino() Tetromino {
	var t Tetromino
	t.shape = Z
	t.orientation = 1
	t.orientations = 2
	t.color = sdl.Color{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF}
	t.bounds = bounds(3, 20)
	t.blocks = t.setOrientation(t.orientation)
	return t
}

// GenerateTetronimo creates a tetronimo of the given type
func GenerateTetronimo(s int32) Tetromino {
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
func (t Tetromino) setOrientation(o int) [4]sdl.Rect {
	var r [4]sdl.Rect

	switch t.shape {
	case I:
		switch o {
		case 1:
			r = [4]sdl.Rect{t.bounds[4], t.bounds[5], t.bounds[6], t.bounds[7]}
		case 2:
			r = [4]sdl.Rect{t.bounds[2], t.bounds[6], t.bounds[10], t.bounds[14]}
		}
	case J:
		switch o {
		case 1:
			r = [4]sdl.Rect{t.bounds[3], t.bounds[4], t.bounds[5], t.bounds[8]}
		case 2:
			r = [4]sdl.Rect{t.bounds[1], t.bounds[4], t.bounds[7], t.bounds[6]}
		case 3:
			r = [4]sdl.Rect{t.bounds[3], t.bounds[6], t.bounds[7], t.bounds[8]}
		case 4:
			r = [4]sdl.Rect{t.bounds[1], t.bounds[4], t.bounds[7], t.bounds[2]}
		}
	case L:
		switch o {
		case 1:
			r = [4]sdl.Rect{t.bounds[3], t.bounds[4], t.bounds[5], t.bounds[6]}
		case 2:
			r = [4]sdl.Rect{t.bounds[1], t.bounds[4], t.bounds[7], t.bounds[0]}
		case 3:
			r = [4]sdl.Rect{t.bounds[5], t.bounds[6], t.bounds[7], t.bounds[8]}
		case 4:
			r = [4]sdl.Rect{t.bounds[1], t.bounds[4], t.bounds[7], t.bounds[8]}
		}
	case O:
		r = [4]sdl.Rect{t.bounds[5], t.bounds[6], t.bounds[9], t.bounds[10]}
	case T:
		switch o {
		case 1:
			r = [4]sdl.Rect{t.bounds[3], t.bounds[4], t.bounds[5], t.bounds[7]}
		case 2:
			r = [4]sdl.Rect{t.bounds[1], t.bounds[4], t.bounds[7], t.bounds[3]}
		case 3:
			r = [4]sdl.Rect{t.bounds[4], t.bounds[6], t.bounds[7], t.bounds[8]}
		case 4:
			r = [4]sdl.Rect{t.bounds[1], t.bounds[4], t.bounds[7], t.bounds[5]}
		}
	case S:
		switch o {
		case 1:
			r = [4]sdl.Rect{t.bounds[6], t.bounds[7], t.bounds[4], t.bounds[5]}
		case 2:
			r = [4]sdl.Rect{t.bounds[0], t.bounds[3], t.bounds[4], t.bounds[7]}
		}

	case Z:
		switch o {
		case 1:
			r = [4]sdl.Rect{t.bounds[3], t.bounds[4], t.bounds[7], t.bounds[8]}
		case 2:
			r = [4]sdl.Rect{t.bounds[4], t.bounds[7], t.bounds[5], t.bounds[2]}
		}

	}

	return r
}

// Blocks - Returns a slice of the rects that make up the tetromino
func (t Tetromino) Blocks() []sdl.Rect {
	return t.blocks[:]
}

// Color - Returns tetromino color
func (t Tetromino) Color() sdl.Color {
	return t.color
}

// Rotate - rotates the tetromino according to TGM rotation rules
func (t *Tetromino) rotate(d float64) {
	for i, pos := range t.blocks {
		var oX, oY int32 = (pos.W * 3) / 2, (pos.W * 3) / 2
		fmt.Printf("Origin X&Y:\t%v\n", oX)
		dX, dY := oX-pos.X, oY-pos.Y

		t.blocks[i].X = int32(cosDegrees(d)*float64(dX)) + int32(-sinDegrees(d)*float64(dY)) + oX
		t.blocks[i].Y = int32(sinDegrees(d)*float64(dX)) + int32(cosDegrees(d)*float64(dY)) + oY
	}
}

// Resize - scale the termino
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
}

// RotateClockwise - calls rotate with input of 90
func (t *Tetromino) RotateClockwise() {
	t.orientation++
	if t.orientation > t.orientations {
		t.orientation = 1
	}
	fmt.Printf("Orientation is %v\n", t.orientation)
	t.blocks = t.setOrientation(t.orientation)
}

// RotateCounterClockwise - calls rotate with input of 270
func (t *Tetromino) RotateCounterClockwise() {
	t.orientation--
	if t.orientation < 1 {
		t.orientation = t.orientations
	}
	fmt.Printf("Orientation is %v\n", t.orientation)
	t.blocks = t.setOrientation(t.orientation)
}

// Drop - tetromino falls one grid space
func (t *Tetromino) Drop() {
	for i := range t.bounds {
		t.bounds[i].Y += t.bounds[i].H
	}

	t.blocks = t.setOrientation(t.orientation)
}

// Draw - draw to passed in renderer
func (t *Tetromino) Draw(r *sdl.Renderer) {
	r.SetDrawColor(0x0, 0x0, 0x0, t.color.A)
	r.FillRects(t.bounds)

	r.SetDrawColor(t.color.R, t.color.G, t.color.B, t.color.A)
	r.FillRects(t.Blocks())
}

func cosDegrees(d float64) float64 {
	return math.Cos(d * math.Pi / 180.0)
}

func sinDegrees(d float64) float64 {
	return math.Sin(d * math.Pi / 180.0)
}
