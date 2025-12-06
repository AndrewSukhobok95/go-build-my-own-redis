package commands

import (
	"math"
	"strconv"
	"time"

	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/engine"
	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/resp"
	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/storage"
)

func handleExpire(args []string, storage *storage.KV, useSeconds bool) resp.Value {
	if len(args) != 2 {
		return errWrongArgs()
	}
	dur, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return resp.NewErrorValue("ERR value is not an integer or out of range")
	}
	var durTime time.Duration
	if useSeconds {
		durTime = time.Duration(dur) * time.Second
	} else {
		durTime = time.Duration(dur) * time.Millisecond
	}
	keyExists := storage.SetExpire(args[0], durTime)
	if keyExists {
		return resp.NewIntValue(1)
	}
	return resp.NewIntValue(0)
}

func handleTTL(args []string, storage *storage.KV, useSeconds bool) resp.Value {
	if len(args) != 1 {
		return errWrongArgs()
	}

	ttlMilli := min(storage.TTL(args[0]), math.MaxInt64)
	if ttlMilli < 0 {
		return resp.NewIntValue(int64(ttlMilli))
	}
	if useSeconds {
		return resp.NewIntValue(int64(ttlMilli / 1000))
	}
	return resp.NewIntValue(int64(ttlMilli))
}

func wrapHandleExpire(ctx *engine.CommandContext, args []string) resp.Value {
	return handleExpire(args, ctx.Storage(), true)
}

func wrapHandlePExpire(ctx *engine.CommandContext, args []string) resp.Value {
	return handleExpire(args, ctx.Storage(), false)
}

func wrapHandleTTL(ctx *engine.CommandContext, args []string) resp.Value {
	return handleTTL(args, ctx.Storage(), true)
}

func wrapHandlePTTL(ctx *engine.CommandContext, args []string) resp.Value {
	return handleTTL(args, ctx.Storage(), false)
}

func init() {
	engine.RegisterCommand("EXPIRE", 2, true, wrapHandleExpire)
	engine.RegisterCommand("PEXPIRE", 2, true, wrapHandlePExpire)
	engine.RegisterCommand("TTL", 1, false, wrapHandleTTL)
	engine.RegisterCommand("PTTL", 1, false, wrapHandlePTTL)
}
