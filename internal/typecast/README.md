# typecast

The `typecast` package provides field-level type coercion for structured JSON log lines.

## Purpose

Log fields are often emitted as strings even when they carry numeric or boolean
semantics (e.g. `"status": "200"`, `"latency": "1.23"`). The `Caster` converts
those fields to their intended Go types so downstream consumers can perform
numeric comparisons, aggregations, or boolean checks without manual parsing.

## Usage

```go
rules := []typecast.Rule{
    {Field: "status",  To: "int"},
    {Field: "latency", To: "float"},
    {Field: "ok",      To: "bool"},
    {Field: "id",      To: "string"},
}

caster := typecast.New(rules)

out, err := caster.Apply(`{"status":"200","latency":"0.45","ok":"true"}`)
// out -> {"latency":0.45,"ok":true,"status":200}
```

## Supported Target Types

| `to`      | Description                                |
|-----------|--------------------------------------------|
| `string`  | Converts any value to its string form      |
| `int`     | Parses integer (also accepts float strings)|
| `float`   | Parses 64-bit floating-point number        |
| `bool`    | Parses `"true"` / `"false"`               |

## Behaviour

- Fields not present in the log line are silently skipped.
- Fields that cannot be cast to the target type are left unchanged.
- Invalid JSON input returns the original line along with an error.
