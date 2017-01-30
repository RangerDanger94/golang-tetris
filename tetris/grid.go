package tetris

import "github.com/veandco/go-sdl2/sdl"

type cell struct {
	color    sdl.Color
	rect     sdl.Rect
	occupied bool
}

// Grid - game board
type Grid struct {
	cellSize  int32
	x         int32
	y         int32
	spawnX    int32
	spawnY    int32
	altSpawnX int32
	altSpawnY int32
	width     int
	height    int
	cells     [][]cell
}

// NewGrid creates new tetris grid
func NewGrid(x, y, cellSize int32, width, height int) Grid {
	var g Grid
	g.x, g.y = x, y
	g.width, g.height = width, height
	g.cellSize = cellSize
	g.createGrid()
	g.spawnX, g.spawnY = g.cells[0][3].rect.X, g.cells[0][3].rect.Y-g.cellSize
	g.altSpawnX, g.altSpawnY = g.spawnX, g.spawnY-g.cellSize*2
	return g
}

// Create an empty row
func (g *Grid) createRow(row int) []cell {
	r := make([]cell, g.width)

	for column := range r {
		r[column].color = Black
		r[column].rect = sdl.Rect{X: g.x + int32(column)*g.cellSize, Y: int32(row) * g.cellSize, W: g.cellSize, H: g.cellSize}
		r[column].occupied = false
	}

	return r
}

func (g *Grid) createGrid() {
	g.cells = make([][]cell, g.height)

	for row := range g.cells {
		g.cells[row] = g.createRow(row)
	}
}

// Draw draws the grid with its locked pieces
func (g Grid) Draw(r *sdl.Renderer) {
	for _, row := range g.cells {
		for _, col := range row {
			//fmt.Printf("%v\n", col.color)
			r.SetDrawColor(col.color.R, col.color.G, col.color.B, col.color.A)
			r.FillRect(&col.rect)
		}
	}
}

// Unoccupied returns true if no elements are occupied
func (g Grid) Unoccupied() bool {
	for _, row := range g.cells {
		for _, col := range row {
			if col.occupied {
				return false
			}
		}
	}

	return true
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
