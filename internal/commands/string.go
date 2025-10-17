package commands

import (
	"strconv"

	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/engine"
	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/resp"
)

func handlePing(ctx *engine.CommandContext, args []string) resp.Value {
	if len(args) == 0 {
		return resp.NewStringValue("PONG")
	}
	if len(args) == 1 {
		return resp.NewBulkValue(args[0])
	}
	return errWrongArgs()
}

func handleEcho(ctx *engine.CommandContext, args []string) resp.Value {
	if len(args) == 1 {
		return resp.NewBulkValue(args[0])
	}
	return errWrongArgs()
}

func handleSet(ctx *engine.CommandContext, args []string) resp.Value {
	if len(args) == 2 {
		ctx.Storage().Set(args[0], args[1])
		return resp.NewStringValue("OK")
	}
	return errWrongArgs()
}

func handleGet(ctx *engine.CommandContext, args []string) resp.Value {
	if len(args) != 1 {
		return errWrongArgs()
	}
	val, isp := ctx.Storage().Get(args[0])
	if !isp {
		return resp.NewNullValue()
	}
	return resp.NewBulkValue(val)
}

func handleDel(ctx *engine.CommandContext, args []string) resp.Value {
	if len(args) == 0 {
		return errWrongArgs()
	}
	return resp.NewIntValue(int64(ctx.Storage().Delete(args...)))
}

func handleType(ctx *engine.CommandContext, args []string) resp.Value {
	if len(args) != 1 {
		return errWrongArgs()
	}
	return resp.NewBulkValue(ctx.Storage().Type(args[0]))
}

func handleExists(ctx *engine.CommandContext, args []string) resp.Value {
	if len(args) == 0 {
		return errWrongArgs()
	}
	return resp.NewIntValue(int64(ctx.Storage().Exists(args...)))
}

func handleIncr(ctx *engine.CommandContext, args []string) resp.Value {
	if len(args) != 1 {
		return errWrongArgs()
	}
	newValue, err := ctx.Storage().Incr(args[0], 1)
	if err != nil {
		return resp.NewErrorValue("ERR " + err.Error())
	}
	return resp.NewIntValue(newValue)
}

func handleDecr(ctx *engine.CommandContext, args []string) resp.Value {
	if len(args) != 1 {
		return errWrongArgs()
	}
	newValue, err := ctx.Storage().Incr(args[0], -1)
	if err != nil {
		return resp.NewErrorValue("ERR " + err.Error())
	}
	return resp.NewIntValue(newValue)
}

func handleIncrBy(ctx *engine.CommandContext, args []string) resp.Value {
	if len(args) != 2 {
		return errWrongArgs()
	}
	incrementInt, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return resp.NewErrorValue("ERR value is not an integer or out of range")
	}
	var newValue int64
	newValue, err = ctx.Storage().Incr(args[0], incrementInt)
	if err != nil {
		return resp.NewErrorValue("ERR " + err.Error())
	}
	return resp.NewIntValue(newValue)
}

func handleDecrBy(ctx *engine.CommandContext, args []string) resp.Value {
	if len(args) != 2 {
		return errWrongArgs()
	}
	incrementInt, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return resp.NewErrorValue("ERR value is not an integer or out of range")
	}
	var newValue int64
	newValue, err = ctx.Storage().Incr(args[0], -incrementInt)
	if err != nil {
		return resp.NewErrorValue("ERR " + err.Error())
	}
	return resp.NewIntValue(newValue)
}

func handleAppend(ctx *engine.CommandContext, args []string) resp.Value {
	if len(args) != 2 {
		return errWrongArgs()
	}
	return resp.NewIntValue(int64(ctx.Storage().Append(args[0], args[1])))
}

func init() {
	engine.RegisterCommand("PING", 0, handlePing)
	engine.RegisterCommand("ECHO", 1, handleEcho)
	engine.RegisterCommand("SET", 2, handleSet)
	engine.RegisterCommand("GET", 1, handleGet)
	engine.RegisterCommand("DEL", -1, handleDel)
	engine.RegisterCommand("TYPE", 1, handleType)
	engine.RegisterCommand("EXISTS", -1, handleExists)
	engine.RegisterCommand("INCR", 1, handleIncr)
	engine.RegisterCommand("DECR", 1, handleDecr)
	engine.RegisterCommand("INCRBY", 2, handleIncrBy)
	engine.RegisterCommand("DECRBY", 2, handleDecrBy)
	engine.RegisterCommand("APPEND", 2, handleAppend)
}
