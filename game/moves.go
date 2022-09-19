package game

import (
	"bareman.net/chess-engine/game/move"
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

func (g *Game) AllValidMoves() []string {
	return g.ValidMoves("")
}

func (g *Game) ValidMoves(pos string) []string {
	var moves []string
	var king string
	color := piece.White
	if !g.WhiteToMove {
		color = piece.Black
	}
	search := make(map[int]piece.Piece, 16)

	for i, p := range g.Board {
		if p == piece.Piece(color|piece.King) {
			king = positionFromIndex(i)
		}
		if p != piece.Empty && p.IsWhite() == g.WhiteToMove {
			search[i] = p
		}
	}

	if pos != "" {
		index := indexFromPosition(pos)
		search = make(map[int]piece.Piece, 1)
		search[index] = g.Board[index]
	}

	// Making each move to see if it will put king in check. Not very efficient, since every move
	// has to check all opposing pieces to see if they can move to the kings square
	for i, p := range search {
		mvs := g.moves(positionFromIndex(i))
		for _, mv := range mvs {
			m, _ := move.EmptyMove(mv)
			g.make(m)
			var atks []string
			if p == piece.Piece(color|piece.King) {
				atks = g.attackers(m.Dest)

				// Checks if castling out of/over check
				if len(atks) == 0 && mdistance(m.OriginIndex(), m.DestIndex()) == 2 && m.OriginIndex()/8 == m.DestIndex()/8 {
					position := positionFromIndex(m.OriginIndex() + (m.DestIndex()-m.OriginIndex())/2)
					atks = append(atks, g.attackers(position)...)
					atks = append(atks, g.attackers(m.Origin)...)
				}
			} else {
				atks = g.attackers(king)
			}
			if len(atks) == 0 {
				moves = append(moves, mv)
			}
			g.Unmake()
		}
	}
	return moves
}

func (g *Game) moves(pos string) []string {
	p := g.Piece(pos)
	start := indexFromPosition(pos)
	var moves []string

	allDirections := []int{
		FrontLeft, Forward, FrontRight,
		Left, Right,
		BackLeft, Backward, BackRight,
	}
	diagonals := []int{
		FrontLeft, FrontRight,
		BackLeft, BackRight,
	}
	orthogonals := []int{
		Forward, Backward,
		Right, Left,
	}

	switch p % 8 {
	case piece.Pawn:
		dir := -1 * (int(p/piece.White)*2 - 3)
		index := start + Forward*dir

		appendPawnMove := func(dest int) {
			mv := pos + positionFromIndex(dest)

			//Handles promotions
			if dest/8 == 7 || dest/8 == 0 {
				promos := []string{"Q", "R", "B", "N"}
				if !p.IsWhite() {
					promos = []string{"q", "r", "b", "n"}
				}
				for _, p := range promos {
					moves = append(moves, mv+p)
				}
			} else {
				moves = append(moves, mv)
			}
		}

		if index < 64 && index >= 0 && g.Board[index] == piece.Empty {
			appendPawnMove(index)

			// in starting row and the two spots in front are open
			if (start/8-dir == 7 || start/8-dir == 0) && g.Board[index+Forward*dir] == piece.Empty {
				moves = append(moves, pos+positionFromIndex(index+Forward*dir))
			}
		}
		// En Passant
		if g.EPTarget != -1 && index+Left == g.EPTarget && mdistance(index, index+Left) == 1 {
			moves = append(moves, pos+positionFromIndex(index+Left))
		}
		if g.EPTarget != -1 && index+Right == g.EPTarget && mdistance(index, index+Right) == 1 {
			moves = append(moves, pos+positionFromIndex(index+Right))
		}
		//Attacking squares
		leftAttack := start + Forward*dir + Left
		if leftAttack < 64 && leftAttack >= 0 &&
			mdistance(start, leftAttack) == 2 &&
			g.Board[leftAttack] != piece.Empty &&
			g.Board[leftAttack].IsWhite() != p.IsWhite() {
			appendPawnMove(leftAttack)
		}
		rightAttack := start + Forward*dir + Right
		if rightAttack < 64 && rightAttack >= 0 &&
			mdistance(start, rightAttack) == 2 &&
			g.Board[rightAttack] != piece.Empty &&
			g.Board[rightAttack].IsWhite() != p.IsWhite() {
			appendPawnMove(rightAttack)
		}
	case piece.Queen:
		for _, dir := range allDirections {
			moves = append(moves, g.slidingMoves(p, start, dir)...)
		}
	case piece.Bishop:
		for _, dir := range diagonals {
			moves = append(moves, g.slidingMoves(p, start, dir)...)
		}
	case piece.Rook:
		for _, dir := range orthogonals {
			moves = append(moves, g.slidingMoves(p, start, dir)...)
		}
	case piece.King:
		for _, dir := range allDirections {
			if start+dir >= 0 && start+dir < 64 && mdistance(start, start+dir) <= 2 && g.Board[start+dir]/8 != p/8 {
				targetPosition := positionFromIndex(start + dir)
				moves = append(moves, pos+targetPosition)
			}
		}
		color := p / piece.White * piece.White
		if p.IsWhite() && g.WKCastle || !p.IsWhite() && g.BKCastle {
			kingSideClear := true
			for i := 1; i <= 2; i++ {
				if g.Board[start+i*Right] != piece.Empty {
					kingSideClear = false
					break
				}
			}
			if kingSideClear && g.Board[start+3*Right] == piece.Piece(color|piece.Rook) {
				moves = append(moves, pos+positionFromIndex(start+Right*2))
			}
		}
		if p.IsWhite() && g.WQCastle || !p.IsWhite() && g.BQCastle {
			queenSideClear := true
			for i := 1; i <= 3; i++ {
				if g.Board[start+i*Left] != piece.Empty {
					queenSideClear = false
					break
				}
			}
			if queenSideClear && g.Board[start+4*Left] == piece.Piece(color|piece.Rook) {
				moves = append(moves, pos+positionFromIndex(start+Left*2))
			}
		}
	case piece.Knight:
		preMoves := []int{
			Forward + Forward + Right,
			Forward + Forward + Left,
			Forward + Left + Left,
			Backward + Left + Left,
			Forward + Right + Right,
			Backward + Right + Right,
			Backward + Backward + Left,
			Backward + Backward + Right,
		}
		for _, mv := range preMoves {
			distance := mdistance(start, start+mv)
			if start+mv >= 0 && start+mv < 64 && distance == 3 && g.Board[start+mv]/8 != p/8 {
				moves = append(moves, pos+positionFromIndex(start+mv))
			}
		}
	}
	return moves
}

func (g *Game) slidingMoves(p piece.Piece, start, dir int) []string {
	var moves []string
	curr := start
	inBoard := curr+dir >= 0 && curr+dir < 64
	crossesBoundary := mdistance(curr, curr+dir) > 2
	for inBoard && !crossesBoundary && (g.Board[curr+dir] == piece.Empty || g.Board[curr+dir].IsWhite() != p.IsWhite()) {
		moves = append(moves, positionFromIndex(start)+positionFromIndex(curr+dir))
		if g.Board[curr+dir] != piece.Empty {
			break
		}

		curr = curr + dir
		crossesBoundary = mdistance(curr, curr+dir) > 2
		inBoard = curr+dir >= 0 && curr+dir < 64
	}
	return moves
}
