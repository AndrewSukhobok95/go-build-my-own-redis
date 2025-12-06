package engine

import (
	"fmt"
	"log"
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

	if !ctx.InReplay() && cmd.isWrite {
		log.Println("Appened to AOF")
		if err := ctx.aof.Append(cmdName, args); err != nil {
			log.Println("AOF append failed:", err)
		}
	}

	if ctx.InTransaction() && cmdName != "MULTI" && cmdName != "EXEC" && cmdName != "DISCARD" {
		ctx.EnqueueCommand(func() resp.Value {
			return cmd.handler(ctx, args)
		})
		return resp.NewStringValue("QUEUED")
	}
	return cmd.handler(ctx, args)
}
