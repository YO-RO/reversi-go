package reversi

import (
	b "github.com/YO-RO/reversi-go/reversi/board"
)

type recorder struct {
	// 1aや2hなどのsignを格納する
	// スキップした場合は空文字
	// Blackは偶数のindex, Whiteは奇数のindex
	memory []string
}

type Record []RecordEntry

type RecordEntry struct {
	Black, White string
}

func (r *recorder) record(sign string, stone b.Stone) {
	if stone == b.None {
		panic("unexpected stone " + stone.String())
	}

	// 黒は偶数の、白は奇数のindexに入れたい
	if stone == b.Black && len(r.memory)%2 != 0 ||
		stone == b.White && len(r.memory)%2 == 0 {
		r.memory = append(r.memory, "")
	}
	r.memory = append(r.memory, sign)
}

func (r *recorder) export() Record {
	rec := make(Record, 0, len(r.memory)/2+1)
	for i := 0; i < len(r.memory); i += 2 {
		var black, white string
		black = r.memory[i]
		if i+1 < len(r.memory) {
			white = r.memory[i+1]
		}
		rec = append(rec, RecordEntry{Black: black, White: white})
	}
	return rec
}
