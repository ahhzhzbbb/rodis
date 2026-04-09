package command

import (
	"rodis/internal/engine"
	"rodis/internal/protocol/resp"
)

type CommandContext struct {
	kv *engine.KeyValue
}

func NewCommandContext(kv *engine.KeyValue) *CommandContext {
	return &CommandContext{kv: kv}
}

type Command interface {
	Execute(args []resp.Value, ctx *CommandContext) resp.Value
}
