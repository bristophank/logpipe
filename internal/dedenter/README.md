# dedenter

The `dedenter` package strips leading whitespace (spaces and tabs) from string
fields inside structured JSON log lines. It can also collapse interior
whitespace runs to a single space.

## Rules

| Field           | Type   | Description                                              |
|-----------------|--------|----------------------------------------------------------|
| `field`         | string | JSON key whose string value will be dedented.            |
| `collapse_inner`| bool   | When `true`, interior whitespace runs become one space.  |

## Example

```go
d := dedenter.New([]dedenter.Rule{
    {Field: "msg", CollapseInner: true},
})
out := d.Apply(`{"msg":"   hello   world"}`)
// out → {"msg":"hello world"}
```

## Streaming

Use `dedenter.Stream(r, w, d)` to process a newline-delimited JSON stream
from an `io.Reader` and write results to an `io.Writer`. Empty lines are
skipped; lines that cannot be parsed as JSON are forwarded unchanged.
