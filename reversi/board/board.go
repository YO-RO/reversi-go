package board

import (
	"regexp"
	"strconv"
	"strings"
)

type Stone int

const (
	None Stone = iota
	Black
	White
)

func (s Stone) Reversed() Stone {
	switch s {
	case Black:
		return White
	case White:
		return Black
	default:
		return None
	}
}

// 2a -> 1, 0
// 6e -> 5, 4
func SignToIndex(sign string) (int, int, bool) {
	re := regexp.MustCompile(`^\s?([1-8])\s?([a-h])\s?$`)
	match := re.FindStringSubmatch(sign)
	if match == nil {
		return 0, 0, false
	}
	row, _ := strconv.Atoi(match[1])
	row -= 1 // インデックスは0はじまり
	col := int(match[2][0] - byte('a'))
	return row, col, true
}

// マスのこと
type Square struct {
	Row, Col int
	Stone    Stone
}

type Board struct {
	grid    [8][8]Stone
	builder strings.Builder
}

func validIndex(row, col int) bool {
	return row >= 0 && row <= 7 && col >= 0 && col <= 7
}

func (b *Board) Grid() [8][8]Stone {
	return b.grid
}

func (b *Board) Get(row, col int) (Stone, bool) {
	if !validIndex(row, col) {
		return None, false
	}
	return b.grid[row][col], true
}

func (b *Board) GetAll() [64]Square {
	var res [64]Square
	for row, line := range b.grid {
		for col, stone := range line {
			i := 8*row + col
			res[i] = Square{
				Row:   row,
				Col:   col,
				Stone: stone,
			}
		}
	}
	return res
}

func (b *Board) LinearExtract(origin, dir [2]int) []Square {
	res := []Square{}
	for i := origin; validIndex(i[0], i[1]); i[0], i[1] = i[0]+dir[0], i[1]+dir[1] {
		stone, _ := b.Get(i[0], i[1])
		res = append(res, Square{
			Row:   i[0],
			Col:   i[1],
			Stone: stone,
		})
	}
	return res
}

func (b *Board) Reverse(row, col int) bool {
	if !validIndex(row, col) || b.grid[row][col] == None {
		return false
	}
	b.grid[row][col] = b.grid[row][col].Reversed()
	return true
}

func (b *Board) Put(row, col int, stone Stone) bool {
	if !validIndex(row, col) || b.grid[row][col] != None {
		return false
	}
	b.grid[row][col] = stone
	return true
}

func (b *Board) PutBySign(sign string, stone Stone) bool {
	if row, col, ok := SignToIndex(sign); ok {
		return b.Put(row, col, stone)
	}
	return false
}

func (b *Board) CountStone() (none, black, white int) {
	for _, square := range b.GetAll() {
		switch square.Stone {
		case None:
			none++
		case Black:
			black++
		case White:
			white++
		default:
			panic("unexpected Stone: Stone must be 0 (None) or 1 (Black) or 2 (White)")
		}
	}
	return none, black, white
}

func (b *Board) String() string {
	b.builder.Reset()
	b.builder.WriteString("\n  a b c d e f g h\n")
	for i, line := range b.grid {
		b.builder.WriteString(strconv.Itoa(i + 1))
		for _, stone := range line {
			switch stone {
			case None:
				b.builder.WriteString(" _")
			case Black:
				b.builder.WriteString(" b")
			case White:
				b.builder.WriteString(" w")
			}
		}
		b.builder.WriteString("\n")
	}
	return b.builder.String()
}
