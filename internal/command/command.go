package command

import (
	"rodis/internal/engine"
	"rodis/internal/protocol/resp"
)

type CommandContext struct {
	kv *engine.KeyValue
	et *engine.ExpireTime
}

func NewCommandContext(kv *engine.KeyValue, et *engine.ExpireTime) *CommandContext {
	return &CommandContext{
		kv: kv,
		et: et,
	}
}

type Command interface {
	Execute(args []resp.Value, ctx *CommandContext) resp.Value
}
