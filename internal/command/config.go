package command

import "rodis/internal/protocol/resp"

type ConfigCommand struct{}

func (c *ConfigCommand) Execute(args []resp.Value, ctx *CommandContext) resp.Value {
	temp := make([]resp.Value, 0)
	return resp.NewArray(temp)
}
