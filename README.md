# logpipe

Lightweight CLI tool for filtering and routing structured log streams in real time.

---

## Installation

```bash
go install github.com/youruser/logpipe@latest
```

Or build from source:

```bash
git clone https://github.com/youruser/logpipe.git
cd logpipe
go build -o logpipe .
```

---

## Usage

Pipe any structured (JSON) log stream into `logpipe` and apply filters or route output to different destinations.

```bash
# Filter logs by level
./myapp | logpipe --level error

# Filter by a specific field value
./myapp | logpipe --match service=auth

# Route errors to a file, everything else to stdout
./myapp | logpipe --level error --out errors.log
```

### Flags

| Flag | Description |
|------|-------------|
| `--level` | Minimum log level to display (`debug`, `info`, `warn`, `error`) |
| `--match` | Filter by field key=value pair |
| `--out` | Write matching logs to a file instead of stdout |
| `--pretty` | Pretty-print JSON output |

### Example Output

```json
{"time":"2024-01-15T10:23:01Z","level":"error","service":"auth","msg":"token expired"}
```

---

## Requirements

- Go 1.21+
- Input must be newline-delimited JSON (NDJSON)

---

## License

MIT © 2024 youruser