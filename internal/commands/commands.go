package commands

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/resp"
	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/storage"
)

func errWrongArgs() resp.Value { return resp.NewErrorValue("ERR wrong number of arguments") }

func handlePing(args []string) resp.Value {
	if len(args) == 0 {
		return resp.NewStringValue("PONG")
	}
	if len(args) == 1 {
		return resp.NewBulkValue(args[0])
	}
	return errWrongArgs()
}

func handleEcho(args []string) resp.Value {
	if len(args) == 1 {
		return resp.NewBulkValue(args[0])
	}
	return errWrongArgs()
}

func handleSet(args []string, storage *storage.Storage) resp.Value {
	if len(args) == 2 {
		storage.Set(args[0], args[1])
		return resp.NewStringValue("OK")
	}
	return errWrongArgs()
}

func handleGet(args []string, storage *storage.Storage) resp.Value {
	if len(args) != 1 {
		return errWrongArgs()
	}
	val, isp := storage.Get(args[0])
	if !isp {
		return resp.NewNullValue()
	}
	return resp.NewBulkValue(val)
}

func handleDel(args []string, storage *storage.Storage) resp.Value {
	if len(args) == 0 {
		return errWrongArgs()
	}
	return resp.NewIntValue(storage.Delete(args...))
}

func handleType(args []string, storage *storage.Storage) resp.Value {
	if len(args) != 1 {
		return errWrongArgs()
	}
	return resp.NewBulkValue(storage.Type(args[0]))
}

func handleExists(args []string, storage *storage.Storage) resp.Value {
	if len(args) == 0 {
		return errWrongArgs()
	}
	return resp.NewIntValue(storage.Exists(args...))
}

func handleKeys(args []string, storage *storage.Storage) resp.Value {
	if len(args) != 1 {
		return errWrongArgs()
	}
	matches, err := storage.Keys(args[0])
	if err != nil {
		return resp.NewErrorValue("ERR invalid pattern")
	}
	values := make([]resp.Value, len(matches))
	for i, m := range matches {
		values[i] = resp.NewBulkValue(m)
	}
	return resp.NewArrayValue(values)
}

func handleFlushdb(args []string, storage *storage.Storage) resp.Value {
	if len(args) != 0 {
		return errWrongArgs()
	}
	storage.Flushdb()
	return resp.NewStringValue("OK")
}

func handleExpire(args []string, storage *storage.Storage, useSeconds bool) resp.Value {
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

func handleTTL(args []string, storage *storage.Storage, useSeconds bool) resp.Value {
	if len(args) != 1 {
		return errWrongArgs()
	}

	ttlMilli := min(storage.TTL(args[0]), math.MaxInt64)
	if ttlMilli < 0 {
		return resp.NewIntValue(int(ttlMilli))
	}
	if useSeconds {
		return resp.NewIntValue(int(ttlMilli / 1000))
	}
	return resp.NewIntValue(int(ttlMilli))
}

func HandleCommand(cmd string, args []string, storage *storage.Storage) resp.Value {
	switch strings.ToUpper(cmd) {
	case "PING":
		return handlePing(args)
	case "ECHO":
		return handleEcho(args)
	case "SET":
		return handleSet(args, storage)
	case "GET":
		return handleGet(args, storage)
	case "DEL":
		return handleDel(args, storage)
	case "TYPE":
		return handleType(args, storage)
	case "EXISTS":
		return handleExists(args, storage)
	case "KEYS":
		return handleKeys(args, storage)
	case "FLUSHDB":
		return handleFlushdb(args, storage)
	case "EXPIRE":
		return handleExpire(args, storage, true)
	case "PEXPIRE":
		return handleExpire(args, storage, false)
	case "TTL":
		return handleTTL(args, storage, true)
	case "PTTL":
		return handleTTL(args, storage, false)
	default:
		return resp.NewErrorValue("ERR unknown command")
	}
}
