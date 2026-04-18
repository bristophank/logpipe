# alerter

Fires structured alerts when a numeric log field exceeds a configured threshold within a rolling time window.

## Rules

Each rule specifies:

| Field | Type | Description |
|-------|------|-------------|
| `field` | string | JSON key to inspect |
| `threshold` | float64 | Value must exceed this to trigger |
| `window_seconds` | int | Rolling window for counting hits |
| `sink` | string | Optional label for grouping |

## Usage

```go
rules := []alerter.Rule{
    {Field: "latency_ms", Threshold: 1000, Window: 30, SinkName: "slow"},
}
a := alerter.New(rules)
a.Check(line, os.Stderr)
```

Alerts are written as JSON lines to the provided `io.Writer`.
