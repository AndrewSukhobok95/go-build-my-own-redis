package engine

import "github.com/AndrewSukhobok95/go-build-my-own-redis/internal/resp"

func DispatchCommand(ctx *CommandContext, cmdName string, args []resp.Value) resp.Value
