package command

import "rodis/internal/protocol/resp"

type ExistsCommand struct{}

func (c *ExistsCommand) Execute(args []resp.Value, ctx *CommandContext) resp.Value {
	if len(args) == 0 {
		return resp.NewError("ERR wrong number of arguments for 'exists' command")
	}

	var result int

	for _, arg := range args {
		ctx.kv.Mu.RLock()
		_, ok := ctx.kv.Kv[arg.Bulk]
		ctx.kv.Mu.RUnlock()
		if ok {
			result++
		}
	}
	return resp.NewInteger(result)
}
