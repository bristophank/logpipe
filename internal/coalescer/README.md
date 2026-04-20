# coalescer

The `coalescer` package merges multiple JSON log fields into a single target field by selecting the first non-empty value from a list of candidate fields.

## Use case

Different log producers may use different field names for the same concept (e.g. `msg`, `message`, `text`, `body`). The coalescer normalises these into a single canonical field.

## Rule fields

| Field        | Type     | Description                                              |
|--------------|----------|----------------------------------------------------------|
| `target`     | string   | Destination field name to write the resolved value into  |
| `candidates` | []string | Ordered list of source fields to check                   |
| `default`    | string   | Fallback value if no candidate has a non-empty value     |
| `delete_src` | bool     | If true, all candidate fields are removed after merging  |

## Example

```go
c := coalescer.New([]coalescer.Rule{
    {
        Target:     "message",
        Candidates: []string{"msg", "text", "body"},
        Default:    "(no message)",
        DeleteSrc:  true,
    },
})

out := c.Apply(`{"text":"hello world"}`)
// out: {"message":"hello world"}
```

## Behaviour

- Empty strings and missing fields are both treated as absent.
- If no candidate matches and no default is set, the target field is not added.
- Non-JSON lines and empty lines are passed through unchanged.
