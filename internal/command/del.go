package command

import "rodis/internal/protocol/resp"

type DelCommand struct{}

func (c *DelCommand) Execute(args []resp.Value, ctx *CommandContext) resp.Value {
	var count int
	if len(args) == 0 {
		return resp.NewBulk("(error) ERR wrong number of arguments for 'del' command")
	}

	for _, arg := range args {
		_, ok := ctx.kv.Kv[arg.Bulk]
		if ok {
			ctx.kv.Mu.Lock()
			delete(ctx.kv.Kv, arg.Bulk)
			ctx.kv.Mu.Unlock()
			count++
		}
	}
	return resp.NewInteger(count)
}
