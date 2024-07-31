package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/YO-RO/reversi-go/reversi"
	b "github.com/YO-RO/reversi-go/reversi/board"
)

func stoneName(s b.Stone) string {
	switch s {
	case b.None:
		return "空"
	case b.Black:
		return "黒"
	case b.White:
		return "白"
	default:
		return ""
	}
}

func boardString(board [8][8]b.Stone, candidates [][2]int) string {
	builder := strings.Builder{}
	builder.WriteString("\n  a b c d e f g h\n")
	for row, line := range board {
		builder.WriteString(strconv.Itoa(row + 1))
		for col, stone := range line {
			switch stone {
			case b.None:
				if slices.Contains(candidates, [2]int{row, col}) {
					builder.WriteString(" .")
				} else {
					builder.WriteString(" _")
				}
			case b.Black:
				builder.WriteString(" b")
			case b.White:
				builder.WriteString(" w")
			}
		}
		builder.WriteString("\n")
	}
	return builder.String()
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	game := reversi.NewReversi()
	for {
		state := game.State()
		fmt.Println(boardString(state.Board, state.Candidates))

		if winner, black, white, ended := game.Result(); ended {
			fmt.Printf("%sの勝ち 黒 %d 白 %d\n", stoneName(winner), black, white)
			return
		}

		if state.SkippedPrev {
			prev := state.CurrStone.Reversed()
			fmt.Println(stoneName(prev), "はスキップしました")
		}
		fmt.Printf("黒: %d, 白: %d\n", state.BlackCnt, state.WhiteCnt)
		fmt.Printf("%v\n", game.ExportRecord())

		// ユーザーからの入力
		for {
			fmt.Print(stoneName(state.CurrStone) + "> ")
			scanner.Scan()
			input := scanner.Text()

			row, col, ok := b.SignToIndex(input)
			if !ok {
				fmt.Println("不正な入力")
				continue
			}
			ok = game.Put(row, col)
			if !ok {
				fmt.Println("その位置には置けません")
				continue
			}
			break
		}
	}
}
