package engine

import (
	"fmt"
	"strings"

	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/resp"
)

func DispatchCommand(ctx *CommandContext, cmdName string, args []string) resp.Value {
	cmdName = strings.ToUpper(cmdName)

	cmd, isp := GetCommand(cmdName)
	if !isp {
		return resp.NewErrorValue("ERR command not found")
	}

	if cmd.arity > 0 && len(args) != cmd.arity {
		return resp.NewErrorValue(fmt.Sprintf("ERR wrong number of arguments for '%s' command", cmdName))
	}
	if cmd.arity < 0 && len(args) < -cmd.arity {
		return resp.NewErrorValue(fmt.Sprintf("ERR wrong number of arguments for '%s' command", cmdName))
	}
	return cmd.handler(ctx, args)
}
