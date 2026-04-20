# replay

The `replay` package re-emits stored log lines to any `io.Writer` at a
controlled rate. It is useful for testing pipelines against real-world
log files without flooding downstream sinks.

## Usage

```go
cfg := replay.Config{
    Rate: 500,  // 500 lines/sec; 0 = unlimited
    Loop: false,
}

// Convenience helper – opens the file for you.
if err := replay.Stream(cfg, "/var/log/app.log", os.Stdout); err != nil {
    log.Fatal(err)
}
```

## Config

| Field  | Type   | Description                                      |
|--------|--------|--------------------------------------------------|
| `Rate` | `int`  | Lines per second to emit. `0` means no throttle. |
| `Loop` | `bool` | Restart from the beginning when EOF is reached.  |

## Notes

- Empty lines are silently skipped.
- `Loop: true` requires the source to implement `io.ReadSeeker` (e.g. `*os.File`).
- The replayer writes each line terminated with `\n` regardless of the
  original line endings in the source file.
