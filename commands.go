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

func handleCommand(cmd string, args []string) Value {
	switch strings.ToUpper(cmd) {
	case "PING":
		return handlePing(args)
	case "ECHO":
		return handleEcho(args)
	default:
		return Value{typ: "error", str: "ERR unknown command"}
	}
}
