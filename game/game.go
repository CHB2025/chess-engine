package game

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"

	"bareman.net/chess-engine/game/piece"
)

const (
	Forward    = 8
	Backward   = -8
	Right      = 1
	Left       = -1
	FrontRight = 9
	FrontLeft  = 7
	BackRight  = -7
	BackLeft   = -9
)

type Game struct {
	Board       [64]piece.Piece
	MoveCount   int
	HalfMove    int
	WhiteToMove bool
	WQCastle    bool
	WKCastle    bool
	BQCastle    bool
	BKCastle    bool
	EPTarget    int
}

func (g Game) String() string {
	result := ""
	for y := 7; y >= 0; y-- {
		result += fmt.Sprint(y+1) + " "
		for _, p := range g.Board[8*y : 8*y+7] {
			result += fmt.Sprintf("%s | ", p)
		}
		result += fmt.Sprintf("%s\n", g.Board[8*y+7])
	}
	result += "  a   b   c   d   e   f   g   h\n"
	//Testing data
	result += "\n"
	result += fmt.Sprintf("White to move: %v\n", g.WhiteToMove)
	result += fmt.Sprintf("Move: %v\n", g.MoveCount)
	result += fmt.Sprintf("HalfMove: %v\n", g.HalfMove)
	result += fmt.Sprintf("White Queenside Castle: %v\n", g.WQCastle)
	result += fmt.Sprintf("White Kingside Castle: %v\n", g.WKCastle)
	result += fmt.Sprintf("Black Queenside Castle: %v\n", g.BQCastle)
	result += fmt.Sprintf("Black Kingside Castle: %v\n", g.BKCastle)
	result += fmt.Sprintf("En Passant target square: %v\n", g.EPTarget)
	return result
}

func (g Game) Move(start, finish string) (Game, error) {
	sIndex := IndexFromPosition(start)
	fIndex := IndexFromPosition(finish)
	if sIndex == -1 || fIndex == -1 {
		return g, fmt.Errorf("Invalid position given. Received %v and %v\n", start, finish)
	}
	move := fIndex - sIndex
	validMoves := g.ValidMoves(sIndex)
	isValid := false
	for _, mv := range validMoves {
		if mv == move {
			isValid = true
			break
		}
	}
	if !isValid {
		return g, fmt.Errorf("Invalid move given. Received %v to %v.\n", start, finish)
	}

	g.Board[fIndex] = g.Board[sIndex]
	g.Board[sIndex] = piece.Empty

	return g, nil
}

func (g Game) GetPiece(position string) piece.Piece {
	index := IndexFromPosition(position)
	if index == -1 {
		fmt.Printf("Invalid position given.")
		return piece.Empty
	}
	return g.Board[index]
}

func (g Game) AllValidMoves() [][2]int {
	moves := [][2]int{}
	for i, p := range g.Board {
		color := piece.White
		if !g.WhiteToMove {
			color = piece.Black
		}
		if int(p)&color == color {
			mvs := g.ValidMoves(i)
			for _, mv := range mvs {
				moves = append(moves, [2]int{i, mv})
			}
		}
	}
	return moves
}

