# Time

> Module path: `gru.time`

Time-related helpers.

Run `gru modules time` for documentation through the CLI.

## Functions

### sleep(seconds: number) -> nil

Pauses the execution for the provided number of seconds.

Errors:

- The argument is not a number or a string convertible to a number.

Example:

```lua
gru.time.sleep(1)
```

### unix() -> number

Returns the current Unix epoch time in seconds.

Example:

```lua
local now = gru.time.unix()
```
