package game

import "bareman.net/chess-engine/game/board"

func NewGame() Game {
	return Game{
		WhiteToMove:    true,
		Board:          board.DefaultBoard(),
		WhiteCanCastle: true,
		BlackCanCastle: true,
	}
}

type Game struct {
	WhiteToMove    bool
	Board          board.Board
	WhiteCanCastle bool
	BlackCanCastle bool
}
