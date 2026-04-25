# cutter

The `cutter` package slices string fields in JSON log lines by index range, similar to Python's string slicing.

## Rules

Each rule specifies:

| Field   | Type   | Description                                              |
|---------|--------|----------------------------------------------------------|
| `field` | string | The JSON field to cut                                    |
| `start` | int    | Start index (inclusive). Negative counts from end.       |
| `end`   | int    | End index (exclusive). `0` means to end of string.       |
| `as`    | string | Optional destination field. Defaults to source field.    |

## Example

```json
{
  "rules": [
    { "field": "request_id", "start": 0, "end": 8, "as": "short_id" },
    { "field": "message", "start": 7 }
  ]
}
```

Input:
```json
{"request_id": "abcd1234-5678-efgh", "message": "ERROR: disk full"}
```

Output:
```json
{"request_id": "abcd1234-5678-efgh", "short_id": "abcd1234", "message": "disk full"}
```

## Notes

- Non-string fields are silently skipped.
- Missing fields are silently skipped.
- Invalid JSON lines are passed through unchanged with an error returned.
