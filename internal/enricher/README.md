# enricher

Adds static or dynamic fields to structured (JSON) log lines.

## Rules

Each rule specifies:
- `field` — the key to insert or overwrite
- `value` — a literal string value (takes precedence over `source`)
- `source` — a dynamic source: `timestamp` (RFC3339 UTC) or `hostname`

## Example config

```json
{
  "enricher": [
    {"field": "env",  "value": "production"},
    {"field": "host", "source": "hostname"},
    {"field": "ts",   "source": "timestamp"}
  ]
}
```

## Behaviour

- Non-JSON lines are passed through unchanged.
- Rules are applied in order; later rules overwrite earlier ones for the same field.
