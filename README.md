# Gru

Gru is a Lua runtime with useful modules and a simple CLI to run Lua scripts, made in Go.

It keeps regular Lua syntax and behavior and adds extra utilities under the global `gru` table,
all with documented type annotations.

```lua

-- Make a PATCH request with JSON body
local resp, err = gru.http.patch("https://jsonplaceholder.typicode.com/posts/1", {
  headers = {
    ["Content-Type"] = "application/json",
  },
  body = {
    title = "Gru's post",
  }
})

if err then
  print("Error:", err)
  return
end

-- Assert response status is as expected
gru.assert.equals(resp.status, 200)

-- Parse JSON
local json = resp.body:json()

-- Assert response JSON is as expected
gru.assert.not_empty(json, "Result JSON")

gru.assert.has_keys(json, { "userId", "id", "title", "body" }, "json")
gru.assert.equals(json.title, "Gru's post", "json.title")

-- Nice success message
print(gru.colors.green("PATCH request successful. Updated title:"), json.title)

-- Build the response file path with a nice API
local filename = "gru-" .. gru.time.unix() .. "-http-response.json"
local distPath = gru.path.join("./", "http", filename)

-- Directly save the JSON in a file
gru.json.dump(distPath, json)
```

## Why Gru

- Lua stays simple and familiar
- Useful modules are available out of the box
- Easy to extend with new Go-powered modules (GruModules and GruModuleFunctions)

## Current Modules

- `gru.assert`: Assertions and runtime validations
- `gru.colors`: Terminal string coloring helpers
- `gru.env`: Environment variable access and management
- `gru.fs`: File system operations
- `gru.http`: HTTP requests
- `gru.time`: Time and date utilities
- `gru.path`: Path manipulation functions
- `gru.json`: JSON parsing and serialization
- `gru.runtime`: Runtime metadata and error hooks
- ~~`gru.zip`: ZIP archive creation utilities~~ (TBD)

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
