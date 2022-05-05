package main

import "bareman.net/chess-engine/engine"

func main() {
	//game1, _ := game.GameFromFEN("rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2")
	//game2, _ := game.GameFromFEN("8/5k2/3p4/1p1Pp2p/pP2Pp1P/P4P1K/8/8 b - - 99 50")

	// fmt.Printf("%v\n", game1)
	// for {
	// 	var position string
	// 	fmt.Print("Enter a position to inspect: ")
	// 	fmt.Scan(&position)
	// 	fmt.Printf("%v\n", position)
	// 	index := game.IndexFromPosition(position)
	// 	moves := game1.ValidMoves(index)
	// 	allowedPositions := []string{}
	// 	for _, mv := range moves {
	// 		allowedPositions = append(allowedPositions, game.PositionFromIndex(index+mv))
	// 	}
	// 	fmt.Printf("%v at %v can move to %v\n", game1.GetPiece(position), position, allowedPositions)
	// }
	engine := &engine.Engine{}
	engine.Run()
}
