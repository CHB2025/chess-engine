package game

import (
	"fmt"

	"bareman.net/chess-engine/game/move"
	"bareman.net/chess-engine/game/piece"
)

func (g *Game) Make(mv string) error {
	move, err := move.EmptyMove(mv)
	if err != nil {
		return err
	}

	isValid := g.IsMoveLegal(mv)
	if !isValid {
		return fmt.Errorf("Invalid move given. Received %v\n", mv)
	}
	g.make(move)

	return nil
}

func (g *Game) make(mv *move.Move) {
	oRow, oCol := coordinates(mv.OriginIndex())
	dRow, dCol := coordinates(mv.DestIndex())
	p := g.Piece(mv.Origin)
	capture := g.Piece(mv.Dest)
	castle := p.Type() == piece.King && oRow == dRow && mdistance(mv.OriginIndex(), mv.DestIndex()) > 1              // piece is king, moving 2 or 3 spaces on one rank
	ep := p.Type() == piece.Pawn && g.EPTarget == mv.DestIndex() && mdistance(mv.OriginIndex(), mv.DestIndex()) == 2 //piece is pawn, moving to target square diagonally

	mv.Capture, mv.Castle, mv.EnPassant = capture, castle, ep
	mv.BoardState = struct {
		WQCastle bool
		WKCastle bool
		BQCastle bool
		BKCastle bool
		EPTarget int
	}{
		WQCastle: g.WQCastle,
		WKCastle: g.WKCastle,
		BQCastle: g.BKCastle,
		BKCastle: g.BKCastle,
		EPTarget: g.EPTarget,
	}

	if p.Type() == piece.Pawn && oCol == dCol && mdistance(mv.OriginIndex(), mv.DestIndex()) == 2 {
		g.EPTarget = mv.DestIndex()/2 + mv.OriginIndex()/2 + mv.DestIndex()%2
	} else {
		g.EPTarget = -1
	}

	g.Board[mv.DestIndex()] = g.Board[mv.OriginIndex()]
	g.Board[mv.OriginIndex()] = piece.Empty
	g.WhiteToMove = !g.WhiteToMove
	if mv.Promotion != piece.Empty {
		g.Board[mv.DestIndex()] = mv.Promotion
	}

	if ep {
		file := mv.Dest[0]
		rank := mv.Origin[1]
		capIndex := indexFromPosition(string(file) + string(rank))
		mv.Capture = g.Board[capIndex]
		g.Board[capIndex] = piece.Empty
	}
	if castle {
		rank := mv.Origin[1]
		file := byte(int(mv.Origin[0]) + (int(mv.Dest[0])-int(mv.Origin[0]))/2)
		rePos := string(file) + string(rank)
		rsPos := "a" + string(rank)
		if file > mv.Origin[0] {
			rsPos = "h" + string(rank)
		}
		reIndex := indexFromPosition(rePos)
		rsIndex := indexFromPosition(rsPos)
		g.Board[reIndex] = g.Board[rsIndex]
		g.Board[rsIndex] = piece.Empty
	}
	if p.Type() == piece.King {
		// Doesn't handle rooks moving
		if p.IsWhite() {
			g.WKCastle = false
			g.WQCastle = false
		} else {
			g.BKCastle = false
			g.BQCastle = false
		}
	}
	g.Moves = append(g.Moves, mv)
	g.MoveCount += 1
	g.incrementHash(mv, p)
}

func (g *Game) Unmake() {

	move := g.Moves[len(g.Moves)-1]
	if move.Promotion == piece.Empty {
		g.incrementHash(move, g.Board[move.DestIndex()])
	} else {
		g.incrementHash(move, piece.Pawn|move.Promotion.Color())
	}
	g.Moves = g.Moves[:len(g.Moves)-1]
	g.MoveCount -= 1

	g.Board[move.OriginIndex()] = g.Board[move.DestIndex()]
	g.Board[move.DestIndex()] = move.Capture
	g.WhiteToMove = !g.WhiteToMove
	g.EPTarget = move.BoardState.EPTarget
	g.WQCastle = move.BoardState.WQCastle
	g.WKCastle = move.BoardState.WKCastle
	g.BKCastle = move.BoardState.BKCastle
	g.BQCastle = move.BoardState.BQCastle
	if move.Promotion != piece.Empty {
		g.Board[move.OriginIndex()] = piece.Pawn | move.Promotion.Color()
	}
	if move.EnPassant {
		g.Board[move.DestIndex()] = piece.Empty
		file := move.Dest[0]
		rank := move.Origin[1]
		index := indexFromPosition(string(file) + string(rank))
		g.Board[index] = move.Capture
	}
	if move.Castle {
		rank := move.Origin[1]
		file := byte(int(move.Origin[0]) + (int(move.Dest[0])-int(move.Origin[0]))/2)
		rePos := string(file) + string(rank)
		rsPos := "a" + string(rank)
		if file > move.Origin[0] {
			rsPos = "h" + string(rank)
		}
		reIndex := indexFromPosition(rePos)
		rsIndex := indexFromPosition(rsPos)
		g.Board[rsIndex] = g.Board[reIndex]
		g.Board[reIndex] = piece.Empty
	}

}
