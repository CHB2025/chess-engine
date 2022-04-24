package board

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Board [8][8]*Piece

func (b Board) String() string {
	result := ""
	for y, row := range b {
		result += fmt.Sprint(8-y) + " "
		for _, p := range row[:len(row)-1] {
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
	result += "  a   b   c   d   e   f   g   h\n"
	return result
}

func (b Board) GetPiece(pos Position) *Piece {

	return b[8-pos.row][pos.col-'a']
}

func (b Board) GetPosition(piece *Piece) *Position {
	var pos *Position
	for y, row := range b {
		for x, piece2 := range row {
			if piece2 == piece {
				pos = &Position{row: 8 - y, col: 'a' + rune(x)}
				break
			}
			if pos != nil {
				break
			}
		}
	}
	return pos
}

func (b Board) Move(start, finish Position) (Board, error) {
	piece := b.GetPiece(start)
	if piece == nil {
		return b, fmt.Errorf("No piece at chosen position")
	}
	validMoves := piece.ValidMoves(b)
	isValid := false
	for _, mv := range validMoves {
		if mv == finish {
			isValid = true
			break
		}
	}
	if !isValid {
		return b, fmt.Errorf("Invalid move. Valid moves for %v are %v", start, validMoves)
	}
	b[8-finish.row][finish.col-'a'] = b[8-start.row][start.col-'a']
	b[8-start.row][start.col-'a'] = nil
	return b, nil
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
