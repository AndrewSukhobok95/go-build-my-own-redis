package commands

import (
	"errors"
	"log"

	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/engine"
	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/resp"
	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/storage"
)

func handleSAdd(ctx *engine.CommandContext, args []string) resp.Value {
	if len(args) < 2 {
		return errWrongArgs()
	}
	added, err := ctx.Storage().SAdd(args[0], args[1:]...)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrWrongType):
			return resp.NewErrorValue("WRONGTYPE Operation against a key holding the wrong kind of value")
		default:
			log.Printf("internal error in SADD: %v", err)
			return resp.NewErrorValue("ERR internal error")
		}
	}
	return resp.NewIntValue(int64(added))
}

func handleSRem(ctx *engine.CommandContext, args []string) resp.Value {
	if len(args) < 2 {
		return errWrongArgs()
	}
	removed, err := ctx.Storage().SRem(args[0], args[1:]...)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrWrongType):
			return resp.NewErrorValue("WRONGTYPE Operation against a key holding the wrong kind of value")
		default:
			log.Printf("internal error in SREM: %v", err)
			return resp.NewErrorValue("ERR internal error")
		}
	}
	return resp.NewIntValue(int64(removed))
}

func handleSMembers(ctx *engine.CommandContext, args []string) resp.Value {
	if len(args) != 1 {
		return errWrongArgs()
	}

	members, err := ctx.Storage().SMembers(args[0])
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrWrongType):
			return resp.NewErrorValue("WRONGTYPE Operation against a key holding the wrong kind of value")
		default:
			log.Printf("internal error in SMEMBERS: %v", err)
			return resp.NewErrorValue("ERR internal error")
		}
	}

	bulks := make([]resp.Value, len(members))
	for i, m := range members {
		bulks[i] = resp.NewBulkValue(m)
	}
	return resp.NewArrayValue(bulks)
}

func handleSIsMember(ctx *engine.CommandContext, args []string) resp.Value {
	if len(args) != 2 {
		return errWrongArgs()
	}

	isMember, err := ctx.Storage().SIsMember(args[0], args[1])
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrWrongType):
			return resp.NewErrorValue("WRONGTYPE Operation against a key holding the wrong kind of value")
		default:
			log.Printf("internal error in SISMEMBER: %v", err)
			return resp.NewErrorValue("ERR internal error")
		}
	}

	if isMember {
		return resp.NewIntValue(1)
	}
	return resp.NewIntValue(0)
}

func init() {
	engine.RegisterCommand("SADD", -2, true, handleSAdd)
	engine.RegisterCommand("SREM", -2, true, handleSRem)
	engine.RegisterCommand("SMEMBERS", 1, false, handleSMembers)
	engine.RegisterCommand("SISMEMBER", 2, false, handleSIsMember)
}
