package command

import (
	"rodis/internal/protocol/resp"
)

type LrangeCommand struct{}

func (c *LrangeCommand) Execute(args []resp.Payload, ctx *CommandContext) resp.Payload {
	key := args[0].Bulk
	elements := ctx.k.GetList(key)
	var res []resp.Payload
	for _, e := range elements {
		res = append(res, resp.NewBulk(e))
	}
	return resp.NewArray(res)
}
