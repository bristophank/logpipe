# clamper

Constrains numeric JSON field values to a configured `[min, max]` range.

## Usage

```go
import "github.com/yourorg/logpipe/internal/clamper"

c := clamper.New([]clamper.Rule{
    {Field: "score",   Min: 0,   Max: 100},
    {Field: "latency", Min: 1,   Max: 5000},
})

out, err := c.Apply(`{"score":150,"latency":-1}`)
// out: {"latency":1,"score":100}
```

## Rules

| Field   | Type    | Description                          |
|---------|---------|--------------------------------------|
| `field` | string  | JSON key to clamp                    |
| `min`   | float64 | Lower bound (inclusive)              |
| `max`   | float64 | Upper bound (inclusive)              |

## Behaviour

- Values already within `[min, max]` are left untouched.
- Values below `min` are set to `min`; values above `max` are set to `max`.
- Fields absent from the log line are silently skipped.
- Non-numeric field values are skipped without error.
- Lines that are not valid JSON are passed through unchanged.
