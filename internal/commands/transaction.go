package commands

import (
	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/engine"
	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/resp"
)

func handleMulti(ctx *engine.CommandContext, args []string) resp.Value {
	if len(args) != 0 {
		return errWrongArgs()
	}
	ctx.BeginTransaction()
	return resp.NewStringValue("OK")
}

func handleExec(ctx *engine.CommandContext, args []string) resp.Value {
	if len(args) != 0 {
		return errWrongArgs()
	}
	if !ctx.InTransaction() {
		return resp.NewErrorValue("ERR EXEC without MULTI")
	}
	results := ctx.ExecuteTransaction()
	return resp.NewArrayValue(results)
}

func handleDiscard(ctx *engine.CommandContext, args []string) resp.Value {
	if len(args) != 0 {
		return errWrongArgs()
	}
	if !ctx.InTransaction() {
		return resp.NewErrorValue("ERR DISCARD without MULTI")
	}
	ctx.DiscardTransaction()
	return resp.NewStringValue("OK")
}

func init() {
	engine.RegisterCommand("MULTI", 0, handleMulti)
	engine.RegisterCommand("EXEC", 0, handleExec)
	engine.RegisterCommand("DISCARD", 0, handleDiscard)
}
