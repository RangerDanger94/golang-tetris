package tetris

import "github.com/veandco/go-sdl2/sdl"

// Grid - game board
type Grid struct {
	cellSize int32
	x        int32
	y        int32
	width    int
	height   int
	cells    []sdl.Rect
}

// NewGrid creates new tetris grid
func NewGrid(x, y, cellSize int32, width, height int) Grid {
	var g Grid
	g.x, g.y = x, y
	g.width, g.height = width, height
	g.cellSize = cellSize
	g.createGrid()
	return g
}

func (g *Grid) createGrid() {
	g.cells = make([]sdl.Rect, g.Area())
	var x, y int32 = g.x, g.y

	for i := 1; i <= g.Area(); i++ {
		g.cells[i-1] = sdl.Rect{X: x, Y: y, W: g.cellSize, H: g.cellSize}
		x += g.cellSize

		if i%int(g.width) == 0 {
			y += g.cellSize
			x = g.x
		}
	}
}

// Cells returns the grid
func (g Grid) Cells() []sdl.Rect {
	return g.cells
}

// Draw draws the grid with its locked pieces
func (g Grid) Draw(r *sdl.Renderer) {
	r.SetDrawColor(0x0, 0x0, 0x0, 0xFF)
	r.FillRects(g.cells)
}

// Area returns width * height
func (g Grid) Area() int {
	return g.width * g.height
}

// Width returns width in elements
func (g Grid) Width() int {
	return g.width
}

// Height returns height in elements
func (g Grid) Height() int {
	return g.height
}

// PixelWidth retunrs width in pixels
func (g Grid) PixelWidth() int32 {
	return int32(g.width) * g.cellSize
}

// PixelHeight returns width in pixels
func (g Grid) PixelHeight() int32 {
	return int32(g.height) * g.cellSize
}

// CellSize returns cellSize
func (g Grid) CellSize() int32 {
	return g.cellSize
}

func (g Grid) X() int32 {
	return g.x
}

func (g Grid) Y() int32 {
	return g.y
}

// Ground rects below the playing area
func (g Grid) Ground() []sdl.Rect {
	ground := make([]sdl.Rect, g.width)

	for i, j := g.Area()-g.width, 0; i < g.Area(); i, j = i+1, j+1 {
		ground[j] = sdl.Rect{X: g.cells[i].X, Y: g.cells[i].Y + g.cellSize, W: g.cellSize, H: g.cellSize}
	}

	return ground
}

// SpawnTetromino on the grid
func (g Grid) SpawnTetromino(t *Tetromino) {
	t.move(g.cells[3].X, g.cells[3].Y-g.cellSize)
}
