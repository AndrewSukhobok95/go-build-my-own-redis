package commands

import (
	"errors"
	"log"
	"strconv"

	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/engine"
	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/resp"
	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/storage"
)

func handlePush(ctx *engine.CommandContext, args []string, isLeft bool) resp.Value {
	if len(args) < 2 {
		return errWrongArgs()
	}

	var (
		n   int
		err error
	)

	if isLeft {
		n, err = ctx.Storage().LPush(args[0], args[1:]...)
	} else {
		n, err = ctx.Storage().RPush(args[0], args[1:]...)
	}

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

func handleLPush(ctx *engine.CommandContext, args []string) resp.Value {
	return handlePush(ctx, args, true)
}

func handleRPush(ctx *engine.CommandContext, args []string) resp.Value {
	return handlePush(ctx, args, false)
}

func handlePop(ctx *engine.CommandContext, args []string, isLeft bool) resp.Value {
	if len(args) != 1 {
		return errWrongArgs()
	}

	var (
		popped string
		err    error
	)

	if isLeft {
		popped, err = ctx.Storage().LPop(args[0])
	} else {
		popped, err = ctx.Storage().RPop(args[0])
	}

	if err != nil {
		switch {
		case errors.Is(err, storage.ErrWrongType):
			return resp.NewErrorValue("WRONGTYPE Operation against a key holding the wrong kind of value")
		default:
			log.Printf("internal error in INCR: %v", err)
			return resp.NewErrorValue("ERR internal error")
		}
	}
	return resp.NewBulkValue(popped)
}

func handleLPop(ctx *engine.CommandContext, args []string) resp.Value {
	return handlePop(ctx, args, true)
}

func handleRPop(ctx *engine.CommandContext, args []string) resp.Value {
	return handlePop(ctx, args, false)
}

func handleLLen(ctx *engine.CommandContext, args []string) resp.Value {
	if len(args) != 1 {
		return errWrongArgs()
	}

	n, err := ctx.Storage().LLen(args[0])
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

func handleLRange(ctx *engine.CommandContext, args []string) resp.Value {
	if len(args) != 3 {
		return errWrongArgs()
	}

	start, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return resp.NewErrorValue("WRONGTYPE index value is not an integer or out of range")
	}

	stop, err := strconv.ParseInt(args[2], 10, 64)
	if err != nil {
		return resp.NewErrorValue("WRONGTYPE index value is not an integer or out of range")
	}

	r, err := ctx.Storage().LRange(args[0], int(start), int(stop))
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrWrongType):
			return resp.NewErrorValue("WRONGTYPE Operation against a key holding the wrong kind of value")
		default:
			log.Printf("internal error in INCR: %v", err)
			return resp.NewErrorValue("ERR internal error")
		}
	}

	bulks := make([]resp.Value, len(r))
	for i, b := range r {
		bulks[i] = resp.NewBulkValue(b)
	}
	return resp.NewArrayValue(bulks)
}

func init() {
	engine.RegisterCommand("LPUSH", -2, true, handleLPush)
	engine.RegisterCommand("RPUSH", -2, true, handleRPush)
	engine.RegisterCommand("LPOP", 1, true, handleLPop)
	engine.RegisterCommand("RPOP", 1, true, handleRPop)
	engine.RegisterCommand("LLEN", 1, false, handleLLen)
	engine.RegisterCommand("LRANGE", 3, false, handleLRange)
}
