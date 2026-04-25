# bouncer

The `bouncer` package filters structured log lines based on per-field allow and block lists.

## Rules

Each rule targets a single JSON field and may specify:

- `allow` — only lines whose field value appears in this list are passed.
- `block` — lines whose field value appears in this list are dropped.

Both lists may be combined on the same rule. Block evaluation happens first.

## Behaviour

- Lines with no matching field are **passed through** unchanged.
- Invalid JSON lines are **passed through** unchanged.
- All rules must pass for a line to be allowed.

## Example

```go
rules := []bouncer.Rule{
    {Field: "level", Allow: []string{"info", "warn", "error"}},
    {Field: "env",   Block: []string{"staging"}},
}
b := bouncer.New(rules)

if b.Allow(line) {
    fmt.Println(line)
}
```
