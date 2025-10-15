package engine

import (
	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/storage"
)

type CommandContext struct {
	storage *storage.Storage
}

func NewCommandContext(storage *storage.Storage) *CommandContext {
	return &CommandContext{storage: storage}
}
