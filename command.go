package command

import (
	"errors"
	"strings"

	"github.com/sampgo/sampgo"
)

// Command struct represents a command
type Command struct {
	Name   string
	Alias  []string
	Prefix string
}

// Context is a struct passed to a command handler
type Context struct {
	Player sampgo.Player
	Args   []string
}

// ErrorContext is a struct passed to the command error callback
type ErrorContext struct {
	Command Command
	Error   error
}

// Func is a custom type for command handler
type Func func(Context) error

// ErrorFunc is a custom type for command error handler
type ErrorFunc func(ErrorContext) bool

// BeforeFunc is a custom type for before func handler
type BeforeFunc Func

// AfterFunc is a custom type for after func handler
type AfterFunc Func

// internalCommand represents command and it's callbacks
type internalCommand struct {
	cmd Command
	fn  []Func
}

var (
	// ErrInvalidCommand is an exported error raised when Name field of Command is empty
	ErrInvalidCommand = errors.New("the required fields for a command was not filled")
	// ErrCantSync is an exported error raised when Sync is fired before Handle
	ErrCantSync = errors.New("a command can only be synced after handle is fired")
)

var (
	// commands stores command handlers
	commands        = make(map[string]internalCommand)
	commandPrefixes []string
)

var (
	// errorFunc is a function fired on command error
	errorFunc ErrorFunc
	// beforeFunc is a function fired before the command handler
	beforeFunc BeforeFunc
	// afterFunc is a function fired after the comman handler
	afterFunc AfterFunc
)

// pushCommandPrefix just pushes a new prefix into commandPrefixes
func pushCommandPrefix(prefix string) {
	for _, cmdPrefix := range commandPrefixes {
		if prefix == cmdPrefix {
			return
		}
	}
	commandPrefixes = append(commandPrefixes, prefix)
}

// NewCommand returns a new command
func NewCommand(command Command) (cmd *Command) {
	pushCommandPrefix("/")
	cmd = &command
	cmd.Prefix = "/"
	return
}

// New returns an empty Command struct
func New() (cmd *Command) {
	pushCommandPrefix("/")
	cmd.Prefix = "/"
	return
}

// SetName sets the command name
func (cmd *Command) SetName(name string) *Command {
	cmd.Name = name
	return cmd
}

// SetAlias sets the aliases for the command
func (cmd *Command) SetAlias(aliases ...string) *Command {
	cmd.Alias = aliases
	return cmd
}

// SetPrefix sets the command prefix
func (cmd *Command) SetPrefix(prefix string) *Command {
	pushCommandPrefix("/")
	cmd.Prefix = prefix
	return cmd
}

// Sync synces command detail changes
func (cmd *Command) Sync() (err error) {
	if cmd.Name == "" {
		err = ErrInvalidCommand
		return
	}
	command, ok := commands[cmd.Name]
	if !ok {
		err = ErrCantSync
		return
	}

	command.cmd = *cmd
	commands[cmd.Name] = command

	for _, alias := range command.cmd.Alias {
		commands[alias] = command
	}

	pushCommandPrefix(command.cmd.Prefix)

	return
}

// SetGeneralCommandErrorFunc sets the general command error handler
func SetGeneralCommandErrorFunc(fn ErrorFunc) {
	errorFunc = fn
}

// SetGeneralCommandBeforeFunc sets the general before func handler
func SetGeneralCommandBeforeFunc(fn BeforeFunc) {
	beforeFunc = fn
}

// SetGeneralCommandAfterFunc sets the general after func handler
func SetGeneralCommandAfterFunc(fn AfterFunc) {
	afterFunc = fn
}

// Handle will register the command, and fire the callback
func (cmd *Command) Handle(fn Func) (err error) {
	if cmd.Name == "" {
		err = ErrInvalidCommand
		return
	}

	if cmd.Prefix == "" {
		cmd.Prefix = "/"
	}

	pushCommandPrefix(cmd.Prefix)

	command, ok := commands[cmd.Name]
	if !ok {
		commands[cmd.Name] = internalCommand{}
	}

	command.cmd = *cmd
	command.fn = append(commands[cmd.Name].fn, fn)

	commands[cmd.Name] = command

	for _, alias := range command.cmd.Alias {
		commands[alias] = command
	}

	return
}

// handler is internal command handler
func handler(p sampgo.Player, text string) bool {
	if len(commands) == 0 {
		return false
	}

	var (
		command *internalCommand = nil
		cmdName string
		args    []string
	)

	for _, prefix := range commandPrefixes {
		cmdName = strings.TrimPrefix(text, prefix)
		args = strings.Split(cmdName, " ")

		if len(args) == 1 {
			args = nil
		} else {
			cmdName = args[0]
			args = append(args[1:])
		}

		cmd, ok := commands[cmdName]
		if !ok {
			continue
		} else {
			command = &cmd
			break
		}
	}

	if command == nil {
		return false
	}

	cmdCtx := Context{p, args}

	var err error

	if beforeFunc != nil {
		err = beforeFunc(cmdCtx)
		if err != nil {
			return false
		}
	}
	for _, fn := range (*command).fn {
		if err := fn(cmdCtx); err != nil {
			if errorFunc != nil {
				if !errorFunc(ErrorContext{command.cmd, err}) {
					return false
				}
			} else {
				return false
			}
		}
	}
	if afterFunc != nil {
		err = afterFunc(cmdCtx)
		if err != nil {
			return false
		}
	}

	return true
}

func init() {
	sampgo.On("playerCommandText", handler)
	sampgo.On("playerText", handler)
}
