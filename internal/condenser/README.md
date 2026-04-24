# condenser

The `condenser` package merges consecutive JSON log lines that share the same
value for a configured key field into a single representative output line.

## Use case

High-volume log streams often emit many repeated entries with the same
severity or category in a row. The condenser collapses these runs so
downstream sinks receive one line per group, optionally annotated with how
many source lines were merged.

## Rules

| Field | Type | Description |
|-------------|--------|-----------------------------------------------------|
| `field` | string | JSON key used to detect consecutive duplicate groups |
| `count_field` | string | Optional key added to the output line with the merge count |

## Behaviour

- Lines with no matching rules are passed through unchanged.
- Invalid JSON lines are passed through unchanged.
- Empty / whitespace-only lines are dropped.
- `Add` returns a non-empty string only when a group boundary is crossed.
- `Flush` must be called after the stream ends to emit the final group.

## Example

```go
c := condenser.New([]condenser.Rule{
    {Field: "level", CountField: "_count"},
})

for _, raw := range lines {
    if out := c.Add(raw); out != "" {
        fmt.Println(out)
    }
}
if out := c.Flush(); out != "" {
    fmt.Println(out)
}
```
