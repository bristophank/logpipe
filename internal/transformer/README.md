# transformer

The `transformer` package applies field-level mutations to structured (JSON) log lines.

## Supported Operations

| Op        | Description                              |
|-----------|------------------------------------------|
| `set`     | Set a field to a static value            |
| `delete`  | Remove a field from the log line         |
| `rename`  | Rename a field, preserving its value     |
| `uppercase` | Convert a field's value to upper case  |
| `lowercase` | Convert a field's value to lower case  |

## Example Config

```json
{
  "transformers": [
    {"field": "level", "op": "uppercase"},
    {"field": "password", "op": "delete"},
    {"field": "msg", "op": "rename", "value": "message"},
    {"field": "env", "op": "set", "value": "production"}
  ]
}
```

Lines that are not valid JSON are passed through unchanged.
