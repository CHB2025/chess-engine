package game

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"bareman.net/chess-engine/game/piece"
)

func Default() *Game {
	game, _ := FromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	return game
}

func FromFEN(fen string) (*Game, error) {
	if !isValidFEN(fen) {
		return nil, fmt.Errorf("Invalid FEN string. Received %v", fen)
	}
	sections := strings.Split(fen, " ")
	var board [64]piece.Piece
	for y, row := range strings.Split(sections[0], "/") {
		var offset int
		for x, symbol := range row {
			if unicode.IsNumber(symbol) {
				num, _ := strconv.Atoi(string(symbol))
				offset += num - 1
				continue
			}

			board[8*(7-y)+x+offset] = piece.FromRune(symbol)
		}
	}

	move, _ := strconv.Atoi(sections[5])
	halfMove, _ := strconv.Atoi(sections[4])

	game := &Game{
		Board:       board,
		MoveCount:   move,
		HalfMove:    halfMove,
		WhiteToMove: sections[1] == "w",
		WKCastle:    strings.Contains(sections[2], "K"),
		WQCastle:    strings.Contains(sections[2], "Q"),
		BKCastle:    strings.Contains(sections[2], "k"),
		BQCastle:    strings.Contains(sections[2], "q"),
		EPTarget:    indexFromPosition(sections[3]),
	}
	return game, nil
}

func isValidFEN(fen string) bool {
	sections := strings.Split(fen, " ")
	if len(sections) != 6 {
		return false
	}

	rows := strings.Split(sections[0], "/")
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
			// Should check if it is one of rnbqkp
			posCount++
		}
	}
	return posCount == 64
	//Should check the rest of the sections too
}
