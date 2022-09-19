package game

import (
	"fmt"
	"math"
	"regexp"
	"strconv"

	"bareman.net/chess-engine/game/move"
)

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
