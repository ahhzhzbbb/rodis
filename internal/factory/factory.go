package factory

import "rodis/internal/command"

type Factory interface {
	CreateCommand() command.Command
}

var CommandRegistry = map[string]func() command.Command{
	"PING":    func() command.Command { return &command.PingCommand{} },
	"COMMAND": func() command.Command { return &command.CommandDocs{} },
	"GET":     func() command.Command { return &command.GetCommand{} },
	"SET":     func() command.Command { return &command.SetCommand{} },
	"DEL":     func() command.Command { return &command.DelCommand{} },
	"EXISTS":  func() command.Command { return &command.ExistsCommand{} },
	"EXPIRE":  func() command.Command { return &command.ExpireCommand{} },
	"CONFIG":  func() command.Command { return &command.ConfigCommand{} },
	"INCR":    func() command.Command { return &command.IncrCommand{} },
	"APPEND":  func() command.Command { return &command.AppendCommand{} },
}
