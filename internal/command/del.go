package command

import "rodis/internal/protocol/resp"

type DelCommand struct{}

func (c *DelCommand) Execute(args []resp.Value, ctx *CommandContext) resp.Value {
	var count int
	if len(args) == 0 {
		return resp.NewError("ERR wrong number of arguments for 'del' command")
	}

	for _, arg := range args {
		key := arg.Bulk

		ctx.et.Mu.RLock()
		_, ok := ctx.et.Et[key]
		ctx.et.Mu.RUnlock()
		if ok {
			ctx.et.Mu.Lock()
			delete(ctx.et.Et, key)
			ctx.et.Mu.Unlock()
		}

		ctx.kv.Mu.RLock()
		_, ok = ctx.kv.Kv[key]
		ctx.kv.Mu.RUnlock()
		if ok {
			ctx.kv.Mu.Lock()
			delete(ctx.kv.Kv, key)
			ctx.kv.Mu.Unlock()
			count++
		}
	}
	return resp.NewInteger(count)
}
