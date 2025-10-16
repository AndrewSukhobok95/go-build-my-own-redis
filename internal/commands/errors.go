package commands

import (
	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/resp"
)

func errWrongArgs() resp.Value { return resp.NewErrorValue("ERR wrong number of arguments") }
