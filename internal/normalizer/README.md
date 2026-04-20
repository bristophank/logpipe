# normalizer

The `normalizer` package standardises field names and values in structured (JSON) log lines.

## Rules

Each rule targets a single field and applies one transformation:

| Transform   | Effect                                                   |
|-------------|----------------------------------------------------------|
| `lowercase` | Converts the field value to lower-case                   |
| `uppercase` | Converts the field value to upper-case                   |
| `trim`      | Strips leading and trailing whitespace from the value    |
| `snake_case`| Lowercases and replaces spaces/hyphens with underscores  |

Optionally, set `rename` to move the result into a different key.

## Example

```go
n := normalizer.New([]normalizer.Rule{
    {Field: "level",   Transform: "lowercase"},
    {Field: "service", Transform: "snake_case", Rename: "svc"},
})

out := n.Apply(`{"level":"ERROR","service":"My-Service"}`)
// out → {"level":"error","svc":"my_service"}
```

## Behaviour

- Lines that are not valid JSON are passed through unchanged.
- Fields that do not exist in the log line are silently skipped.
- Non-string field values are left untouched.
