package game

import (
	"fmt"
	"math"
	"regexp"
	"strconv"

	"bareman.net/chess-engine/game/move"
)

const (
	colMask = 0b00000111
)

func coordinates(index int) (int, int) {
	col := index & colMask
	row := index >> 3
	return row, col
}

func indexFromPosition(pos string) int {
	reg := regexp.MustCompile(move.PositionRegex)
	if !reg.MatchString(pos) {
		return -1
	}
	row, _ := strconv.Atoi(string(pos[1])) // Guaranteed by regex
	col := int(pos[0] - 'a')
	return (row-1)<<3 + col
}

func positionFromIndex(index int) string {
	if index >= 64 || index < 0 {
		return "-"
	}
	row, col := coordinates(index)
	return fmt.Sprintf("%v%v", string('a'+rune(col)), row+1)
}

func mdistance(start, finish int) int {
	sRow, sCol := coordinates(start)
	fRow, fCol := coordinates(finish)
	y := math.Abs(float64(sCol - fCol))
	x := math.Abs(float64(sRow - fRow))
	return int(x + y)
}
