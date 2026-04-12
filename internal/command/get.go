package command

import (
	"rodis/internal/protocol/resp"
)

type GetCommand struct{}

func (c *GetCommand) Execute(args []resp.Value, ctx *CommandContext) resp.Value {
	if len(args) != 1 {
		return resp.NewError("ERR wrong number of arguments for 'get' command")
	}
	if ctx == nil || ctx.k == nil {
		return resp.NewError("ERR internal server error")
	}

	key := args[0].Bulk

	value, ok := ctx.k.Get(key)

	if !ok {
		return resp.NewNullBulk()
	}

	return resp.NewBulk(value)
}
