package commands

import (
	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/engine"
	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/resp"
)

func handleKeys(ctx *engine.CommandContext, args []string) resp.Value {
	if len(args) != 1 {
		return errWrongArgs()
	}
	matches, err := ctx.Storage().Keys(args[0])
	if err != nil {
		return resp.NewErrorValue("ERR invalid pattern")
	}
	values := make([]resp.Value, len(matches))
	for i, m := range matches {
		values[i] = resp.NewBulkValue(m)
	}
	return resp.NewArrayValue(values)
}

func handleFlushdb(ctx *engine.CommandContext, args []string) resp.Value {
	if len(args) != 0 {
		return errWrongArgs()
	}
	ctx.Storage().Flushdb()
	return resp.NewStringValue("OK")
}

func init() {
	engine.RegisterCommand("KEYS", 1, handleKeys)
	engine.RegisterCommand("FLUSHDB", 0, handleFlushdb)
}
