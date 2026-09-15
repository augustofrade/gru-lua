# Json

> Module path: `gru.json`

JSON parsing, serialization and file output helpers.

Run `gru modules json` for documentation through the CLI.

## Summary

- `stringify` converts Lua tables into JSON strings.
- `parse` converts JSON strings into Lua values.
- `dump` writes formatted JSON to disk and creates parent directories when needed.

## Functions

### stringify(table: table) -> string

Converts a Lua table to a JSON string.

Properties with the string value `"nil"` are converted to JSON `null`.

Errors:

- The argument is not a table.
- The table cannot be serialized to JSON.

Example:

```lua
local json = gru.json.stringify({ name = "Gru", active = true })
```

### parse(json: string) -> table

Parses a JSON string into a Lua value.

Errors:

- The argument is not a string.
- The input is not valid JSON.

Example:

```lua
local data = gru.json.parse('{"name":"Gru"}')
```

### dump(path: string, data: table) -> GruError

Serializes a table to indented JSON and writes it to a file.

Errors:

- The path is not a string.
- The data argument is not a table.
- The table cannot be serialized to JSON.
- The parent directory cannot be created.
- The process does not have permission to write the file.

Example:

```lua
local err = gru.json.dump("./dist/result.json", { ok = true })
if err then
  print(err)
end
```