func (g Game) ValidMoves(start int) []int {
	p := g.Board[start]
	dir := -1 * (int(p/8)*2 - 3)
	var moves []int
	switch p % 8 {
	case piece.Pawn:
		index := start + Forward*dir
		if index < 64 && index >= 0 && g.Board[index] == piece.Empty {
			moves = append(moves, Forward*dir)

			if start/8-dir == 7 || (start)/8-dir == 0 { // in starting row
				moves = append(moves, Forward*2*dir)
			}
			if g.EPTarget != -1 && index+Left == g.EPTarget && g.WhiteToMove == (p/8 == 1) {
				moves = append(moves, Forward*dir+Left)
			}
			if g.EPTarget != -1 && index+Right == g.EPTarget && g.WhiteToMove == (p/8 == 1) {
				moves = append(moves, Forward*dir+Right)
			}
		}
	case piece.Queen:
		moves = append(moves, g.slidingMoves(p, start, Forward, -1)...)
		moves = append(moves, g.slidingMoves(p, start, Backward, -1)...)
		moves = append(moves, g.slidingMoves(p, start, Left, -1)...)
		moves = append(moves, g.slidingMoves(p, start, Right, -1)...)
		moves = append(moves, g.slidingMoves(p, start, FrontRight, -1)...)
		moves = append(moves, g.slidingMoves(p, start, FrontLeft, -1)...)
		moves = append(moves, g.slidingMoves(p, start, BackRight, -1)...)
		moves = append(moves, g.slidingMoves(p, start, BackLeft, -1)...)
	case piece.Bishop:
		moves = append(moves, g.slidingMoves(p, start, FrontRight, -1)...)
		moves = append(moves, g.slidingMoves(p, start, FrontLeft, -1)...)
		moves = append(moves, g.slidingMoves(p, start, BackRight, -1)...)
		moves = append(moves, g.slidingMoves(p, start, BackLeft, -1)...)
	case piece.Rook:
		moves = append(moves, g.slidingMoves(p, start, Forward, -1)...)
		moves = append(moves, g.slidingMoves(p, start, Backward, -1)...)
		moves = append(moves, g.slidingMoves(p, start, Left, -1)...)
		moves = append(moves, g.slidingMoves(p, start, Right, -1)...)
	case piece.King:
		moves = append(moves, g.slidingMoves(p, start, Forward, 1)...)
		moves = append(moves, g.slidingMoves(p, start, Backward, 1)...)
		moves = append(moves, g.slidingMoves(p, start, Left, 1)...)
		moves = append(moves, g.slidingMoves(p, start, Right, 1)...)
		moves = append(moves, g.slidingMoves(p, start, FrontRight, 1)...)
		moves = append(moves, g.slidingMoves(p, start, FrontLeft, 1)...)
		moves = append(moves, g.slidingMoves(p, start, BackRight, 1)...)
		moves = append(moves, g.slidingMoves(p, start, BackLeft, 1)...)
	case piece.Knight:
		preMoves := []int{
			Forward + FrontRight,
			Forward + FrontLeft,
			Left + FrontLeft,
			Left + BackLeft,
			Right + FrontRight,
			Right + BackRight,
			Backward + BackLeft,
			Backward + BackRight,
		}
		for _, mv := range preMoves {
			distance := int(math.Abs(float64(start/8-(start+mv)/8))) + int(math.Abs(float64(start%8-(start+mv)%8)))
			if start+mv >= 0 && start+mv < 64 && distance == 3 && g.Board[start+mv]/8 != p/8 {
				moves = append(moves, mv)
			}
		}
	}
	return moves
}

func (g Game) slidingMoves(p piece.Piece, start, dir, limit int) []int {
	var moves []int
	move := dir
	startRow := start / 8
	startCol := start % 8
	inBoard := start+move >= 0 && start+move < 64
	withinLimit := limit == -1 || len(moves) < limit
	crossesBoundary := math.Abs(float64(startRow-(start+move)/8))+math.Abs(float64(startCol-(start+move)%8)) > 2
	for inBoard && withinLimit && !crossesBoundary && g.Board[start+move]/8 != p/8 {
		moves = append(moves, move)
		if g.Board[start+move] != piece.Empty {
			break
		}
		crossesBoundary = math.Abs(float64((start+move+dir)/8-(start+move)/8))+math.Abs(float64((start+move+dir)%8-(start+move)%8)) <= 2
		move += dir
		inBoard = start+move >= 0 && start+move < 64
		withinLimit = limit == -1 || len(moves) < limit

	}
	return moves
}

func DefaultGame() Game {
	game, _ := GameFromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	return game
}

func GameFromFEN(fen string) (Game, error) {
	if !isValidFEN(fen) {
		return Game{}, fmt.Errorf("Invalid FEN string. Received %v", fen)
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

	game := Game{
		Board:       board,
		MoveCount:   move,
		HalfMove:    halfMove,
		WhiteToMove: sections[1] == "w",
		WKCastle:    strings.Contains(sections[2], "K"),
		WQCastle:    strings.Contains(sections[2], "Q"),
		BKCastle:    strings.Contains(sections[2], "k"),
		BQCastle:    strings.Contains(sections[2], "q"),
		EPTarget:    IndexFromPosition(sections[3]),
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

func IndexFromPosition(pos string) int {
	if len(pos) != 2 {
		return -1
	}
	row, err := strconv.Atoi(string(pos[1]))
	col := rune(pos[0])
	if err != nil || col < 'a' || col > 'h' || row < 1 || row > 8 {
		return -1
	}
	return (row-1)*8 + int(col-'a')
}

func PositionFromIndex(index int) string {
	if index >= 64 || index < 0 {
		return ""
	}
	col := index % 8
	row := index / 8
	return fmt.Sprintf("%v%v", string('a'+rune(col)), row+1)
}
