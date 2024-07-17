package reversi

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
