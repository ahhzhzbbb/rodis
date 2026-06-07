# rodis

## project structure
```
.
├── cmd
│   ├── rodis-cli
│   └── rodis-server
│       ├── dump.rdb
│       └── main.go
├── dump.rdb
├── go.mod
├── go.sum
├── internal
│   ├── command
│   │   ├── append.go
│   │   ├── commandDocs.go
│   │   ├── command.go
│   │   ├── config.go
│   │   ├── del.go
│   │   ├── exists.go
│   │   ├── expire.go
│   │   ├── get.go
│   │   ├── incr.go
│   │   ├── lpop.go
│   │   ├── lpush.go
│   │   ├── lrange.go
│   │   ├── ping.go
│   │   ├── rpop.go
│   │   ├── rpush.go
│   │   └── set.go
│   ├── engine
│   │   ├── keyValue.go
│   │   ├── list.go
│   │   ├── object.go
│   │   ├── quickList.go
│   │   ├── shardmap.go
│   │   └── zipList.go
│   ├── factory
│   │   └── factory.go
│   ├── protocol
│   │   └── resp
│   │       ├── encoder.go
│   │       ├── parser.go
│   │       └── payload.go
│   └── server
│       ├── config.go
│       ├── handler.go
│       └── server.go
├── logs
│   ├── bench.log
│   ├── error.log
│   └── strace.log
├── README.md
├── roadmap
├── scripts
│   └── strace.sh
└── temp

13 directories, 41 files
```