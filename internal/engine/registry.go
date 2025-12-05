package engine

import (
	"fmt"
	"strings"

	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/resp"
)

type CommandHandler func(ctx *CommandContext, args []string) resp.Value

type Command struct {
	name    string
	arity   int
	isWrite bool
	handler CommandHandler
}

var registry map[string]*Command = make(map[string]*Command)

func AllCommands() map[string]*Command {
	return registry
}

func RegisterCommand(name string, arity int, isWrite bool, handler CommandHandler) {
	name = strings.ToUpper(name)
	if name == "" {
		panic("command name cannot be empty")
	}
	if handler == nil {
		panic(fmt.Sprintf("command %q has nil handler", name))
	}
	if _, exists := registry[name]; exists {
		panic(fmt.Sprintf("command %q already registered", name))
	}
	registry[name] = &Command{name: name, arity: arity, isWrite: isWrite, handler: handler}
}

func GetCommand(name string) (*Command, bool) {
	cmd, isp := registry[name]
	return cmd, isp
}
