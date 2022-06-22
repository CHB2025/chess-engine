package engine

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"bareman.net/chess-engine/game"
	"bareman.net/chess-engine/game/move"
)

type Engine struct {
	mu        sync.Mutex
	game      *game.Game
	isDebug   bool
	isRunning bool
}

func (e *Engine) Run() {
	e.isRunning = true
	reader := bufio.NewReader(os.Stdin)
	for e.isRunning {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("info string %s\n", err)
			continue
		}
		e.handleCommand(strings.TrimSpace(input)) //Can this cause race conditions?
	}
}

func (e *Engine) handleCommand(command string) bool {
	if len(command) == 0 {
		return false
	}
	e.sendCommand("info string " + command)
	whiteSpace := regexp.MustCompile(`\s+`)
	parts := whiteSpace.Split(strings.TrimSpace(command), -1)
	switch parts[0] {
	case "uci":
		e.sendCommand("id name Random Engine")
		e.sendCommand("id author Caleb B")
		// send options
		// Can handle opening and endgame books and more custom settings as well, but I think only Hash is really required.
		e.sendCommand("option name Hash type spin default 1 min 1 max 128") // Not sure what this should change about the engine. Have to read up about it.
		e.sendCommand("option name Ponder type check default true")         // Remove if engine doesn't support Pondering
		e.sendCommand("option name UCI_ShowCurrLine type check default false")

		e.sendCommand("uciok")
	case "debug":
		if len(parts) >= 2 && parts[1] == "on" {
			go e.setIsDebug(true)
		} else if len(parts) >= 2 && parts[1] == "off" {
			go e.setIsDebug(false)
		}
	case "isready":
		e.mu.Lock()
		defer e.mu.Unlock() //Ready when this process can lock the state?
		e.sendCommand("readyok")
	case "setoption":
		// TODO
	case "ucinewgame":
		// TODO
	case "position":
		go e.handlePosition(parts[1:])
	case "go":
		go e.handleGo(parts[1:])
	case "stop":
		go e.handleStop()
	case "ponderhit":
		// TODO
	case "board":
		e.mu.Lock()
		fmt.Println(e.game)
		e.mu.Unlock()
	case "fen":
		fmt.Println(e.game.ToFEN())
	case "undo":
		e.mu.Lock()
		e.game.Unmake()
		e.mu.Unlock()
	case "quit":
		e.mu.Lock()
		e.isRunning = false
		e.mu.Unlock()
	}
	return true
}

func (e *Engine) setIsDebug(val bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.isDebug = val
}

func (e *Engine) handlePosition(command []string) {
	if len(command) == 0 {
		return
	}
	e.mu.Lock()
	defer e.mu.Unlock()

	movesIndex := 0
	switch strings.ToLower(command[0]) {
	case "fen":
		g, err := game.FromFEN(strings.Join(command[1:7], " "))
		if err != nil {
			e.sendCommand("info string " + err.Error())
			return
		}
		e.game = g
		movesIndex = 7

	case "startpos":
		e.game = game.Default()
		movesIndex = 1
	}
	if len(command) > movesIndex {
		for _, mv := range command[movesIndex+1:] {
			err := e.game.Make(mv)
			if err != nil {
				e.sendCommand("info string " + err.Error())
				break
			}
		}
	}
	fmt.Print(e.game)
}

func (e *Engine) handleGo(options []string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	moveReg := regexp.MustCompile(move.MoveRegex)
	var moves []string
	var ponder, infinite bool
	var wtime, btime, winc, binc, movestogo, depth, nodes, mate, movetime int
	var opt string

	for len(options) > 0 {
		opt, options = options[0], options[1:]
		var err error
		switch strings.ToLower(opt) {
		case "searchmoves":
			for len(options) > 0 && moveReg.MatchString(options[0]) {
				moves = append(moves, options[0])
				options = options[1:]
			}
		case "ponder":
			ponder = true
		case "wtime":
			wtime, err = strconv.Atoi(options[0])
			options = options[1:]
		case "btime":
			btime, err = strconv.Atoi(options[0])
			options = options[1:]
		case "winc":
			winc, err = strconv.Atoi(options[0])
			options = options[1:]
		case "binc":
			binc, err = strconv.Atoi(options[0])
			options = options[1:]
		case "movestogo":
			movestogo, err = strconv.Atoi(options[0])
			options = options[1:]
		case "depth":
			depth, err = strconv.Atoi(options[0])
			options = options[1:]
		case "nodes":
			nodes, err = strconv.Atoi(options[0])
			options = options[1:]
		case "mate":
			mate, err = strconv.Atoi(options[0])
			options = options[1:]
		case "movetime":
			movetime, err = strconv.Atoi(options[0])
			options = options[1:]
		case "infinite":
			infinite = true
		}

		if err != nil {
			e.sendCommand("info string Invalid go command")
			break
		}
	}

	fmt.Sprintf(
		`Running search with the following parameters:
		Search moves: %s
		Ponder: %v
		White time: %v
		Black Time: %v
		White Increment: %v
		Black Increment: %v
		Moves til Time Change: %v
		Depth: %v
		Nodes: %v
		Mate: %v
		Move Time: %v
		Infinite: %v`,
		moves, ponder, wtime, btime, winc, binc, movestogo, depth, nodes, mate, movetime, infinite)
	// e.sendCommand(output)

	mvs := e.game.AllValidMoves()
	ind := rand.Intn(len(mvs))
	e.sendCommand(fmt.Sprintf("bestmove %v\n", mvs[ind]))
}

func (e *Engine) handleStop() {
	moves := e.game.AllValidMoves()
	ind := rand.Intn(len(moves))
	e.sendCommand(fmt.Sprintf("bestmove %v\n", moves[ind]))
}

func (e *Engine) sendCommand(command string) bool {
	if e.isRunning {
		fmt.Printf("%s\n", command)
		return true
	}
	return false
}
