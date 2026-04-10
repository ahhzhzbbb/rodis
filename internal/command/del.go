package command

import "rodis/internal/protocol/resp"

type DelCommand struct{}

func (c *DelCommand) Execute(args []resp.Value, ctx *CommandContext) resp.Value {
	var count int
	if len(args) == 0 {
		return resp.NewBulk("(error) ERR wrong number of arguments for 'del' command")
	}

	for _, arg := range args {
		key := arg.Bulk

		_, ok := ctx.et.Et[key]
		if ok {
			ctx.et.Mu.Lock()
			delete(ctx.et.Et, key)
			ctx.et.Mu.Unlock()
		}

		_, ok = ctx.kv.Kv[key]
		if ok {
			ctx.kv.Mu.Lock()
			delete(ctx.kv.Kv, key)
			ctx.kv.Mu.Unlock()
			count++
		}
	}
	return resp.NewInteger(count)
}
