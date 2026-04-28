package command

import "rodis/internal/protocol/resp"

type LpopCommand struct{}

func (c *LpopCommand) Execute(args []resp.Payload, ctx *CommandContext) resp.Payload {
	if len(args) > 2 || len(args) == 0 {
		return resp.NewError("ERR wrong number of arguments for 'lrange' command")
	}

	if ctx == nil || ctx.k == nil {
		return resp.NewError("ERR internal server error")
	}

	key := args[0].Bulk
	count := "1"
	if len(args) == 2 {
		count = args[1].Bulk
	}
	elements, poped, err := ctx.k.PopList(key, count, true)
	if err != nil {
		return resp.NewError(err.Error())
	}
	if !poped {
		return resp.NewNullBulk()
	}

	var res []resp.Payload
	for _, e := range elements {
		res = append(res, resp.NewBulk(e))
	}
	return resp.NewArray(res)
}
