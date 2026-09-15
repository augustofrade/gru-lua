# Path

> Module path: `gru.path`

Path manipulation helpers based on the host operating system.

Run `gru modules path` for documentation through the CLI.

## Summary

- Path separators and absolute-path rules follow the current OS.
- Most functions operate only on path strings and do not access the file system.
- `absolute` and `resolve` may depend on the current working directory.

## Custom Types

### GruPathInfo

Properties:

- `dir: string`: Full directory of the path.
- `file: string`: Filename with extension.
- `ext: string`: Extension of the path.

## Functions

### basename(path: string) -> string, GruError

Returns the last portion of the path.

Example:

```lua
local name = gru.path.basename("/tmp/file.txt")
```

### dirname(path: string) -> string, GruError

Returns the directory name of the path.

Example:

```lua
local dir = gru.path.dirname("/tmp/file.txt")
```

### extname(path: string) -> string, GruError

Returns the file extension of the path.

Example:

```lua
local ext = gru.path.extname("archive.tar.gz")
```

### is_absolute(path: string) -> boolean, GruError

Returns whether the path is absolute.

Example:

```lua
local abs = gru.path.is_absolute("/tmp/file.txt")
```

### is_relative(path: string) -> boolean, GruError

Returns whether the path is relative.

Example:

```lua
local rel = gru.path.is_relative("docs/readme.md")
```

### join(...: string) -> string, GruError

Joins path elements using the correct separator for the current OS.

Errors:

- No path parts are provided.
- Any argument is not a string.

Example:

```lua
local full = gru.path.join(".", "docs", "readme.md")
```

### parse(path: string) -> GruPathInfo, GruError

Parses a path into a GruPathInfo.

Example:

```lua
local parsed = gru.path.parse("/tmp/file.txt")
print(parsed.dir)
print(parsed.file)
print(parsed.ext)
```

### absolute(path: string) -> string, GruError

Returns an absolute representation of the path.

Errors:

- The argument is not a string.
- The absolute path cannot be resolved.

Example:

```lua
local full = gru.path.absolute("./README.md")
```

### clean(path: string) -> string, GruError

Normalizes a path by removing redundant separators and resolving `.` and `..` segments when possible.

Example:

```lua
local normalized = gru.path.clean("./tmp/../docs//readme.md")
```

### stem(path: string) -> string, GruError

Returns the file name without its extension.

Example:

```lua
local stem = gru.path.stem("archive.tar.gz")
```

### relative(base_path: string, target_path: string) -> string, GruError

Returns a relative path from `base_path` to `target_path`.

Errors:

- Any argument is not a string.
- A relative path cannot be computed between the two paths.

Example:

```lua
local rel = gru.path.relative("/tmp", "/tmp/docs/readme.md")
```

### resolve(...: string) -> string, GruError

Joins, cleans and converts the provided path parts into an absolute path.

Errors:

- Any argument is not a string.
- The resulting absolute path cannot be resolved.

Example:

```lua
local resolved = gru.path.resolve(".", "docs", "readme.md")
```
