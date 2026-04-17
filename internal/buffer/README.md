# buffer

Package `buffer` provides a thread-safe, fixed-capacity ring buffer for log lines.

## Behaviour

- **Capacity**: set at construction time via `New(capacity int)`.
- **Overflow**: when the buffer is full, the oldest line is silently overwritten and the internal `dropped` counter is incremented.
- **Flush**: returns all buffered lines in insertion order and resets the buffer to empty.
- **Thread-safe**: all operations are protected by a mutex and safe for concurrent use.

## Usage

```go
buf := buffer.New(1000)

// producer
buf.Write(logLine)

// consumer (e.g. on a ticker)
for _, line := range buf.Flush() {
    sink.Write(line)
}

// observability
fmt.Println("dropped:", buf.Dropped())
```

## Integration

The buffer sits between the pipeline and slow sinks (e.g. file or network sinks) to absorb bursts without blocking the main read loop. The `metrics` collector can periodically sample `Dropped()` to surface overflow events.
