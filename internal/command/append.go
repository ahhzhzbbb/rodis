package command

import "rodis/internal/protocol/resp"

type AppendCommand struct{}

func (c *AppendCommand) Execute(args []resp.Value, ctx *CommandContext) resp.Value {
	if len(args) != 2 {
		return resp.NewError("ERR wrong number of arguments for 'append' command")
	}

	if ctx == nil || ctx.k == nil {
		return resp.NewError("ERR internal server error")
	}

	key := args[0].Bulk

	value := args[1].Bulk

	oldValue, ok := ctx.k.Get(key)
	if ok {
		value = oldValue + value
		ctx.k.DelValue(key)
	}

	ctx.k.Set(key, value)
	return resp.NewInteger(len(value))
}
