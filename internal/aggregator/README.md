# aggregator

The `aggregator` package groups structured log lines by a chosen field and emits count summaries at a configurable flush interval.

## Usage

```go
a := aggregator.New("level", 5*time.Second, func(line string) {
    fmt.Println(line) // {"field":"level","value":"error","count":42}
})

a.Add(`{"level":"error","msg":"something failed"}`)
a.Add(`{"level":"info","msg":"all good"}`)

// manual flush (window=0 disables auto-flush)
a.Flush()

// stop background goroutine when window > 0
a.Stop()
```

## Output format

Each flushed summary is a JSON object:

```json
{"field":"level","value":"error","count":3}
```

## Notes

- Lines with invalid JSON or a missing group field are silently ignored.
- Pass `window = 0` to disable the background ticker and call `Flush()` manually.
- `Stop()` must be called when using a non-zero window to avoid goroutine leaks.
