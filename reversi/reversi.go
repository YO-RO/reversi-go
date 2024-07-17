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
	CurrStone Stone
	Winner    Stone
	Board     [8][8]Stone
	b         strings.Builder
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

func (r *Reversi) String() string {
	r.b.Reset()
	r.b.WriteString("\n  a b c d e f g h\n")
	for i, line := range r.Board {
		r.b.WriteString(strconv.Itoa(i + 1))
		for _, stone := range line {
			switch stone {
			case None:
				r.b.WriteString(" ◌")
			case Black:
				r.b.WriteString(" ●")
			case White:
				r.b.WriteString(" ○")
			}
		}
		r.b.WriteString("\n")
	}
	return r.b.String()
}
