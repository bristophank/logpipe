# classifier

The `classifier` package assigns a category label to structured log lines based on field matching rules.

## Rules

Each rule specifies:
- `field` — the JSON key to inspect
- `equals` — exact match against the field value
- `contains` — substring match against the field value
- `category` — the label to assign when the rule matches

The first matching rule wins. Lines that match no rule are passed through unchanged.

## Usage

```go
rules := []classifier.Rule{
    {Field: "level", Equals: "error",   Category: "critical"},
    {Field: "msg",   Contains: "timeout", Category: "network"},
}
c := classifier.New("category", rules)

out, err := c.Apply(`{"level":"error","msg":"disk full"}`)
// out: {"level":"error","msg":"disk full","category":"critical"}
```

## Streaming

```go
err := c.Stream(os.Stdin, os.Stdout)
```

Empty lines are skipped. Lines with invalid JSON are passed through unchanged.
