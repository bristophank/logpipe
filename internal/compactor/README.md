# compactor

The `compactor` package removes fields with empty, null, or zero values from structured JSON log lines.

## Rules

Each rule targets either a specific field or all fields, and specifies which value types to drop:

| Field      | Type   | Description                              |
|------------|--------|------------------------------------------|
| `field`    | string | Target field name; empty means all fields |
| `drop_null`  | bool | Remove fields whose value is `null`      |
| `drop_empty` | bool | Remove fields whose value is `""`        |
| `drop_zero`  | bool | Remove fields whose numeric value is `0` |
| `drop_false` | bool | Remove fields whose boolean value is `false` |

## Example

```json
{
  "compactor": [
    { "drop_null": true, "drop_empty": true },
    { "field": "retries", "drop_zero": true }
  ]
}
```

Input:
```json
{"level":"info","msg":"","err":null,"retries":0,"host":"web-1"}
```

Output:
```json
{"level":"info","host":"web-1"}
```

## Behaviour

- Lines that are not valid JSON are passed through unchanged.
- Multiple rules are applied in order; each rule is independent.
- If `field` is empty the rule applies to every key in the object.
