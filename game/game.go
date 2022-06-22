package game

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode"

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

func (g *Game) Make(mv string) error {
	move, err := move.EmptyMove(mv)
	if err != nil {
		return err
	}

	validMoves := g.ValidMoves(move.Origin)
	isValid := false
	for _, m := range validMoves {
		if m == mv {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("Invalid move given. Received %v\n", mv)
	}
	fmt.Printf("Making move %v\n", move)
	g.make(move)

	return nil
}

func (g *Game) make(mv *move.Move) {
	//fmt.Printf("Moving from %v (%v) to %v (%v)\n", mv.Origin, mv.OriginIndex(), mv.Dest, mv.DestIndex())
	p := g.Piece(mv.Origin)
	capture := g.Piece(mv.Dest)
	castle := p%piece.White == piece.King && mv.OriginIndex()/8 == mv.DestIndex()/8 && mdistance(mv.OriginIndex(), mv.DestIndex()) > 1 // piece is king, moving 2 or 3 spaces on one rank
	ep := p%piece.White == piece.Pawn && g.EPTarget == mv.DestIndex() && mdistance(mv.OriginIndex(), mv.DestIndex()) == 2              //piece is pawn, moving to target square diagonally

	mv.Capture, mv.Castle, mv.EnPassant, mv.EPTarget = capture, castle, ep, g.EPTarget

	if p%piece.White == piece.Pawn && mv.Origin[0] == mv.Dest[0] && mdistance(mv.OriginIndex(), mv.DestIndex()) == 2 {
		g.EPTarget = mv.DestIndex()/2 + mv.OriginIndex()/2
	} else {
		g.EPTarget = -1
	}

	g.Board[mv.DestIndex()] = g.Board[mv.OriginIndex()]
	g.Board[mv.OriginIndex()] = piece.Empty
	g.WhiteToMove = !g.WhiteToMove

	if ep {
		file := mv.Dest[0]
		rank := mv.Origin[1]
		capIndex := indexFromPosition(string(file) + string(rank))
		mv.Capture = g.Board[capIndex]
		g.Board[capIndex] = piece.Empty
	}
	if castle {
		fmt.Printf("Castle")
		rank := mv.Origin[1]
		file := (mv.Dest[0]-mv.Origin[0])/2 + mv.Origin[0]
		rePos := string(file) + string(rank)
		rsPos := "a" + string(rank)
		if file == 'g' {
			rsPos = "h" + string(rank)
		}
		reIndex := indexFromPosition(rePos)
		rsIndex := indexFromPosition(rsPos)
		g.Board[reIndex] = g.Board[rsIndex]
		g.Board[rsIndex] = piece.Empty
	}
	g.Moves = append(g.Moves, mv)
}

func (g *Game) Unmake() {
	move := g.Moves[len(g.Moves)-1]
	g.Moves = g.Moves[:len(g.Moves)-1]

	g.Board[move.OriginIndex()] = g.Board[move.DestIndex()]
	g.Board[move.DestIndex()] = move.Capture
	g.WhiteToMove = !g.WhiteToMove
	g.EPTarget = move.EPTarget
	if move.EnPassant {
		g.Board[move.DestIndex()] = piece.Empty
		file := move.Dest[0]
		rank := move.Origin[1]
		index := indexFromPosition(string(file) + string(rank))
		g.Board[index] = move.Capture
	}
	if move.Castle {

	}

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
	var moves []string
	for index, pi := range g.Board {
		if pi != piece.Empty && pi/piece.White != p/piece.White {
			moves = append(moves, g.moves(positionFromIndex(index))...)
		}
	}
	var attackers []string
	for _, mv := range moves {
		if mv[2:] == position {
			attackers = append(attackers, mv[:2])
		}
	}
	return attackers
}

func (g *Game) AllValidMoves() []string {
	var moves []string
	var king string
	color := piece.White
	if !g.WhiteToMove {
		color = piece.Black
	}
	for i, p := range g.Board {
		if p == piece.Piece(color|piece.King) {
			king = positionFromIndex(i)
			break
		}
	}
	for i, p := range g.Board {
		if p == piece.Empty || p.IsWhite() != g.WhiteToMove {
			continue
		}
		mvs := g.moves(positionFromIndex(i))
		for _, mv := range mvs {
			m, _ := move.EmptyMove(mv)
			g.make(m)
			var atks []string
			if p == piece.Piece(color|piece.King) {
				atks = g.attackers(m.Dest)
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

func (g *Game) ValidMoves(pos string) []string {
	p := g.Piece(pos)
	if p == piece.Empty || p.IsWhite() != g.WhiteToMove {
		return nil
	}
	var moves []string
	var king string
	color := piece.White
	if !g.WhiteToMove {
		color = piece.Black
	}
	for i, p := range g.Board {
		if p == piece.Piece(color|piece.King) {
			king = positionFromIndex(i)
			break
		}
	}
	mvs := g.moves(pos)
	for _, mv := range mvs {
		m, _ := move.EmptyMove(mv)
		g.make(m)
		var atks []string
		if p == piece.Piece(color|piece.King) {
			atks = g.attackers(m.Dest)
		} else {
			atks = g.attackers(king)
		}
		if len(atks) == 0 {
			moves = append(moves, mv)
		}
		g.Unmake()
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
		dir := -1 * (int(p/8)*2 - 3)
		// fmt.Printf("dir: %v\n", dir)
		index := start + Forward*dir
		if index < 64 && index >= 0 && g.Board[index] == piece.Empty {
			moves = append(moves, pos+positionFromIndex(index))

			if start/8-dir == 7 || (start)/8-dir == 0 && g.Board[start+2*Forward*dir] == piece.Empty { // in starting row
				moves = append(moves, pos+positionFromIndex(index+Forward*dir))
			}
		}
		if g.EPTarget != -1 && index+Left == g.EPTarget && mdistance(index, index+Left) == 1 {
			moves = append(moves, pos+positionFromIndex(index+Left))
		}
		if g.EPTarget != -1 && index+Right == g.EPTarget && mdistance(index, index+Right) == 1 {
			moves = append(moves, pos+positionFromIndex(index+Right))
		}
		leftAttack := start + Forward*dir + Left
		if leftAttack < 64 && leftAttack >= 0 &&
			mdistance(start, leftAttack) == 2 &&
			g.Board[leftAttack] != piece.Empty &&
			g.Board[leftAttack].IsWhite() != p.IsWhite() {
			moves = append(moves, pos+positionFromIndex(leftAttack))
		}
		rightAttack := start + Forward*dir + Right
		if rightAttack < 64 && rightAttack >= 0 &&
			mdistance(start, rightAttack) == 2 &&
			g.Board[rightAttack] != piece.Empty &&
			g.Board[rightAttack].IsWhite() != p.IsWhite() {
			moves = append(moves, pos+positionFromIndex(rightAttack))
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
			Backward + Backward + Left,
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

func Default() *Game {
	game, _ := FromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	return game
}

func FromFEN(fen string) (*Game, error) {
	if !isValidFEN(fen) {
		return nil, fmt.Errorf("Invalid FEN string. Received %v", fen)
	}
	sections := strings.Split(fen, " ")
	var board [64]piece.Piece
	for y, row := range strings.Split(sections[0], "/") {
		var offset int
		for x, symbol := range row {
			if unicode.IsNumber(symbol) {
				num, _ := strconv.Atoi(string(symbol))
				offset += num - 1
				continue
			}

			board[8*(7-y)+x+offset] = piece.FromRune(symbol)
		}
	}

	move, _ := strconv.Atoi(sections[5])
	halfMove, _ := strconv.Atoi(sections[4])

	game := &Game{
		Board:       board,
		MoveCount:   move,
		HalfMove:    halfMove,
		WhiteToMove: sections[1] == "w",
		WKCastle:    strings.Contains(sections[2], "K"),
		WQCastle:    strings.Contains(sections[2], "Q"),
		BKCastle:    strings.Contains(sections[2], "k"),
		BQCastle:    strings.Contains(sections[2], "q"),
		EPTarget:    indexFromPosition(sections[3]),
	}
	return game, nil
}

func isValidFEN(fen string) bool {
	sections := strings.Split(fen, " ")
	if len(sections) != 6 {
		return false
	}

	rows := strings.Split(sections[0], "/")
	if len(rows) != 8 {
		return false
	}

	posCount := 0
	for _, row := range rows {
		for _, symbol := range row {
			if unicode.IsNumber(symbol) {
				num, err := strconv.Atoi(string(symbol))
				if err != nil {
					return false
				}
				posCount += num
				continue
			}
			// Should check if it is one of rnbqkp
			posCount++
		}
	}
	return posCount == 64
	//Should check the rest of the sections too
}

func indexFromPosition(pos string) int {
	reg := regexp.MustCompile(move.PositionRegex)
	if !reg.MatchString(pos) {
		return -1
	}
	row, _ := strconv.Atoi(string(pos[1])) // Guaranteed by regex
	col := rune(pos[0])
	return (row-1)*8 + int(col-'a')
}

func positionFromIndex(index int) string {
	if index >= 64 || index < 0 {
		return "-"
	}
	col := index % 8
	row := index / 8
	return fmt.Sprintf("%v%v", string('a'+rune(col)), row+1)
}

func mdistance(start, finish int) int {
	y := math.Abs(float64(start/8 - finish/8))
	x := math.Abs(float64(start%8 - finish%8))
	return int(x + y)
}
