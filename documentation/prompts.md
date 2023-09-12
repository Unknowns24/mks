### Prompt structure

Prompts are questions that MKS will ask the user every time the template is being installed

#### `Prompt types:`

| Type        | Description                                                                                       | Answer type           | Has validation |
| ----------- | ------------------------------------------------------------------------------------------------- | --------------------- | -------------- |
| **replace** | The answer of the prompt will be used to find and replace the specified placeholder on every file | mixed [String, Int]   | Yes            |
| **extend**  | The answer of the prompt will determine if use a specific extend file or not                      | confirmation [Yes/No] | No             |

-   `replace prompt structure`

```json
{
	"type": "replace",
	"prompt": "prompt question",
	"placeholder": "%%PACKAGE_NAME%%",
	"validate": "<validationType>"
}
```

-   `extends prompt structure`

```json
{
	"type": "extend",
	"prompt": "Do you require extra functions?",
	"extendFile": "<folderName>.<fileName>.extends"
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
