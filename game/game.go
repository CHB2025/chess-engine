package game

import (
	"fmt"

	"bareman.net/chess-engine/game/move"
	"bareman.net/chess-engine/game/piece"
)

type Game struct {
	Board       [64]piece.Piece
	Moves       []*move.Move
	MoveCount   int
	HalfMove    int
	WhiteToMove bool
	WQCastle    bool
	WKCastle    bool
	BQCastle    bool
	BKCastle    bool
	EPTarget    int
}

func (g *Game) String() string {
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

func (g *Game) Piece(position string) piece.Piece {
	index := indexFromPosition(position)
	if index == -1 {
		// fmt.Printf("Invalid position given.")
		return piece.Empty
	}
	return g.Board[index]
}

func (g *Game) attackers(position string) []string {
	p := g.Piece(position)
	if p == piece.Empty {
		if g.WhiteToMove {
			p = piece.Black
		} else {
			p = piece.White
		}
	}
	var moves []string
	for index, pi := range g.Board {
		if pi != piece.Empty && pi/piece.White != p/piece.White {
			moves = append(moves, g.moves(positionFromIndex(index))...)
		}
	}
	var attackers []string
	for _, mv := range moves {
		if mv[2:4] == position {
			attackers = append(attackers, mv[:2])
		}
	}
	return attackers
}

func (g *Game) ToFEN() string {
	var boardString, playerToMove, castlingRights, epPosition string
	var numEmptySquares int
	for y := 7; y >= 0; y-- {
		for x := 0; x < 8; x++ {
			p := g.Board[8*y+x]
			if p == piece.Empty {
				numEmptySquares++
				continue
			}
			if numEmptySquares > 0 {
				boardString += fmt.Sprint(numEmptySquares)
				numEmptySquares = 0
			}
			boardString += fmt.Sprint(p)
		}
		if numEmptySquares > 0 {
			boardString += fmt.Sprint(numEmptySquares)
			numEmptySquares = 0
		}

		if y != 0 {
			boardString += "/"
		}
	}
	if g.WhiteToMove {
		playerToMove = "w"
	} else {
		playerToMove = "b"
	}
	if g.WKCastle {
		castlingRights += "K"
	}
	if g.WQCastle {
		castlingRights += "Q"
	}
	if g.BKCastle {
		castlingRights += "k"
	}
	if g.BQCastle {
		castlingRights += "q"
	}
	if castlingRights == "" {
		castlingRights = "-"
	}
	epPosition = positionFromIndex(g.EPTarget)
	if epPosition == "" {
		epPosition = "-"
	}

	return fmt.Sprintf("%v %v %v %v %v %v", boardString, playerToMove, castlingRights, epPosition, g.HalfMove, g.MoveCount)
}

func (g *Game) Perft(depth int) int {
	if depth == 0 {
		return 1
	}
	moves := g.AllValidMoves()
	var moveCount int
	for _, mv := range moves {
		m, _ := move.EmptyMove(mv)
		g.make(m)
		moveCount += g.Perft(depth - 1)
		g.Unmake()
	}
	return moveCount
}

func (g *Game) DividedPerft(depth int) map[string]int {
	if depth == 0 {
		return make(map[string]int)
	}
	moves := g.AllValidMoves()
	results := make(map[string]int)
	for _, mv := range moves {
		m, _ := move.EmptyMove(mv)
		g.make(m)
		results[mv] = g.Perft(depth - 1)
		g.Unmake()
	}
	return results
}
