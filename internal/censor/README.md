# censor

The `censor` package replaces specific field values in JSON log lines with a configurable mask string.

## Use Case

Use `censor` to suppress known sensitive values such as usernames, roles, or status codes before logs are forwarded to external sinks.

## Rules

Each rule specifies:

| Field    | Description                                          |
|----------|------------------------------------------------------|
| `field`  | The JSON key to inspect                              |
| `values` | List of values to censor (case-insensitive match)    |
| `mask`   | Replacement string (defaults to `[CENSORED]`)        |

## Example

```go
c := censor.New([]censor.Rule{
    {Field: "role", Values: []string{"admin", "root"}, Mask: "[REDACTED]"},
})

out, err := c.Apply(`{"role":"admin","action":"drop_table"}`)
// out: {"action":"drop_table","role":"[REDACTED]"}
```

## Notes

- Non-string field values are skipped.
- Lines with invalid JSON are returned unchanged with an error.
- If no rules are configured, all lines pass through unmodified.
