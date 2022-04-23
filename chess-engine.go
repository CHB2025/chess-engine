package main

import (
	"fmt"

	game "bareman.net/chess-engine/game/board"
)

func main() {
	board, _ := game.BoardFromFEN("r1b1k1nr/p2p1pNp/n2B4/1p1NP2P/6P1/3P1Q2/P1P1K3/q5b1")

	fmt.Printf("%v", board)

}
