# counter

The `counter` package tracks the frequency of distinct values for one or more
JSON fields across a stream of log lines.

## Usage

```go
rules := []counter.Rule{
    {Field: "level"},
    {Field: "service", Alias: "svc"},
}
c := counter.New(rules)

// Feed individual lines
_ = c.Add(`{"level":"info","service":"api"}`)
_ = c.Add(`{"level":"error","service":"worker"}`)
_ = c.Add(`{"level":"info","service":"api"}`)

// Read current counts (keyed by alias when set)
snap := c.Snapshot()
// snap["level"]["info"]  == 2
// snap["svc"]["api"]     == 2

// Reset for the next window
c.Reset()
```

## Streaming

```go
err := c.Stream(os.Stdin, os.Stdout)
```

All lines are forwarded to the writer unchanged. Invalid JSON lines are passed
through without interrupting the stream.

## Rules

| Field   | Type   | Description                              |
|---------|--------|------------------------------------------|
| `field` | string | JSON key to count distinct values for    |
| `alias` | string | Output key in `Snapshot()`; optional     |
