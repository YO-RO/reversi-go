package reversi

import (
	b "github.com/YO-RO/reversi-go/reversi/board"
)

type Reversi struct {
	CurrStone   b.Stone
	AutoSkipped bool
	end         bool
	board       b.Board
}

func NewReversi() Reversi {
	board := b.Board{}
	board.PutByLoc("4d", b.White)
	board.PutByLoc("4e", b.Black)
	board.PutByLoc("5d", b.Black)
	board.PutByLoc("5e", b.White)
	return Reversi{
		CurrStone: b.Black,
		board:     board,
	}
}

func (r *Reversi) testPut(row, col int, stone b.Stone) [][2]int {
	if s, ok := r.board.Get(row, col); !ok || s != b.None {
		return [][2]int{}
	}

	dirs := [...][2]int{
		{1, 0},
		{1, 1},
		{0, 1},
		{-1, 1},
		{-1, 0},
		{-1, -1},
		{0, -1},
		{1, -1},
	}
	// slice of [2]int{row, col}
	willReverseLocs := make([][2]int, 0)
	for _, dir := range dirs {
		candedates := make([][2]int, 0, 6)
		origin := [2]int{row, col}
	CHECK_REVERSE:
		// originのマスは除外するため[1:]
		for _, m := range r.board.LinearExtract(origin, dir)[1:] {
			switch m.Stone {
			case stone.Reversed():
				candedates = append(
					candedates,
					[2]int{m.Row, m.Col},
				)
			case stone:
				willReverseLocs = append(willReverseLocs, candedates...)
				break CHECK_REVERSE
			default:
				break CHECK_REVERSE
			}
		}
	}
	return willReverseLocs
}

func (r *Reversi) mustSkipCount(stone b.Stone) int {
	mustSkip := func(stone b.Stone) bool {
		for _, square := range r.board.GetAll() {
			row, col := square.Row, square.Col
			if square.Stone != b.None {
				continue
			}
			locToBeAbleToReverse := r.testPut(row, col, stone)
			if locToBeAbleToReverse == nil || len(locToBeAbleToReverse) == 0 {
				continue
			}
			return false
		}
		return true
	}

	switch {
	case mustSkip(stone) && mustSkip(stone.Reversed()):
		return 2
	case mustSkip(stone):
		return 1
	default:
		return 0
	}
}

func (r *Reversi) Put(row, col int) bool {
	r.AutoSkipped = false

	willReverseLocs := r.testPut(row, col, r.CurrStone)
	if willReverseLocs == nil || len(willReverseLocs) == 0 {
		return false
	}

	r.board.Put(row, col, r.CurrStone)
	for _, loc := range willReverseLocs {
		r.board.Reverse(loc[0], loc[1])
	}

	switch r.mustSkipCount(r.CurrStone.Reversed()) {
	case 0:
		r.CurrStone = r.CurrStone.Reversed()
	case 1:
		r.AutoSkipped = true
	case 2:
		r.end = true
	}
	return true
}

func (r *Reversi) Skip() {
	r.CurrStone = r.CurrStone.Reversed()
}

func (r *Reversi) Result() (b.Stone, bool) {
	if !r.end {
		return b.None, false
	}

	var winner b.Stone
	_, blackCnt, whiteCnt := r.board.CountStone()
	switch {
	case blackCnt > whiteCnt:
		winner = b.Black
	case blackCnt < whiteCnt:
		winner = b.White
	default:
		winner = b.None
	}
	return winner, r.end
}

func (r *Reversi) String() string {
	board := r.board.String()
	return board
}
