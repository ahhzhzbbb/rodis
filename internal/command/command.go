package command

import (
	"rodis/internal/engine"
	"rodis/internal/protocol/resp"
)

type CommandContext struct {
	k *engine.KeyValue
}

func NewCommandContext(k *engine.KeyValue) *CommandContext {
	return &CommandContext{
		k: k,
	}
}

type Command interface {
	Execute(args []resp.Value, ctx *CommandContext) resp.Value
}
