# truncator

The `truncator` package shortens string fields in JSON log lines that exceed a configured maximum byte length.

## Usage

```go
tr := truncator.New([]truncator.Rule{
    {Field: "msg",  MaxLen: 128},
    {Field: "body", MaxLen: 512},
})

out := tr.Apply(line)
```

## Behaviour

- Fields not present in the log line are silently skipped.
- Non-string fields are left untouched.
- Lines that are not valid JSON are returned as-is.
- If no rules are configured the line is returned unchanged.
