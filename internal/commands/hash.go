package commands

import (
	"errors"
	"log"

	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/engine"
	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/resp"
	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/storage"
)

func handleHSet(ctx *engine.CommandContext, args []string) resp.Value {
	if len(args) != 3 {
		return errWrongArgs()
	}
	isNew, err := ctx.Storage().HSet(args[0], args[1], args[2])
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrWrongType):
			return resp.NewErrorValue("WRONGTYPE Operation against a key holding the wrong kind of value")
		default:
			log.Printf("internal error in HSET: %v", err)
			return resp.NewErrorValue("ERR internal error")
		}
	}
	return resp.NewIntValue(int64(isNew))
}

func handleHGet(ctx *engine.CommandContext, args []string) resp.Value {
	if len(args) != 2 {
		return errWrongArgs()
	}
	val, exists, err := ctx.Storage().HGet(args[0], args[1])
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrWrongType):
			return resp.NewErrorValue("WRONGTYPE Operation against a key holding the wrong kind of value")
		default:
			log.Printf("internal error in HGET: %v", err)
			return resp.NewErrorValue("ERR internal error")
		}
	}
	if !exists {
		return resp.NewNullValue()
	}
	return resp.NewBulkValue(val)
}

func handleHGetAll(ctx *engine.CommandContext, args []string) resp.Value {
	if len(args) != 1 {
		return errWrongArgs()
	}
	flattenedHash, err := ctx.Storage().HGetAll(args[0])
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrWrongType):
			return resp.NewErrorValue("WRONGTYPE Operation against a key holding the wrong kind of value")
		default:
			log.Printf("internal error in HGETALL: %v", err)
			return resp.NewErrorValue("ERR internal error")
		}
	}
	bulks := make([]resp.Value, len(flattenedHash))
	for i, v := range flattenedHash {
		bulks[i] = resp.NewBulkValue(v)
	}
	return resp.NewArrayValue(bulks)
}

func init() {
	engine.RegisterCommand("HSET", 3, true, handleHSet)
	engine.RegisterCommand("HGET", 2, false, handleHGet)
	engine.RegisterCommand("HGETALL", 1, false, handleHGetAll)
}
