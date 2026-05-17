# Gru

Gru is an experimental Lua runtime built with Go.

It keeps regular Lua syntax and behavior and adds extra utilities under the global `gru` table,
all with documented type annotations.

```lua
local basename = gru.path.basename("~/Downloads/image.jpg")
local fullPath = gru.path.join("~/Documents", "assets", "sprites", basename)

print(gru.colors.light_blue(fullPath))
```

## Why Gru

- Lua stays simple and familiar
- Useful modules are available out of the box
- Easy to extend with new Go-powered modules (GruModules and GruModuleFunctions)

## Current Modules

- `gru.colors`
- `gru.time`
- `gru.path`
- `gru.json`

## CLI

Gru also provides a CLI under the same binary of the runtime.

- `gru init <path>`: initializes a git repository in the target path and generates Gru type files. Defaults to the current dir. **Requires git**.
- `gru types <path>`: generates Gru type annotations at the target path. Defaults to the current dir.
- `gru modules [module]`: lists all modules or details from one module
- `gru eval "<code>"`: evaluates Lua code directly from the terminal
- `gru help`: shows help information
- `gru <file.lua>`: runs a Lua file when no CLI command matches

## Build From Source

To build the project from source, run `make`.

### Dev Installation

To build and install the binary for the current user, run `make install-dev`.
The binary will be moved to `~/.local/bin` and available globally.
