package command

import (
	"fmt"
	"strings"

	"github.com/sampgo/sampgo"
)

type command struct {
	Handler interface{}
}

var commands = make(map[string]command)

// RegisterCmd allows you to register a command.
func RegisterCmd(name string, handler interface{}) error {
	_, ok := commands["/"+name]
	if ok {
		return fmt.Errorf("command name already exists")
	}

	commands["/"+name] = command{Handler: handler}
	return nil
}

// HandleCmds handles commands.
func HandleCmds(p sampgo.Player, cmd string) error {
	evt, ok := commands[cmd]
	if !ok {
		return fmt.Errorf("command hasnt got a handler")
	}

	fn, ok := evt.Handler.(func(sampgo.Player, string) error)
	if !ok {
		return fmt.Errorf("command is formatted incorrectly")
	}

	c := strings.Split(cmd, " ")
	params := []string{}

	for i, param := range c {
		if i == 0 {
			continue
		}

		params = append(params, param)
	}

	go fn(p, params)
	return nil
}
