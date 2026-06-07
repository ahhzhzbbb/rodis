package command

import "rodis/internal/protocol/resp"

type LInsertCommand struct{}

func (c *LInsertCommand) Execute(args []resp.Payload, ctx *CommandContext) resp.Payload {
	if len(args) != 4 {
		return resp.NewError("ERR wrong number of arguments for 'linsert' command")
	}
	if ctx == nil || ctx.k == nil {
		return resp.NewError("ERR internal server error")
	}

	key := args[0].Bulk
	position := args[1].Bulk
	pivot := args[2].Bulk
	element := args[3].Bulk

	switch position {
	case "BEFORE":
		pos = engine.ListInsertBefore
	case "AFTER":
		pos = engine.ListInsertAfter
	default:
		return resp.NewError("ERR syntax error")
	}

	count, err := ctx.k.ListInsert(key, pos, pivot, value)
	if err != nil {
		return resp.NewError(err.Error())
	}

	return resp.NewInteger(count)
}
