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

## Notes

- Logfmt bare keys (no `=`) are stored as `true`.
- Quoted values have surrounding `"` stripped.
- Empty lines return an error.
