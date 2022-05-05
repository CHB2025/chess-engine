package piece

import "unicode"

const (
	Empty  = 0
	King   = 1
	Pawn   = 2
	Knight = 3
	Bishop = 4
	Rook   = 5
	Queen  = 6

	White = 8
	Black = 16
)

type Piece uint8

func (p Piece) String() string {
	var sym rune
	switch p % 8 {
	case Empty:
		sym = ' '
	case King:
		sym = 'k'
	case Pawn:
		sym = 'p'
	case Knight:
		sym = 'n'
	case Bishop:
		sym = 'b'
	case Rook:
		sym = 'r'
	case Queen:
		sym = 'q'
	}
	if p/8 == 1 {
		sym = unicode.ToUpper(sym)
	}
	return string(sym)
}

func FromRune(r rune) Piece {
	var p Piece
	switch unicode.ToLower(r) {
	case 'k':
		p = King
	case 'p':
		p = Pawn
	case 'n':
		p = Knight
	case 'b':
		p = Bishop
	case 'r':
		p = Rook
	case 'q':
		p = Queen
	}
	if unicode.IsUpper(r) {
		p = p | White
	} else {
		p = p | Black
	}
	return p
}

// func (p *Piece) IsValidMove(board Board, move Position) bool {
// 	moves := p.ValidMoves(board)
// 	for _, mv := range moves {
// 		if mv == move {
// 			return true
// 		}
// 	}
// 	return false
// }

// func (p *Piece) ValidMoves(board Board) []Position {
// 	moves := []Position{}
// 	pos := board.GetPosition(p)
// 	if pos == nil {
// 		return moves
// 	}

// 	switch p.Symbol {
// 	case 'p':
// 		dir := 1
// 		if !p.IsWhite {
// 			dir = -1
// 		}
// 		forward := Position{row: pos.row + dir, col: pos.col}
// 		if board.GetPiece(forward) == nil {
// 			moves = append(moves, forward)
// 			if pos.row-dir == 1 || pos.row-dir == 8 {
// 				double := Position{row: pos.row + 2*dir, col: pos.col}
// 				if board.GetPiece(double) == nil {
// 					moves = append(moves, double)
// 				}
// 			}
// 		}
// 	case 'r':
// 		moves = append(moves, movesFrom(board, *pos, 1, 0, -1)...)
// 		moves = append(moves, movesFrom(board, *pos, -1, 0, -1)...)
// 		moves = append(moves, movesFrom(board, *pos, 0, 1, -1)...)
// 		moves = append(moves, movesFrom(board, *pos, 0, -1, -1)...)
// 	case 'b':
// 		moves = append(moves, movesFrom(board, *pos, 1, 1, -1)...)
// 		moves = append(moves, movesFrom(board, *pos, 1, -1, -1)...)
// 		moves = append(moves, movesFrom(board, *pos, -1, 1, -1)...)
// 		moves = append(moves, movesFrom(board, *pos, -1, -1, -1)...)
// 	case 'q':
// 		for y := -1; y < 2; y++ {
// 			for x := -1; x < 2; x++ {
// 				moves = append(moves, movesFrom(board, *pos, y, x, -1)...)
// 			}
// 		}
// 	case 'k':
// 		for y := -1; y < 2; y++ {
// 			for x := -1; x < 2; x++ {
// 				moves = append(moves, movesFrom(board, *pos, y, x, 1)...)
// 			}
// 		}
// 	case 'n':
// 		moves = append(moves, movesFrom(board, *pos, 1, 2, 1)...)
// 		moves = append(moves, movesFrom(board, *pos, 2, 1, 1)...)
// 		moves = append(moves, movesFrom(board, *pos, -1, 2, 1)...)
// 		moves = append(moves, movesFrom(board, *pos, -2, 1, 1)...)
// 		moves = append(moves, movesFrom(board, *pos, -2, -1, 1)...)
// 		moves = append(moves, movesFrom(board, *pos, -1, -2, 1)...)
// 		moves = append(moves, movesFrom(board, *pos, 1, -2, 1)...)
// 		moves = append(moves, movesFrom(board, *pos, 2, -1, 1)...)
// 	}
// 	return moves
// }

// func (p *Piece) IsUnderAttack(board Board) bool {
// 	pos := board.GetPosition(p)
// 	if pos == nil {
// 		return false
// 	}
// 	for _, row := range board {
// 		for _, piece := range row {
// 			if piece != nil && piece.IsWhite != p.IsWhite && piece.IsValidMove(board, *pos) {
// 				return true
// 			}
// 		}
// 	}
// 	return false
// }

// func movesFrom(board Board, start Position, rInc, cInc, limit int) []Position {
// 	piece := board.GetPiece(start)
// 	freeSpaces := []Position{}
// 	pos := Position{start.row + rInc, start.col + rune(cInc)}

// 	isWithinBoard := 'a' <= pos.col && pos.col <= 'h' && 1 <= pos.row && pos.row <= 8
// 	isFreeSpace := isWithinBoard && (board.GetPiece(pos) == nil || board.GetPiece(pos).IsWhite != piece.IsWhite)
// 	isWithinLimit := limit == -1 || len(freeSpaces) < limit

// 	for isWithinBoard && isFreeSpace && isWithinLimit {
// 		freeSpaces = append(freeSpaces, pos)
// 		if board.GetPiece(pos) != nil {
// 			break
// 		}
// 		pos = Position{pos.row + rInc, pos.col + rune(cInc)}
// 		isWithinBoard = 'a' <= pos.col && pos.col <= 'h' && 1 <= pos.row && pos.row <= 8
// 		isFreeSpace = isWithinBoard && (board.GetPiece(pos) == nil || board.GetPiece(pos).IsWhite != piece.IsWhite)
// 		isWithinLimit = limit == -1 || len(freeSpaces) < limit
// 	}

// 	return freeSpaces
// }
