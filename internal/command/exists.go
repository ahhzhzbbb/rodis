package command

import "rodis/internal/protocol/resp"

type ExistsCommand struct{}

func (c *ExistsCommand) Execute(args []resp.Value, ctx *CommandContext) resp.Value {
	if len(args) == 0 {
		return resp.NewBulk("(error) ERR wrong number of arguments for 'exists' command")
	}

	var result int

	for _, arg := range args {
		_, ok := ctx.kv.Kv[arg.Bulk]
		if ok {
			result++
		}
	}
	return resp.NewInteger(result)
}
