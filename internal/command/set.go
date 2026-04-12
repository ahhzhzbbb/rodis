package command

import "rodis/internal/protocol/resp"

type SetCommand struct{}

func (c *SetCommand) Execute(args []resp.Value, ctx *CommandContext) resp.Value {
	if len(args) != 2 {
		return resp.NewError("ERR wrong number of arguments for 'set' command")
	}
	if ctx == nil || ctx.kv == nil {
		return resp.NewError("ERR internal server error")
	}

	key := args[0].Bulk

	value := args[1].Bulk

	ctx.et.Mu.RLock()
	_, ok := ctx.et.Et[key]
	ctx.et.Mu.RUnlock()
	if ok {
		ctx.et.Mu.Lock()
		delete(ctx.et.Et, key)
		ctx.et.Mu.Unlock()
	}

	ctx.kv.Mu.Lock()
	ctx.kv.Kv[key] = value
	ctx.kv.Mu.Unlock()

	return resp.NewString("OK")
}
