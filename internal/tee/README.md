# tee

The `tee` package fans out a single log stream to multiple named `io.Writer` sinks simultaneously.

## Usage

```go
te := tee.New()
te.Add("stdout", os.Stdout)
te.Add("file",   fileWriter)

// Write a single line to all sinks
te.Write([]byte(`{"level":"info","msg":"started"}` + "\n"))

// Or stream continuously from a reader
te.Stream(os.Stdin)
```

## API

| Method | Description |
|--------|-------------|
| `New()` | Create a new Tee |
| `Add(name, w)` | Register a named writer |
| `Remove(name)` | Unregister a writer |
| `Write(p)` | Fan out bytes to all writers |
| `Stream(r)` | Read lines from r and fan out |
| `Len()` | Number of registered writers |

Writers are safe to add/remove concurrently with ongoing writes.
