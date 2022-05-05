package engine

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Engine struct {
	Mode    string
	Running bool
}

func (self *Engine) Run() {
	self.Running = true
	reader := bufio.NewReader(os.Stdin)
	for self.Running {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("info string %s\n", err)
			continue
		}
		self.HandleCommand(strings.TrimSpace(input))
	}
}

func (self *Engine) HandleCommand(command string) bool {
	if len(command) == 0 {
		return false
	}
	whiteSpace := regexp.MustCompile("\\s+")
	parts := whiteSpace.Split(command, -1)
	switch parts[0] {
	case "uci":
		self.SendCommand("id name Random Engine\n")
		self.SendCommand("id author Caleb B\n")
		// send options
		self.SendCommand("uciok\n")
	case "debug":
		// TODO
	case "isready":
		self.SendCommand("readyok") // Not really, but oh well
	case "setoption":
		// TODO
	case "ucinewgame":
		// TODO
	case "position":
		// TODO
	case "go":
		// TODO
	case "stop":
		// TODO
	case "ponderhit":
		// TODO
	case "quit":
		self.Running = false
	}
	return true
}

func (self *Engine) SendCommand(command string) bool {
	if self.Running {
		fmt.Printf("%s\n", command)
		return true
	}
	return false
}
