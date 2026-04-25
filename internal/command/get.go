package command

import (
	"rodis/internal/protocol/resp"
)

type GetCommand struct{}

func (c *GetCommand) Execute(args []resp.Payload, ctx *CommandContext) resp.Payload {
	if len(args) != 1 {
		return resp.NewError("ERR wrong number of arguments for 'get' command")
	}
	if ctx == nil || ctx.k == nil {
		return resp.NewError("ERR internal server error")
	}

	key := args[0].Bulk

	value, ok, err := ctx.k.GetString(key)
	if err != nil {
		return resp.NewError(err.Error())
	}

	if !ok {
		return resp.NewNullBulk()
	}

	return resp.NewBulk(value)
}
