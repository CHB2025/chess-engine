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

func (g *Game) IsMoveLegal(mv string) bool {
	m, err := move.EmptyMove(mv)
	if err != nil {
		return false
	}
	p := g.Piece(m.Origin)
	if p.IsWhite() != g.WhiteToMove {
		return false
	}
	color := piece.White
	if !g.WhiteToMove {
		color = piece.Black
	}

	var king string
	for i, p := range g.Board {
		if p == piece.Piece(color|piece.King) {
			king = positionFromIndex(i)
		}
	}

	g.make(m)
	inCheck := false
	if p == piece.Piece(color|piece.King) {
		inCheck = g.IsAttacked(m.Dest)

		oRow, _ := coordinates(m.OriginIndex())
		dRow, _ := coordinates(m.DestIndex())
		// Checks if castling out of/over check
		if !inCheck && mdistance(m.OriginIndex(), m.DestIndex()) == 2 && oRow == dRow {
			position := positionFromIndex(m.OriginIndex() + (m.DestIndex()-m.OriginIndex())/2)
			inCheck = g.IsAttacked(position) || g.IsAttacked(m.Origin)
		}
	} else {
		inCheck = g.IsAttacked(king)
	}
	g.Unmake()
	return !inCheck
}

func (g *Game) AllLegalMoves() []string {
	return g.LegalMoves("")
}

func (g *Game) LegalMoves(pos string) []string {
	var moves, search []string
	if pos != "" {
		search = []string{pos}
	} else {
		search = make([]string, 64)
		for i := 0; i < 64; i++ {
			search[i] = positionFromIndex(i)
		}
	}

	for i, ps := range search {
		if g.Board[i] == piece.Empty || g.Board[i].IsWhite() != g.WhiteToMove {
			continue
		}
		mvs := g.moves(ps)
		for _, mv := range mvs {
			if g.IsMoveLegal(mv) {
				moves = append(moves, mv)
			}
		}
	}
	return moves
}

func (g *Game) PseudoLegalMoves(pos string) []string {
	var moves []string
	search := make(map[int]string, 64)
	if pos != "" {
		search[indexFromPosition(pos)] = pos
	} else {
		for i := 0; i < 64; i++ {
			if g.Board[i] != piece.Empty && g.Board[i].IsWhite() == g.WhiteToMove {
				search[i] = positionFromIndex(i)
			}
		}
	}

	for i, ps := range search {
		if g.Board[i] == piece.Empty || g.Board[i].IsWhite() != g.WhiteToMove {
			continue
		}
		moves = append(moves, g.moves(ps)...)
	}
	return moves
}

func (g *Game) moves(pos string) []string {
	p := g.Piece(pos)
	start := indexFromPosition(pos)

	switch p.Type() {
	case piece.Pawn:
		return g.pawnMoves(start, p.Color())
	case piece.Queen:
		return g.queenMoves(start, p.Color())
	case piece.Bishop:
		return g.bishopMoves(start, p.Color())
	case piece.Rook:
		return g.rookMoves(start, p.Color())
	case piece.King:
		return g.kingMoves(start, p.Color())
	case piece.Knight:
		return g.knightMoves(start, p.Color())
	default:
		return []string{}
	}
}

