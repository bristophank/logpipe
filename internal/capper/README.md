# capper

The `capper` package limits how many times a given field value may pass through the pipeline within a session. Once the cap is reached for a specific value, further lines carrying that value are dropped.

## Rules

Each rule specifies:

| Field  | Type   | Description                                      |
|--------|--------|--------------------------------------------------|
| field  | string | JSON field to inspect                            |
| cap    | int    | Maximum number of allowed occurrences            |

## Behaviour

- Lines with no matching field pass through unconditionally.
- Each unique value of `field` is tracked independently.
- Once a value's count exceeds `cap`, all subsequent lines with that value are dropped.
- Call `Reset()` to clear all counters.

## Example

```go
c := capper.New([]capper.Rule{
    {Field: "level", Cap: 3},
})

allowed := c.Allow(`{"level":"error","msg":"boom"}`)
```
