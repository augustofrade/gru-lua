# Colors Module

> Module path: `gru.colors`

Simple module with straightforward functions to color text in the terminal.

Run `gru modules colors` for documentation through the CLI.

## Functions

All functions have the same signature:

`<color>(text: string) -> string`

Available colors:

- black
- red
- green
- yellow
- blue
- magenta
- cyan
- white
- light_black
- light_red
- light_green
- light_yellow
- light_blue
- light_magenta
- light_cyan
- light_white

Example:

```lua
print(gru.colors.green("success"))
print(gru.colors.light_red("failure"))
```
