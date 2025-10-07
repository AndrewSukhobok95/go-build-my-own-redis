package main

import "strings"

func handlePing(args []string) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}
	if len(args) == 1 {
		return Value{typ: "bulk", bulk: args[0]}
	}
	return Value{typ: "error", str: "ERR wrong number of arguments"}
}

func handleEcho(args []string) Value {
	if len(args) == 1 {
		return Value{typ: "bulk", bulk: args[0]}
	}
	return Value{typ: "error", str: "ERR wrong number of arguments, only one is expected"}
}

func handleSet(args []string, storage *Storage) Value {
	if len(args) == 2 {
		storage.Set(args[0], args[1])
		return Value{typ: "string", str: "OK"}
	}
	return Value{typ: "error", str: "ERR wrong number of arguments, two are expected"}
}

func handleGet(args []string, storage *Storage) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments, one is expected"}
	}
	val, isp := storage.Get(args[0])
	if !isp {
		return Value{typ: "null"}
	}
	return Value{typ: "bulk", bulk: val}
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
	default:
		return Value{typ: "error", str: "ERR unknown command"}
	}
}
