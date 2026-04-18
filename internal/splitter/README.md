# splitter

The `splitter` package fans out a single structured log line to one or more
named output sinks based on a JSON field value.

## Rules

Each rule specifies:

| Field   | Description                              |
|---------|------------------------------------------|
| `field` | JSON key to inspect                      |
| `value` | Expected string value                    |
| `sink`  | Destination sink name                    |

Multiple rules can match the same line — it will be written to all matching
sinks.

## Fallback

If no rule matches and a fallback `io.Writer` is provided, the line is written
there instead.

## Example

```go
rules := []splitter.Rule{
    {Field: "level", Value: "error", Sink: "stderr"},
    {Field: "service", Value: "auth",  Sink: "auth-log"},
}
s := splitter.New(rules, sinks, fallback)
s.Write(line)
```
