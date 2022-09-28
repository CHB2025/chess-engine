package game

import (
	"fmt"
	"sync"

	"bareman.net/chess-engine/game/move"
	"bareman.net/chess-engine/game/piece"
)

type Game struct {
	Mu          sync.Mutex
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
	hashKeys    [781]uint64
	Hash        uint64
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

func (g *Game) Score() int {
	score := 0
	for _, p := range g.Board {
		if p.IsWhite() == g.WhiteToMove {
			score += p.Score()
		} else {
			score -= p.Score()
		}
	}
	return 0
}

// Will ignore En-Passant
func (g *Game) Attackers(position string) []string {

	p := g.Piece(position)
	if p == piece.Empty {
		if g.WhiteToMove {
			p = piece.Black
		} else {
			p = piece.White
		}
	}

	atks := []string{}

	start := indexFromPosition(position)

	// Check pawns
	pMoves := g.pawnMoves(start, p.Color())
	for _, m := range pMoves {
		if g.Piece(m[2:4]).Type() == piece.Pawn {
			atks = append(atks, m[2:4])
		}
	}

	// Check king
	kMoves := g.kingMoves(start, p.Color())
	for _, m := range kMoves {
		if g.Piece(m[2:4]).Type() == piece.King {
			atks = append(atks, m[2:4])
		}
	}

	// Check bishop/half queen
	bMoves := g.bishopMoves(start, p.Color())
	for _, m := range bMoves {
		if g.Piece(m[2:4]).Type() == piece.Bishop || g.Piece(m[2:4]).Type() == piece.Queen {
			atks = append(atks, m[2:4])
		}
	}

	// Check rook/other half queen
	rMoves := g.rookMoves(start, p.Color())
	for _, m := range rMoves {
		if g.Piece(m[2:4]).Type() == piece.Rook || g.Piece(m[2:4]).Type() == piece.Queen {
			atks = append(atks, m[2:4])
		}
	}

	// Check knight
	nMoves := g.knightMoves(start, p.Color())
	for _, m := range nMoves {
		if g.Piece(m[2:4]).Type() == piece.Knight {
			atks = append(atks, m[2:4])
		}
	}

	return atks
}

// Will Ignore EnPassant
func (g *Game) IsAttacked(position string) bool {

	p := g.Piece(position)
	if p == piece.Empty {
		if g.WhiteToMove {
			p = piece.Black
		} else {
			p = piece.White
		}
	}

	start := indexFromPosition(position)

	// Check pawns
	pMoves := g.pawnMoves(start, p.Color())
	for _, m := range pMoves {
		if g.Piece(m[2:4]).Type() == piece.Pawn {
			return true
		}
	}

	// Check king
	kMoves := g.kingMoves(start, p.Color())
	for _, m := range kMoves {
		if g.Piece(m[2:4]).Type() == piece.King {
			return true
		}
	}

	// Check bishop/half queen
	bMoves := g.bishopMoves(start, p.Color())
	for _, m := range bMoves {
		if g.Piece(m[2:4]).Type() == piece.Bishop || g.Piece(m[2:4]).Type() == piece.Queen {
			return true
		}
	}

	// Check rook/other half queen
	rMoves := g.rookMoves(start, p.Color())
	for _, m := range rMoves {
		if g.Piece(m[2:4]).Type() == piece.Rook || g.Piece(m[2:4]).Type() == piece.Queen {
			return true
		}
	}

	// Check knight
	nMoves := g.knightMoves(start, p.Color())
	for _, m := range nMoves {
		if g.Piece(m[2:4]).Type() == piece.Knight {
			return true
		}
	}

	return false
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
	transpositions := make(map[uint64]int)
	return g.perft(depth, transpositions)
}

func (g *Game) perft(depth int, transpositions map[uint64]int) int {
	if depth == 0 {
		return 1
	}

	moves := g.AllLegalMoves()
	var moveCount int
	for _, mv := range moves {
		m, _ := move.EmptyMove(mv)
		g.make(m)

		mvs, ok := transpositions[g.Hash]
		if !ok {
			mvs = g.perft(depth-1, transpositions)
			transpositions[g.Hash] = mvs
		}
		moveCount += mvs
		g.Unmake()
	}
	return moveCount
}

func (g *Game) DividedPerft(depth int) map[string]int {
	if depth == 0 {
		return make(map[string]int)
	}
	transpositions := make(map[uint64]int)
	moves := g.PseudoLegalMoves("")
	results := make(map[string]int)
	for _, mv := range moves {
		err := g.Make(mv)
		if err != nil {
			continue
		}
		results[mv] = g.perft(depth-1, transpositions)
		g.Unmake()
	}
	return results
}
