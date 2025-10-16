package engine

import (
	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/storage"
)

type CommandContext struct {
	storage *storage.KV
}

func NewCommandContext(storage *storage.KV) *CommandContext {
	return &CommandContext{storage: storage}
}
