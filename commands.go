package main

import (
	"strings"
)

func newErrorValue(s string) Value { return Value{typ: "error", str: s} }

func errWrongArgs() Value { return newErrorValue("ERR wrong number of arguments") }

func newStringValue(s string) Value { return Value{typ: "string", str: s} }

func newIntValue(n int) Value { return Value{typ: "integer", num: n} }

func newBulkValue(s string) Value { return Value{typ: "bulk", bulk: s} }

func newArrayValue(v []Value) Value { return Value{typ: "array", array: v} }

func newNullValue() Value { return Value{typ: "null"} }

func handlePing(args []string) Value {
	if len(args) == 0 {
		return newStringValue("PONG")
	}
	if len(args) == 1 {
		return newBulkValue(args[0])
	}
	return errWrongArgs()
}

func handleEcho(args []string) Value {
	if len(args) == 1 {
		return newBulkValue(args[0])
	}
	return errWrongArgs()
}

func handleSet(args []string, storage *Storage) Value {
	if len(args) == 2 {
		storage.Set(args[0], args[1])
		return newStringValue("OK")
	}
	return errWrongArgs()
}

func handleGet(args []string, storage *Storage) Value {
	if len(args) != 1 {
		return errWrongArgs()
	}
	val, isp := storage.Get(args[0])
	if !isp {
		return newNullValue()
	}
	return newBulkValue(val)
}

func handleDel(args []string, storage *Storage) Value {
	if len(args) == 0 {
		return errWrongArgs()
	}
	return newIntValue(storage.Delete(args...))
}

func handleType(args []string, storage *Storage) Value {
	if len(args) != 1 {
		return errWrongArgs()
	}
	return newBulkValue(storage.Type(args[0]))
}

func handleExists(args []string, storage *Storage) Value {
	if len(args) == 0 {
		return errWrongArgs()
	}
	return newIntValue(storage.Exists(args...))
}

func handleKeys(args []string, storage *Storage) Value {
	if len(args) != 1 {
		return errWrongArgs()
	}
	matches, err := storage.Keys(args[0])
	if err != nil {
		return newErrorValue("ERR invalid pattern")
	}
	values := make([]Value, len(matches))
	for i, m := range matches {
		values[i] = newBulkValue(m)
	}
	return newArrayValue(values)
}

func handleFlushdb(args []string, storage *Storage) Value {
	if len(args) != 0 {
		return errWrongArgs()
	}
	storage.Flushdb()
	return newStringValue("OK")
}

func handleCommand(cmd string, args []string, storage *Storage) Value {
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
	default:
		return Value{typ: "error", str: "ERR unknown command"}
	}
}
