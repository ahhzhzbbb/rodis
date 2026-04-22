package command

import "rodis/internal/protocol/resp"

type AppendCommand struct{}

func (c *AppendCommand) Execute(args []resp.Payload, ctx *CommandContext) resp.Payload {
	if len(args) != 2 {
		return resp.NewError("ERR wrong number of arguments for 'append' command")
	}

	if ctx == nil || ctx.k == nil {
		return resp.NewError("ERR internal server error")
	}

	key := args[0].Bulk

	value := args[1].Bulk

	result, err := ctx.k.Append(key, value)
	if err != nil {
		return resp.NewError("ERR this is not string")
	}

	return resp.NewInteger(result)
}
