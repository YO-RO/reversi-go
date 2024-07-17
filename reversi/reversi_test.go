package reversi_test

import (
	"testing"

	"github.com/YO-RO/reversi-go/reversi"
)

func TestNewReversi(t *testing.T) {
	r := reversi.NewReversi()
	if r.CurrStone != reversi.Black {
		t.Errorf("reversi.Reversi.CurrentStone == %v, want %v", r.CurrStone, reversi.Black)
	}
	// 4d white, 4f black
	// 5d black, 5f white
	if r.Board[3][3] != reversi.White || r.Board[3][4] != reversi.Black ||
		r.Board[4][3] != reversi.Black || r.Board[4][4] != reversi.White {
		b := [8][8]reversi.Stone{
			3: {3: reversi.White, 4: reversi.Black},
			4: {3: reversi.Black, 4: reversi.White},
		}
		t.Errorf("reversi.Reversi.Board ==\n%v, want \n%v", r.Board, b)
	}

}
