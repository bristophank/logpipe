# parser

Parses raw log lines into `map[string]any` for downstream processing.

## Supported Formats

| Format   | Description                          |
|----------|--------------------------------------|
| `json`   | Standard JSON objects                |
| `logfmt` | `key=value` pairs (logfmt style)     |
| `auto`   | Detects JSON vs logfmt automatically |

## Usage

```go
p := parser.New(parser.FormatAuto)
m, err := p.Parse(`level=info msg="started"`)
// m == map[string]any{"level": "info", "msg": "started"}
```

## Error Handling

`Parse` returns a non-nil error in the following cases:

| Condition              | Error                          |
|------------------------|--------------------------------|
| Empty or blank line    | `ErrEmptyLine`                 |
| Invalid JSON           | wrapped `encoding/json` error  |
| Malformed logfmt input | wrapped logfmt parser error    |

Callers can check for expected empty-line skips using `errors.Is`:

```go
m, err := p.Parse(line)
if errors.Is(err, parser.ErrEmptyLine) {
    continue // skip blank lines
}
```

## Notes

- Logfmt bare keys (no `=`) are stored as `true`.
- Quoted values have surrounding `"` stripped.
- Empty lines return `ErrEmptyLine`.
