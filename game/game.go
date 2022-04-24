package game

import (
	"fmt"

	"bareman.net/chess-engine/game/board"
)

func NewGame() Game {
	b := board.DefaultBoard()
	var whiteKing, blackKing *board.Piece
	for _, row := range b {
		for _, piece := range row {
			if piece != nil && piece.Symbol == 'k' {
				if piece.IsWhite {
					whiteKing = piece
				} else {
					blackKing = piece
				}
			}
		}
	}
	return Game{
		WhiteToMove:    true,
		Board:          b,
		WhiteCanCastle: true,
		BlackCanCastle: true,
		WhiteKing:      whiteKing,
		BlackKing:      blackKing,
	}
}

type Game struct {
	WhiteToMove    bool
	Board          board.Board
	WhiteCanCastle bool
	BlackCanCastle bool
	WhiteKing      *board.Piece
	BlackKing      *board.Piece
}

func (g *Game) Move(start, finish board.Position) error {
	piece := g.Board.GetPiece(start)
	if piece == nil {
		return fmt.Errorf("No piece at chosen position")
	}
	if piece.IsWhite != g.WhiteToMove {
		return fmt.Errorf("Wrong color piece chosen")
	}

	newBoard, err := g.Board.Move(start, finish)
	if err != nil {
		return fmt.Errorf("Error moving piece: %v", err.Error())
	}
	king := g.WhiteKing
	if !g.WhiteToMove {
		king = g.BlackKing
	}
	if king.IsUnderAttack(newBoard) && king.IsUnderAttack(g.Board) {
		return fmt.Errorf("Illegal move. King is in check")
	} else if king.IsUnderAttack(newBoard) {
		return fmt.Errorf("Illegal move. Putting king in check")
	}
	g.Board = newBoard
	g.WhiteToMove = !g.WhiteToMove
	return nil
}

func (g *Game) Run() {

	for {
		fmt.Print(g.Board)
		color := "White"
		if !g.WhiteToMove {
			color = "Black"
		}
		fmt.Printf("%v to move\n", color)

		var move string
		fmt.Print("Enter a move: ")
		_, err := fmt.Scan(&move)
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err.Error())
			continue
		}
		if len(move) != 4 {
			fmt.Printf("Improperly formatted move.\n")
			continue
		}

		start, pErr := board.PositionFromString(move[:2])
		finish, pErr := board.PositionFromString(move[2:])
		if pErr != nil {
			fmt.Printf("Improperly formatted move.\n")
			continue
		}
		gErr := g.Move(start, finish)
		if gErr != nil {
			fmt.Printf("%v\n", gErr.Error())
		}
	}

}