func (g *Game) pawnMoves(start int, color piece.Piece) []string {
	dir := -1 * (int(color/piece.White)*2 - 3)
	index := start + Forward*dir
	moves := []string{}

	appendMoves := func(dest int) {
		mv := positionFromIndex(start) + positionFromIndex(dest)
		row, _ := coordinates(dest)
		//Handles promotions
		if row == 7 || row == 0 {
			promos := []string{"Q", "R", "B", "N"}
			if color != piece.White {
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
		appendMoves(index)

		row, _ := coordinates(start)

		// in starting row and the two spots in front are open
		if (row-dir == 7 || row-dir == 0) && g.Board[index+Forward*dir] == piece.Empty {
			appendMoves(index + Forward*dir)
		}
	}
	// En Passant
	if g.EPTarget != -1 && index+Left == g.EPTarget && mdistance(index, index+Left) == 1 {
		appendMoves(index + Left)
	}
	if g.EPTarget != -1 && index+Right == g.EPTarget && mdistance(index, index+Right) == 1 {
		appendMoves(index + Right)
	}
	//Attacking squares
	leftAttack := start + Forward*dir + Left
	if leftAttack < 64 && leftAttack >= 0 &&
		mdistance(start, leftAttack) == 2 &&
		g.Board[leftAttack] != piece.Empty &&
		g.Board[leftAttack].Color() != color {
		appendMoves(leftAttack)
	}
	rightAttack := start + Forward*dir + Right
	if rightAttack < 64 && rightAttack >= 0 &&
		mdistance(start, rightAttack) == 2 &&
		g.Board[rightAttack] != piece.Empty &&
		g.Board[rightAttack].Color() != color {
		appendMoves(rightAttack)
	}
	return moves
}

func (g *Game) kingMoves(start int, color piece.Piece) []string {
	allDirections := []int{
		FrontLeft, Forward, FrontRight,
		Left, Right,
		BackLeft, Backward, BackRight,
	}

	moves := []string{}

	for _, dir := range allDirections {
		if start+dir >= 0 && start+dir < 64 && mdistance(start, start+dir) <= 2 && g.Board[start+dir].Color() != color {
			targetPosition := positionFromIndex(start + dir)
			moves = append(moves, positionFromIndex(start)+targetPosition)
		}
	}
	row, col := coordinates(start)
	// Castling
	if col == 4 && (color.IsWhite() && g.WKCastle && row == 0 || !color.IsWhite() && g.BKCastle && row == 7) {
		kingSideClear := true
		for i := 1; i <= 2; i++ {
			if g.Board[start+i*Right] != piece.Empty {
				kingSideClear = false
				break
			}
		}
		if kingSideClear && g.Board[start+3*Right] == piece.Piece(color|piece.Rook) {
			moves = append(moves, positionFromIndex(start)+positionFromIndex(start+Right*2))
		}
	}
	if col == 4 && (color.IsWhite() && g.WQCastle && row == 0 || !color.IsWhite() && g.BQCastle && row == 7) {
		queenSideClear := true
		for i := 1; i <= 3; i++ {
			if g.Board[start+i*Left] != piece.Empty {
				queenSideClear = false
				break
			}
		}
		if queenSideClear && g.Board[start+4*Left] == piece.Piece(color|piece.Rook) {
			moves = append(moves, positionFromIndex(start)+positionFromIndex(start+Left*2))
		}
	}

	return moves
}

func (g *Game) bishopMoves(start int, color piece.Piece) []string {
	diagonals := []int{
		FrontLeft, FrontRight,
		BackLeft, BackRight,
	}

	moves := []string{}
	for _, dir := range diagonals {
		moves = append(moves, g.slidingMoves(start, dir, color)...)
	}
	return moves
}

func (g *Game) rookMoves(start int, color piece.Piece) []string {
	orthogonals := []int{
		Forward, Backward,
		Right, Left,
	}

	moves := []string{}
	for _, dir := range orthogonals {
		moves = append(moves, g.slidingMoves(start, dir, color)...)
	}
	return moves
}

func (g *Game) queenMoves(start int, color piece.Piece) []string {
	return append(g.rookMoves(start, color), g.bishopMoves(start, color)...)
}

func (g *Game) knightMoves(start int, color piece.Piece) []string {
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

	moves := []string{}
	for _, mv := range preMoves {
		distance := mdistance(start, start+mv)
		if start+mv >= 0 && start+mv < 64 && distance == 3 && g.Board[start+mv].Color() != color {
			moves = append(moves, positionFromIndex(start)+positionFromIndex(start+mv))
		}
	}
	return moves
}

func (g *Game) slidingMoves(start, dir int, color piece.Piece) []string {
	var moves []string
	curr := start
	inBoard := curr+dir >= 0 && curr+dir < 64
	crossesBoundary := mdistance(curr, curr+dir) > 2
	for inBoard && !crossesBoundary && (g.Board[curr+dir] == piece.Empty || g.Board[curr+dir].IsWhite() != color.IsWhite()) {
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
