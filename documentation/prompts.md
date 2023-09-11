### Prompt structure

```json
{
	"type": "replace",
	"default": "",
	"placeholder": "%%PACKAGE_NAME%%",
	"validate": "<validationType>"
}
```

#### `ValidationTypes:`

| Validation                                                | Min Range   | Max Range  |      Type      |
| --------------------------------------------------------- | ----------- | ---------- | :------------: |
| number                                                    | -2147483648 | 2147483647 |      int       |
| numberRange(min int, max int)                             | -2147483648 | 2147483647 |      int       |
| alphabet(caseSensitive int, minLenght int, maxLenght int) | 0           | 65536      | string [az-AZ] |
| none                                                      |             |            |                |

**alphabet `caseSensitive` posible param values:**

| Value | Description        |
| ----- | ------------------ |
| 0     | Mayus only         |
| 1     | Minus only         |
| 2     | Not case sensitive |
