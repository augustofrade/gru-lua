# Runtime

> Module path: `gru.runtime`

Runtime metadata and error hook helper.

Run `gru modules runtime` for documentation through the CLI.

## Functions

### version() -> string

Returns the current Gru runtime build string.

Example:

```lua
print(gru.runtime.version())
```

### on_error(callback: function) -> nil

Registers a callback to be triggered on runtime errors.

Errors:

- The argument is not a function.

Example:

```lua
gru.runtime.on_error(function()
  print("Program crashed")
end)
```
