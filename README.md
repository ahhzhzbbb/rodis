# rodis

## project structure
```
.
├── cmd
│   ├── redis-cli
│   └── redis-server
│       ├── dump.rdb
│       └── main.go
├── go.mod
├── internal
│   ├── command
│   │   ├── commandDocs.go
│   │   ├── command.go
│   │   ├── config.go
│   │   ├── del.go
│   │   ├── exists.go
│   │   ├── expire.go
│   │   ├── get.go
│   │   ├── incr.go
│   │   ├── ping.go
│   │   └── set.go
│   ├── engine
│   │   └── keyValue.go
│   ├── factory
│   │   └── factory.go
│   ├── protocol
│   │   └── resp
│   │       ├── encoder.go
│   │       ├── parser.go
│   │       └── value.go
│   └── server
│       ├── config.go
│       ├── handler.go
│       └── server.go
├── README.md
├── redis-server
├── roadmap
└── temp

11 directories, 25 files
```