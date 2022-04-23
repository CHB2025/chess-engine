package board

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Board [8][8]*Piece

func (b *Board) String() string {
	result := ""
	for _, row := range b {
		for _, p := range row[:len(row)-2] {
			if p != nil {
				result += fmt.Sprintf("%s | ", p)
			} else {
				result += "  | "
			}
		}
		if row[len(row)-1] != nil {
			result += fmt.Sprintf("%s\n", row[len(row)-1])
		} else {
			result += " \n"
		}
	}
	return result
}

func (b *Board) Move(move string) error {
	// convert move to two positions, get piece at first position, check if second position is a valid move for that piece, then
	return nil
}

func DefaultBoard() Board {
	board, _ := BoardFromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR")
	return board
}

func BoardFromFEN(str string) (Board, error) {
	if !isValidFEN(str) {
		return Board{}, fmt.Errorf("Invalid FEN string. Received %v", str)
	}
	board := Board{}
	for y, row := range strings.Split(str, "/") {
		var offset int
		for x, symbol := range row {
			if unicode.IsNumber(symbol) {
				num, _ := strconv.Atoi(string(symbol))
				offset += num - 1
				continue
			}

			board[y][x+offset] = &Piece{
				Symbol:  unicode.ToLower(symbol),
				IsWhite: unicode.IsUpper(symbol),
			}
		}
	}
	return board, nil
}

func isValidFEN(fen string) bool {
	rows := strings.Split(fen, "/")
	if len(rows) != 8 {
		return false
	}

	posCount := 0
	for _, row := range rows {
		for _, symbol := range row {
			if unicode.IsNumber(symbol) {
				num, err := strconv.Atoi(string(symbol))
				if err != nil {
					return false
				}
				posCount += num
				continue
			}
			// Should check if it is one of rnbqkbnrp
			posCount++
		}
	}
	return posCount == 64
}
