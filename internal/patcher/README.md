# patcher

The `patcher` package conditionally sets fields in JSON log lines based on match rules.

## Rules

Each rule specifies:

| Field    | Description                                          |
|----------|------------------------------------------------------|
| `field`  | The JSON field to evaluate                           |
| `op`     | Match operator: `eq`, `contains`, or `exists`        |
| `value`  | Value to match (not used for `exists`)               |
| `target` | The field to write when the condition matches        |
| `patch`  | The string value to assign to `target`               |

## Operators

- **eq** — exact string equality
- **contains** — substring match
- **exists** — field is present (any value)

## Usage

```go
rules := []patcher.Rule{
    {
        Field:  "level",
        Op:     "eq",
        Value:  "error",
        Target: "alert",
        Patch:  "true",
    },
}
p := patcher.New(rules)
out := p.Apply(`{"level":"error","msg":"something failed"}`)
// out: {"level":"error","msg":"something failed","alert":"true"}
```

## Notes

- Lines that are not valid JSON are passed through unchanged.
- Multiple rules are evaluated independently; all matching rules are applied.
- If no rules are configured, every line is returned as-is.
