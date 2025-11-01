package engine

import (
	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/resp"
	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/storage"
)

type CommandContext struct {
	storage       *storage.KV
	inTransaction bool
	queued        []func() resp.Value
}

func NewCommandContext(storage *storage.KV) *CommandContext {
	return &CommandContext{
		storage:       storage,
		inTransaction: false,
		queued:        make([]func() resp.Value, 0),
	}
}

func (c *CommandContext) Storage() *storage.KV {
	return c.storage
}

func (c *CommandContext) InTransaction() bool {
	return c.inTransaction
}

func (c *CommandContext) BeginTransaction() {
	c.inTransaction = true
	c.queued = c.queued[:0]
}

func (c *CommandContext) DiscardTransaction() {
	c.inTransaction = false
	c.queued = c.queued[:0]
}

func (c *CommandContext) EnqueueCommand(fn func() resp.Value) {
	c.queued = append(c.queued, fn)
}

func (c *CommandContext) ExecuteTransaction() []resp.Value {
	results := make([]resp.Value, 0, len(c.queued))
	for i, fn := range c.queued {
		results = append(results, fn())
		c.queued[i] = nil
	}
	c.queued = c.queued[:0]
	c.inTransaction = false
	return results
}
