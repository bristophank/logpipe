# throttle

The `throttle` package suppresses repeated identical log lines within a configurable cooldown window.

## Behaviour

- If `cooldown` is zero, all lines pass through unconditionally.
- The first occurrence of a line is always forwarded.
- Subsequent identical lines are suppressed until the cooldown duration has elapsed since the last forwarded occurrence.
- Different lines are tracked independently.

## Usage

```go
th := throttle.New(5 * time.Second)

for _, line := range lines {
    if th.Allow(line) {
        fmt.Println(line)
    }
}
```

## Config fields

| Field      | Type     | Description                          |
|------------|----------|--------------------------------------|
| `cooldown` | duration | Minimum gap between identical lines  |
