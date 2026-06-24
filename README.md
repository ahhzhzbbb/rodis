<div align="center">

  ![Rodis Logo]([https://via.placeholder.com/150](https://images.viblo.asia/2a451245-3e33-415a-9ae2-93339784df41.png))

  # 🚀 Rodis

  <p><strong>A Redis-compatible in-memory data structure store written in Go</strong></p>

  <p>
    <img src="https://badges.aleen42.com/golang.svg" alt="Go" />
    <img src="https://img.shields.io/github/go-mod/go-version/hoangmp/rodis" alt="Go Version" />
    <img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT" />
    <img src="https://img.shields.io/badge/Protocol-RESP-green.svg" alt="Protocol: RESP" />
  </p>

  <p>
    <a href="#features">Features</a> •
    <a href="#architecture">Architecture</a> •
    <a href="#quick-start">Quick Start</a> •
    <a href="#supported-commands">Commands</a> •
    <a href="#performance">Performance</a> •
    <a href="#benchmarks">Benchmarks</a>
  </p>

</div>

---

## 📋 Overview

Rodis is a high-performance, Redis-compatible in-memory data store written in Go. It implements the **RESP (REdis Serialization Protocol)** and a custom storage engine designed for efficient string and list operations. Rodis features thread-safe concurrent access, a custom `QuickList` data structure, AOF persistence, and optimized expiration policies.

## 🚀 Features

- **🔑 Key-Value Storage**: High-performance string operations with support for `SET`, `GET`, `DEL`, `EXISTS`, `INCR`, and `APPEND`.
- **📑 List Operations**: Advanced list support using a custom `QuickList` implementation (`LPUSH`, `RPUSH`, `LPOP`, `RPOP`, `LRANGE`, `LINSERT`).
- **⏱️ TTL / Expiration**: Both lazy (on access) and active (timed background task) key expiration.
- **🛡️ Thread-Safe**: Built on a sharded hash map with fine-grained locking for high concurrency.
- **💾 Persistence (AOF)**: Append-Only File logging of mutations for data durability.
- **🔌 RESP Protocol**: Communicates via the standard RESP protocol, allowing standard Redis clients to be used.

---

## 🏗️ Architecture

Rodis is built with a modular and clean architecture, separating concerns between the network protocol, the core command engine, and the internal data structures.

```
┌────────────────────┐
│   Client (e.g.,   │
│   redis-cli)      │
└────────┬───────────┘
         │ TCP (:6379)
         ▼
┌────────────────────┐
│   Server (Handler) │  Handles network I/O using Go's goroutines
│   (TCP, RESP)      │  ── cmd/rodis-server/main.go ──
└────────┬───────────┘
         │
         ▼
┌────────────────────┐
│  Protocol (RESP)     │  Parses and encodes RESP streams (bulk strings,
│  Parser / Encoder    │  arrays, errors) into Go structs.
│   internal/protocol/resp  │
└────────┬───────────┘
         │
         ▼
┌────────────────────┐
│   Command Layer    │  Maps command names (e.g., "SET") to logic
│  (Factory Pattern) │  via a command registry (internal/factory).
│   internal/command │
└────────┬───────────┘
         │
         ▼
┌────────────────────┐
│   Storage Engine     │  Thread-safe sharded hash map for key-value
│   (KeyValue + lists) │  with custom data structures (QuickList, ZipList).
└────────────────────┘
```

### Core Data Structures

- **Sharded Hash Map (`internal/engine/shardmap.go`)**: A high-performance, thread-safe hash map built on top of `tidwall/rhh` (Robin Hood Hashing). It partitions data across multiple shards to minimize lock contention and maximize throughput.
- **QuickList (`internal/engine/quickList.go`)**: Rodis's custom implementation of a Redis List. It is a linked list of `ZipList` nodes, optimizing memory usage while maintaining good performance for both random access and sequential traversal.
- **ZipList (`internal/engine/zipList.go`)**: A compact, byte-packed list that reduces memory overhead by storing data contiguously. It is used to store many strings within a single node to boost memory efficiency.
- **AOF (`internal/engine/aof.go`)**: The persistence layer that appends every write command to a file for data durability.

---

## ⚡ Quick Start

### Prerequisites

- Go **1.22** or later installed on your system.
- A Redis client (such as `redis-cli`) for testing.

### Installation & Running

1. **Clone the repository**

   ```bash
   git clone https://github.com/hoangmp/rodis.git
   cd rodis
   ```

2. **Run the server**

   ```bash
   cd cmd/rodis-server
   go run main.go
   ```

   You should see the welcome banner and the server starting up on port `6379`:
   ```
   Welcome to Rodis! Server is running on port :6379
    ____   ___  ____ ___ ____
   |  _ \ / _ \|  _ \_ _/ ___|
   | |_) | | | | | | | |\___ \
   |  _ <| |_| | |_| | | ___) |
   |_| \_\\___/|____/___|____/ 
   ```

3. **Connect with a client**

   ```bash
   redis-cli -p 6379
   ```

4. **Try some commands**

   ```bash
   127.0.0.1:6379> PING
   "PONG"
   127.0.0.1:6379> SET mykey "hello"
   OK
   127.0.0.1:6379> GET mykey
   "hello"
   127.0.0.1:6379> LPUSH mylist a b c
   (integer) 3
   127.0.0.1:6379> LRANGE mylist 0 -1
   1) "c"
   2) "b"
   3) "a"
   ```

---

## 🖥️ Supported Commands

Rodis implements a key subset of the Redis command set.

### String Commands

| Command | Description | Status |
| :--- | :--- | :--- |
| `GET key` | Get the value of a key. | ✅ |
| `SET key value` | Set the string value of a key. | ✅ |
| `INCR key` | Increment the integer value of a key by one. | ✅ |
| `APPEND key value` | Append a value to a key. | ✅ |
| `EXISTS key [...keys]` | Check if a key exists. | ✅ |
| `DEL key [...keys]` | Delete one or more keys. | ✅ |
| `EXPIRE key seconds` | Set a key's time to live in seconds. | ✅ |

### List Commands

| Command | Description | Status |
| :--- | :--- | :--- |
| `LPUSH key element [element ...]` | Insert all the specified values at the head of the list stored at key. | ✅ |
| `RPUSH key element [element ...]` | Insert all the specified values at the tail of the list. | ✅ |
| `LPOP key [count]` | Remove and get the first elements in a list. | ✅ |
| `RPOP key [count]` | Remove and get the last elements in a list. | ✅ |
| `LRANGE key start stop` | Get a range of elements from a list. | ✅ |
| `LINSERT key BEFORE/AFTER pivot element` | Insert an element before or after a pivot in a list. | ✅ |
| `CONFIG` | Get configuration parameters. | ✅ |

---

## 📈 Performance

Rodis is designed with performance in mind. In benchmarks, it demonstrates the `QuickList`'s efficiency by achieving competitive throughput, thanks to its thread-safe, sharded architecture and the hybrid `QuickList` structure.

### Key Optimizations

- **Sharded Map**: Distributes keys across `N` shards (up to `runtime.NumCPU() * 16`) to avoid the single global lock throttling heavy concurrent workloads.
- **QuickList Amortization**: The `QuickList`'s hierarchical structure (linked list of pipelines) provides O(1) push/pop with O(N) scans for head-tail, where `engine` of particular are `using respon` with `resp.` for `from the`.
* **Pipeline Support**: Supports pipelining requests over TCP for high-throughput workloads.

### Benchmark Comparison (from `internal/benchmark/usageList_test.go`)

| Metric | LinkedList | ZipList | QuickList |
| :--- | :--- | :--- | :--- |
| **Memory Usage** | High | Low | Low |
| **Push/Pop (head/tail)** | Fast | Fast (tail), Slow (head) | Fast |
| **Sequential Traversal** | Slow | Fast | Fast |
| **Cache Locality** | Poor | Good | Good |

---

## 💻 Development

### Project Structure

```
rodis/
├── cmd/
│   └── rodis-server/
│       └── main.go          # Entry point for the server
├── internal/
│   ├── benchmark/
│   │   └── usageList_test.go  # Benchmarks for QuickList vs ZipList
│   ├── command/             # Command implementations (SET, GET, etc.)
│   ├── engine/              # Core storage (QuickList, ZipList, ShardedMap, AOF)
│   ├── factory/             # Command registry (Factory pattern)
│   ├── protocol/resp/       # RESP parser and encoder
│   └── server/              # TCP server and request handler
├── scripts/
│   └── strace.sh            # System call analysis script
├── logs/                    # Server logs
├── go.mod
└── README.md
```

### Running Tests

```bash
cd rodis
go test ./...
```

---

## 📜 Roadmap

- [x] Implement core String commands (SET, GET, DEL, INCR, APPEND)
- [x] Implement core List commands (LPUSH, RPUSH, LPOP, RPOP, LRANGE, LINSERT)
- [x] Implement AOF (Append-Only File) persistence
- [x] Lazy and Active Expiration policies
- [ ] Implement Hash data structure (`HSET`, `HGET`, `HDEL`)
- [ ] Implement Set data structure (`SADD`, `SMEMBERS`, `SISMEMBER`)
- [ ] Implement Pub/Sub messaging
- [ ] Replication (Leader/Follower)
- [ ] Full support for Redis Streams

---

## 🤝 Contributing

Contributions are welcome! If you have a feature request, bug report, or want to improve the code:

1.  Fork this repository.
2.  Create a new branch for your feature or fix.
3.  Submit a pull request detailing your changes.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<p align="center">Built with ❤️ by <a href="mailto:ahhzhzbbb@gmail.com">hoangmp</a></p>
