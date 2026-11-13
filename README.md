# 🐝 Worker Colony

**Async Task Queue & Worker Engine in Go**

A concurrent task processing system implementing the worker pool pattern. It demonstrates practical application of Go's concurrency primitives to build a thread-safe, asynchronous job processor with strict state management.

## 🏗️ Architecture

The project follows standard Go layout conventions, separating executable code from internal library logic.

```text
worker-colony/
├── cmd/server/              # Application entry point
│   └── main.go             # HTTP server setup & worker initialization
├── internal/
│   ├── handler/            # HTTP handlers & worker orchestration
│   │   └── handler.go      # TaskHandler: queue management, worker lifecycle
│   ├── memory-storage/     # Thread-safe storage layer
│   │   ├── memory.go       # InMemoryStorage with sync.RWMutex
│   │   └── memory_test.go  # Storage concurrency tests
│   └── task/               # Domain models & state machine
│       ├── task.go         # Task struct & transition rules
│       └── task_test.go    # State machine validation
├── go.mod
└── README.md
```

## 💡 Core Concepts & Design

### The Worker Pool Pattern
Instead of spawning a new goroutine for every incoming request, Worker Colony uses a fixed set of workers that consume tasks from a shared channel. This approach prevents resource exhaustion and provides predictable performance under load.

### Concurrency & Safety
- **Backpressure via Buffered Channels**: The task queue is a buffered channel. When full, it naturally rejects new work (`503 Service Unavailable`), preventing unbounded memory growth.
- **Thread-Safe State Management**: An in-memory map protected by `sync.RWMutex` allows high-concurrency reads (status checks) while ensuring exclusive access for writes (state updates).
- **Lock-Free Operations**: Task IDs are generated using `atomic.Int64`, avoiding mutex overhead for simple counters.

### Deterministic State Machine
Tasks follow a strict lifecycle: `Pending → Running → Completed/Failed`. Invalid transitions are rejected at the code level, ensuring data integrity without complex database constraints.

## 🚀 Quick Start

### Prerequisites
- Go 1.21+ (required for `PathValue` in HTTP routing)

### Build & Run
```bash
# From project root
go run ./cmd/server/main.go
```

This starts:
- HTTP server on `localhost:8080`
- 3 worker goroutines processing tasks concurrently 
(Both the HTTP Port and amount of Workers adjustable)

### API Usage

**Submit a Task:**
```bash
curl -X POST http://localhost:8080/task \
  -d "Process this data"
```
*Response: `202 Accepted` with Task ID*
```text
task is being processed => ID: 1
```

**Check Task Status:**
```bash
curl http://localhost:8080/task/1
```
*Response:*
```json
{
  "ID": 1,
  "Status": 2,
  "CreatedAt": "2026-11-13T...",
  "Description": "Process this data"
}
```
*Status codes: `0=Pending`, `1=Running`, `2=Completed`, `3=Failed`*

## 🧪 Testing

```bash
go test ./...
```

Tests cover:
- Task state transition validity
- Concurrent storage access patterns
- Memory storage CRUD operations

## 🔮 Future Enhancements

- [ ] Persistent storage backend (PostgreSQL/Redis)
- [ ] Priority queue implementation
- [ ] Worker health monitoring & auto-scaling
- [ ] Retry logic with exponential backoff
- [ ] Prometheus metrics integration
- [ ] Graceful shutdown handling

## 📄 License

MIT

---

**Built by [Boalq](https://github.com/Boalq)**