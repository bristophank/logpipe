# sequencer

The `sequencer` package injects monotonically increasing sequence numbers into structured JSON log lines.

## Usage

```go
seq := sequencer.New([]sequencer.Rule{
    {Field: "seq", Start: 1, Step: 1},
})

out, err := seq.Apply(`{"msg":"hello"}`)
// out: {"msg":"hello","seq":1}
```

## Rule fields

| Field   | Type | Description                          |
|---------|------|--------------------------------------|
| `Field` | string | JSON key to write the counter into |
| `Start` | int  | Initial counter value (default `0`) |
| `Step`  | int  | Increment per line (default `1`)    |

## Streaming

```go
err := seq.Stream(os.Stdin, os.Stdout)
```

Empty lines are skipped. Lines with invalid JSON are passed through unchanged.

## Thread safety

The sequencer is safe for concurrent use. Multiple goroutines may call `Apply` simultaneously.
