# stencil

The `stencil` package injects new fields into JSON log lines by rendering
Go `text/template` expressions against the existing field values.

## Usage

```go
import "github.com/yourorg/logpipe/internal/stencil"

s, err := stencil.New([]stencil.Rule{
    {
        Target:    "summary",
        Template:  "[{{.level}}] {{.msg}}",
        Overwrite: false,
    },
})
if err != nil {
    log.Fatal(err)
}

out := s.Apply(`{"level":"error","msg":"disk full"}`)
// out => {"level":"error","msg":"disk full","summary":"[error] disk full"}
```

## Rules

| Field | Type | Description |
|-------|------|-------------|
| `target` | string | Destination field written into the log object |
| `template` | string | Go `text/template` expression; reference fields with `{{.fieldname}}` |
| `overwrite` | bool | When `true`, replaces an existing target field (default: `false`) |

## Notes

- Non-JSON lines are passed through unchanged.
- Template execution errors are silently skipped for that rule.
- Missing fields render as an empty string (zero value).
