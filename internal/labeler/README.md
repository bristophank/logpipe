# labeler

The `labeler` package adds static key-value labels to every structured JSON log line passing through the pipeline.

## Use case

Use the labeler when you need to annotate all logs with fixed metadata such as environment, region, service name, or deployment version — without modifying the upstream log producer.

## Rules

Each rule has:

| Field   | Description                                |
|---------|--------------------------------------------|
| `key`   | JSON field name to add (or overwrite).     |
| `value` | Static string value to assign to the key.  |

Rules with an empty `key` are silently ignored.

## Behaviour

- Empty lines are passed through unchanged.
- Lines that are not valid JSON are returned as-is with an error (non-fatal; the pipeline continues).
- If a key already exists in the log line it is **overwritten** by the label value.

## Example config

```json
{
  "labels": [
    {"key": "env",     "value": "production"},
    {"key": "service", "value": "api-gateway"},
    {"key": "region",  "value": "eu-west-1"}
  ]
}
```
