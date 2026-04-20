# dispatcher

Routes structured log lines to named sinks based on field-value rules.

## How it works

Each incoming JSON line is evaluated against an ordered list of `Rule` entries.
The first rule whose `field` matches the expected `value` wins and the line is
written to the named sink. If no rule matches, the line goes to the
`DefaultSink`. Lines are silently dropped when no default is configured.

Non-JSON lines bypass rule evaluation and go directly to the default sink.

## Config

```json
{
  "rules": [
    {"field": "level", "value": "error", "sink": "errors"},
    {"field": "level", "value": "warn",  "sink": "warnings"}
  ],
  "default_sink": "general"
}
```

## Usage

```go
d, err := dispatcher.New(cfg, sinks)
if err != nil { log.Fatal(err) }

// single line
sinkName, err := d.Dispatch(line)

// streaming
n, err := dispatcher.Stream(os.Stdin, d)
```

## Notes

- Rules are evaluated in order; first match wins.
- An unknown sink name in a rule or as the default causes `New` to return an error.
- Empty lines are ignored and do not count toward the dispatched total.
