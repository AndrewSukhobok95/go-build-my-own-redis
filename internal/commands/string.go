package commands

import (
	"errors"
	"log"
	"strconv"

	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/engine"
	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/resp"
	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/storage"
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

	val, isp, err := ctx.Storage().Get(args[0])
	if err != nil {
		if errors.Is(err, storage.ErrWrongType) {
			return resp.NewErrorValue("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		log.Printf("internal error in GET: %v", err)
		return resp.NewErrorValue("ERR internal error")
	}

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

func doIncr(ctx *engine.CommandContext, key string, delta int64) resp.Value {
	newValue, err := ctx.Storage().Incr(key, delta)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrWrongType):
			return resp.NewErrorValue("WRONGTYPE Operation against a key holding the wrong kind of value")
		case errors.Is(err, storage.ErrNotInteger), errors.Is(err, storage.ErrOverflow):
			return resp.NewErrorValue("ERR value is not an integer or out of range")
		default:
			log.Printf("internal error in INCR: %v", err)
			return resp.NewErrorValue("ERR internal error")
		}
	}
	return resp.NewIntValue(newValue)
}

func handleIncr(ctx *engine.CommandContext, args []string) resp.Value {
	if len(args) != 1 {
		return errWrongArgs()
	}
	return doIncr(ctx, args[0], 1)
}

func handleDecr(ctx *engine.CommandContext, args []string) resp.Value {
	if len(args) != 1 {
		return errWrongArgs()
	}
	return doIncr(ctx, args[0], -1)
}

func handleIncrBy(ctx *engine.CommandContext, args []string) resp.Value {
	if len(args) != 2 {
		return errWrongArgs()
	}
	incrementInt, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return resp.NewErrorValue("ERR value is not an integer or out of range")
	}
	return doIncr(ctx, args[0], incrementInt)
}

func handleDecrBy(ctx *engine.CommandContext, args []string) resp.Value {
	if len(args) != 2 {
		return errWrongArgs()
	}
	incrementInt, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return resp.NewErrorValue("ERR value is not an integer or out of range")
	}
	return doIncr(ctx, args[0], -incrementInt)
}

func handleAppend(ctx *engine.CommandContext, args []string) resp.Value {
	if len(args) != 2 {
		return errWrongArgs()
	}
	n, err := ctx.Storage().Append(args[0], args[1])
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrWrongType):
			return resp.NewErrorValue("WRONGTYPE Operation against a key holding the wrong kind of value")
		default:
			log.Printf("internal error in INCR: %v", err)
			return resp.NewErrorValue("ERR internal error")
		}
	}
	return resp.NewIntValue(int64(n))
}

func init() {
	engine.RegisterCommand("PING", 0, false, handlePing)
	engine.RegisterCommand("ECHO", 1, false, handleEcho)
	engine.RegisterCommand("SET", 2, true, handleSet)
	engine.RegisterCommand("GET", 1, false, handleGet)
	engine.RegisterCommand("DEL", -1, true, handleDel)
	engine.RegisterCommand("TYPE", 1, false, handleType)
	engine.RegisterCommand("EXISTS", -1, false, handleExists)
	engine.RegisterCommand("INCR", 1, true, handleIncr)
	engine.RegisterCommand("DECR", 1, true, handleDecr)
	engine.RegisterCommand("INCRBY", 2, true, handleIncrBy)
	engine.RegisterCommand("DECRBY", 2, true, handleDecrBy)
	engine.RegisterCommand("APPEND", 2, true, handleAppend)
}
