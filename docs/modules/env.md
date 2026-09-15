# Env

> Module path: `gru.env`

Environment variable helpers.

Run `gru modules env` for documentation through the CLI.

## Summary

- Values are read from and written to the current process environment.
- Mutating functions may return an error string on failure.
- `lookup` distinguishes between an empty value and a missing variable.

## Functions

### set(key: string, value: string) -> string?

Sets an environment variable.

Returns `nil` on success or an error string on failure.

Example:

```lua
local err = gru.env.set("APP_ENV", "development")
if err then
  print(err)
end
```

### get(key: string) -> string

Returns the value of an environment variable.

If the variable is not set, an empty string is returned.

Example:

```lua
local home = gru.env.get("HOME")
```

### clear() -> nil

Deletes all environment variables from the current process.

Example:

```lua
gru.env.clear()
```

### unset(key: string) -> string?

Unsets an environment variable.

Returns `nil` on success or an error string on failure.

Example:

```lua
local err = gru.env.unset("APP_ENV")
if err then
  print(err)
end
```

### lookup(key: string) -> string, boolean

Looks up an environment variable and returns its value along with whether it exists.

Example:

```lua
local value, exists = gru.env.lookup("HOME")
```

### all() -> table

Returns all environment variables as an array-like table of `key=value` strings.

Example:

```lua
local entries = gru.env.all()
```
