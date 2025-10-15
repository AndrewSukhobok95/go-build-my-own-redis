package engine

import "github.com/AndrewSukhobok95/go-build-my-own-redis/internal/resp"

type Command struct {
	name    string
	arity   int
	handler func()
}

type CommandHandler func(ctx *CommandContext, args []resp.Value) resp.Value

var registry map[string]*Command

func RegisterCommand(name string, arity int, handler CommandHandler)

func GetCommand(name string) (*Command, bool)
