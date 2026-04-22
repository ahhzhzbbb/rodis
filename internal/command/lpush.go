package command

import "rodis/internal/protocol/resp"

type LpushCommand struct{}

func (c *LpushCommand) Execute(args []resp.Payload, ctx *CommandContext) resp.Payload {
	if len(args) < 2 {
		return resp.NewError("ERR wrong number of arguments for 'lpush' command")
	}

	if ctx == nil || ctx.k == nil {
		return resp.NewError("ERR internal server error")
	}

	key := args[0].Bulk

	elements := args[1:]

	ctx.k.Set(key, "list", elements)

	return resp.NewInteger(len(args))
}
