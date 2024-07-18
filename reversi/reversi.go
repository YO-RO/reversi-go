package reversi

import (
	b "github.com/YO-RO/reversi-go/reversi/board"
)

type Reversi struct {
	CurrStone     b.Stone
	Skipped       bool
	skippedInARow bool
	end           bool
	board         b.Board
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

func (r *Reversi) testPut(row, col int, stopStone b.Stone) [][2]int {
	if s, ok := r.board.Get(row, col); !ok || s != b.None {
		return [][2]int{}
	}
	if stopStone != r.CurrStone && stopStone != b.None {
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
	found := make([][2]int, 0)
	for _, dir := range dirs {
		candedates := make([][2]int, 0)

		origin := [2]int{row, col}
		// originのマスは除外するため[1:]
		checkingMasses := r.board.LinearExtract(origin, dir)[1:]
		for _, m := range checkingMasses {
			if m.Stone == r.CurrStone.Reversed() {
				candedates = append(
					candedates,
					[2]int{m.Row, m.Col},
				)
			} else if m.Stone == stopStone {
				found = append(found, candedates...)
				break
			} else {
				break
			}
		}
	}
	return found
}

func (r *Reversi) currentMustSkip() bool {
	for _, mass := range r.board.GetAll() {
		row, col := mass.Row, mass.Col
		if mass.Stone != b.None {
			continue
		}
		reverseStones := r.testPut(row, col, b.None)
		if reverseStones == nil || len(reverseStones) == 0 {
			continue
		}
		return false
	}
	return true
}

func (r *Reversi) Put(row, col int) bool {
	r.Skipped = false

	reverseStonesIdx := r.testPut(row, col, r.CurrStone)
	if reverseStonesIdx == nil ||
		len(reverseStonesIdx) == 0 {
		return false
	}

	r.board.Put(row, col, r.CurrStone)
	for _, i := range reverseStonesIdx {
		r.board.Put(i[0], i[1], r.CurrStone)
	}

	r.CurrStone = r.CurrStone.Reversed()

	if r.currentMustSkip() {
		r.Skipped = true
		r.CurrStone = r.CurrStone.Reversed()
		if r.currentMustSkip() {
			r.skippedInARow = true
		}
		return true
	}
	return true
}

func (r *Reversi) Skip() {
	r.CurrStone = r.CurrStone.Reversed()
}

func (r *Reversi) Result() (b.Stone, bool) {
	noneCnt, blackCnt, whiteCnt := r.board.CountStone()
	if noneCnt > 0 && !r.skippedInARow {
		return b.None, false
	}
	var result b.Stone
	r.end = true
	switch {
	case blackCnt > whiteCnt:
		result = b.Black
	case blackCnt < whiteCnt:
		result = b.White
	default:
		result = b.None
	}
	return result, r.end
}

func (r *Reversi) String() string {
	board := r.board.String()
	return board
}
