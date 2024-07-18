package reversi

import (
	"strconv"
	"strings"
)

type Stone int

const (
	None Stone = iota
	Black
	White
)

type Reversi struct {
	CurrStone     Stone
	Skipped       bool
	skippedInARow bool
	end           bool
	result        Stone // Noneは引き分け
	Board         [8][8]Stone
	builder       strings.Builder
}

func NewReversi() Reversi {
	return Reversi{
		CurrStone: Black,
		Board: [8][8]Stone{
			3: {3: White, 4: Black},
			4: {3: Black, 4: White},
		},
	}
}

func (r *Reversi) testPut(row, col int, stopStone Stone) [][2]int {
	if row < 0 || row > 7 || col < 0 || col > 7 {
		return [][2]int{}
	}
	if r.Board[row][col] != None {
		return [][2]int{}
	}
	if stopStone != r.CurrStone && stopStone != None {
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
	candedates := make([][2]int, 0)
	for _, dir := range dirs {
		candedates = candedates[:0]
		rr, cc := row, col
		for {
			rr, cc = rr+dir[0], cc+dir[1]
			if rr < 0 || rr > 7 || cc < 0 || cc > 7 {
				break
			}
			if r.Board[rr][cc] == r.otherStone() {
				candedates = append(
					candedates,
					[2]int{rr, cc},
				)
			} else if r.Board[rr][cc] == stopStone {
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
	for row, line := range r.Board {
		for col, stone := range line {
			if stone != None {
				continue
			}
			reverseStonesPos := r.testPut(row, col, None)
			if reverseStonesPos == nil || len(reverseStonesPos) == 0 {
				continue
			}
			return false
		}
	}
	return true
}

func (r *Reversi) Put(row, col int) bool {
	r.Skipped = false

	reverseStonesPos := r.testPut(row, col, r.CurrStone)
	if reverseStonesPos == nil ||
		len(reverseStonesPos) == 0 {
		return false
	}

	r.Board[row][col] = r.CurrStone
	for _, pos := range reverseStonesPos {
		r.Board[pos[0]][pos[1]] = r.CurrStone
	}
	r.CurrStone = r.otherStone()

	if r.currentMustSkip() {
		r.Skipped = true
		r.CurrStone = r.otherStone()
		if r.currentMustSkip() {
			r.skippedInARow = true
		}
		return true
	}
	return true
}

func (r *Reversi) otherStone() Stone {
	switch r.CurrStone {
	case Black:
		return White
	case White:
		return Black
	default:
		return None
	}
}

func (r *Reversi) Skip() {
	r.CurrStone = r.otherStone()
}

func (r *Reversi) Result() (bool, Stone) {
	if r.end {
		return true, r.result
	}

	zeroMass := true
	blackCnt, whiteCnt := 0, 0
	for _, line := range r.Board {
		for _, mass := range line {
			switch mass {
			case None:
				zeroMass = false
			case Black:
				blackCnt++
			case White:
				whiteCnt++
			}
		}
	}

	if !zeroMass && !r.skippedInARow {
		return false, None
	}

	r.end = true
	switch {
	case blackCnt > whiteCnt:
		r.result = Black
	case blackCnt < whiteCnt:
		r.result = White
	default:
		r.result = None
	}
	return r.end, r.result
}

func (r *Reversi) String() string {
	r.builder.Reset()
	r.builder.WriteString("\n  a b c d e f g h\n")
	for i, line := range r.Board {
		r.builder.WriteString(strconv.Itoa(i + 1))
		for _, stone := range line {
			switch stone {
			case None:
				r.builder.WriteString(" ◌")
			case Black:
				r.builder.WriteString(" ●")
			case White:
				r.builder.WriteString(" ○")
			}
		}
		r.builder.WriteString("\n")
	}
	return r.builder.String()
}
