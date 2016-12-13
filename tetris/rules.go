package tetris

import (
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const tetronimos int32 = 7

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
	return GenerateTetronimo(tS)
}

// GetTGMGravityMap get tgm grav rules
func GetTGMGravityMap() map[int]float64 {
	return tgmGravity
}

// SpawnTetromino on the board
func SpawnTetromino(b []sdl.Rect, t *Tetromino) {
	t.move(b[3].X, b[3].Y-b[3].H)
}
