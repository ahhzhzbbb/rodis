package command

import "rodis/internal/protocol/resp"

type CommandContext struct{}

type Command interface {
	Execute(args []resp.Value, ctx *CommandContext) resp.Value
}
