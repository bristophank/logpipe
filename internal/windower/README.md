# windower

Accumulates numeric field values within a fixed **tumbling window** and emits a
summary JSON line when the window expires.

## Rules

| Field    | Type            | Description                              |
|----------|-----------------|------------------------------------------|
| `field`  | `string`        | Source numeric field to accumulate       |
| `window` | `time.Duration` | Length of the tumbling window            |
| `alias`  | `string`        | Output field name in the summary line    |

## Behaviour

- `Add(line)` ingests a JSON log line. If the elapsed time since the first
  ingested line exceeds `window`, the accumulated totals are flushed and the
  method returns the summary line plus `true`.
- `Flush()` can be called at any time to emit and reset the current bucket.
- Non-numeric values and missing fields are silently ignored.
- Lines with invalid JSON are skipped without error.

## Example

```go
w := windower.New([]windower.Rule{
    {Field: "bytes", Alias: "total_bytes", Window: time.Minute},
})

if summary, ok := w.Add(line); ok {
    fmt.Println(summary) // {"total_bytes":12345}
}
```
