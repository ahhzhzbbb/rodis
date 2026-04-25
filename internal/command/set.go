package command

import "rodis/internal/protocol/resp"

type SetCommand struct{}

func (c *SetCommand) Execute(args []resp.Payload, ctx *CommandContext) resp.Payload {
	if len(args) != 2 {
		return resp.NewError("ERR wrong number of arguments for 'set' command")
	}

	if ctx == nil || ctx.k == nil {
		return resp.NewError("ERR internal server error")
	}

	key := args[0].Bulk

	value := args[1].Bulk

	err := ctx.k.SetString(key, value)
	if err != nil {
		return resp.NewError(err.Error())
	}
	return resp.NewString("OK")
}
