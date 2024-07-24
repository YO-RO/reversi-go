package reversi

import (
	b "github.com/YO-RO/reversi-go/reversi/board"
)

type Reversi struct {
	currStone b.Stone
	skipped   bool
	ended     bool
	board     b.Board
}

func NewReversi() Reversi {
	board := b.Board{}
	board.PutBySign("4d", b.White)
	board.PutBySign("4e", b.Black)
	board.PutBySign("5d", b.Black)
	board.PutBySign("5e", b.White)
	return Reversi{
		currStone: b.Black,
		board:     board,
	}
}

type State struct {
	CurrStone   b.Stone
	SkippedPrev bool

	Board      [8][8]b.Stone
	Candidates [][2]int

	BlackCnt, WhiteCnt int
}

func (r *Reversi) State() State {
	_, b, w := r.board.CountStone()
	return State{
		CurrStone:   r.currStone,
		SkippedPrev: r.skipped,

		Board:      r.board.Grid(),
		Candidates: r.putCandidates(r.currStone),

		BlackCnt: b,
		WhiteCnt: w,
	}
}

func (r *Reversi) testPut(row, col int, stone b.Stone) ([][2]int, bool) {
	if s, ok := r.board.Get(row, col); !ok || s != b.None {
		return nil, false
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
	stonesToReverse := make([][2]int, 0)
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
				stonesToReverse = append(stonesToReverse, candedates...)
				break CHECK_REVERSE
			default:
				break CHECK_REVERSE
			}
		}
	}

	if len(stonesToReverse) > 0 {
		return stonesToReverse, true
	} else {
		return nil, false
	}
}

func (r *Reversi) skipCount(stone b.Stone) int {
	requireSkip := func(stone b.Stone) bool {
		return len(r.putCandidates(stone)) == 0
	}

	switch {
	case requireSkip(stone) && requireSkip(stone.Reversed()):
		return 2
	case requireSkip(stone):
		return 1
	default:
		return 0
	}
}

func (r *Reversi) putCandidates(stone b.Stone) [][2]int {
	candidates := [][2]int{}
	for _, square := range r.board.GetAll() {
		if square.Stone != b.None {
			continue
		}
		_, canPut := r.testPut(square.Row, square.Col, stone)
		if canPut {
			candidates = append(candidates, [2]int{square.Row, square.Col})
		}
	}
	return candidates
}

func (r *Reversi) Put(row, col int) bool {
	r.skipped = false

	stoneIdxsToReverse, canPut := r.testPut(row, col, r.currStone)
	if !canPut {
		return false
	}

	r.board.Put(row, col, r.currStone)
	for _, stoneIdx := range stoneIdxsToReverse {
		r.board.Reverse(stoneIdx[0], stoneIdx[1])
	}

	switch r.skipCount(r.currStone.Reversed()) {
	case 0:
		r.currStone = r.currStone.Reversed()
	case 1:
		r.skipped = true
	case 2:
		r.ended = true
	}
	return true
}

func (r *Reversi) Result() (b.Stone, bool) {
	if !r.ended {
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
	return winner, r.ended
}
