# MKS template fs

mks's templates supports the following extensions:

(<span style="color:orange">\*</span>) are required

-   `main.load`
-   `main.unload`
-   `this.depends`
-   `this.prompts`
-   `this.goconfig`
-   `this.envconfig`
-   `...<folders>.<file>.extends`
-   `...<folders>.<file>.template` (<span style="color:orange">\*</span>)

## Extension descriptions

| files                            | description                                                                                  |               required               |
| -------------------------------- | -------------------------------------------------------------------------------------------- | :----------------------------------: |
| main.load                        | File code will be executed at the start of main function                                     |                                      |
| main.unload                      | File code will be executed at the end of main function                                       |                                      |
| this.prompts                     | Prompts to configure the template rightly                                                    |                                      |
| this.depends                     | Dependencies list of all templates that are required to run the current one                  |                                      |
| this.goconfig                    | Configurations that will be copied inside the Config struct in utils/config.go               |                                      |
| this.envconfig                   | Configurations that will be copied at the end of the app.env                                 |                                      |
| ...\<folders\>.\<file\>.extends  | Code that will extends the functions inside the specified file acording the file name format |                                      |
| ...\<folders\>.\<file\>.template | File/s that will contain the template code                                                   | <span style="color:orange">\*</span> |

#### `.extends file example:`

-   \<folder1\>.\<folder2\>.\<folderN\>.\<filename\>.extends => this file code will be copied at the bottom of the file: <span style="color:orange">src/</span>folder1/folder2/folderN/filename.go

#### `.template file example:`

-   \<folder1\>.\<folder2\>.\<folderN\>.\<filename\>.template => this file code will be copied in: <span style="color:orange">src/</span>folder1/folder2/folderN/filename.go

### Prompts structure (this.prompts)

```json
{
    "type": "replace",
    "default": "",
    "placeholder": "%PACKAGE_NAME%",
    "validate": validationType
}
```

#### `ValidationTypes:`

| Validation                                                 | Min Range   | Max Range  |    Type     |
| ---------------------------------------------------------- | ----------- | ---------- | :---------: |
| number                                                     | -2147483648 | 2147483647 |     int     |
| numberRange(min int, max int)                              | -2147483648 | 2147483647 |     int     |
| alphabet(caseSensitive bool, minLenght int, maxLenght int) | 0           | 65536      | string [AZ] |
| none                                                       |             |            |             |