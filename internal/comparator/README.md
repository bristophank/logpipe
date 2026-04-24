# comparator

The `comparator` package evaluates numeric comparisons on JSON log fields and writes a boolean result field.

## Rules

Each rule specifies:
- `field` — source numeric field to read
- `op` — comparison operator: `gt`, `lt`, `gte`, `lte`, `eq`, `neq`
- `value` — the numeric threshold to compare against
- `target` — destination field for the boolean result (defaults to `<field>_<op>`)

## Example

```go
c := comparator.New([]comparator.Rule{
    {Field: "latency_ms", Op: "gt", Value: 500, Target: "is_slow"},
    {Field: "status",     Op: "eq", Value: 200, Target: "is_ok"},
})

out, err := c.Apply(`{"latency_ms": 750, "status": 200}`)
// out: {"latency_ms":750,"status":200,"is_slow":true,"is_ok":true}
```

## Behaviour

- Non-numeric fields are silently skipped.
- Invalid JSON is returned unchanged with an error.
- String-encoded numbers (e.g. `"42"`) are coerced to float before comparison.
- Multiple rules are applied in order; later rules may overwrite earlier results if they share a target field.
