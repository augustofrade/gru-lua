# Fs

> Module path: `gru.fs`

File system helpers for reading, writing and inspecting files and directories.

Run `gru modules fs` for documentation through the CLI.

## Summary

- Most file system operations return `nil, GruError` on failure.
- Paths may be relative or absolute unless stated otherwise.
- Some functions return structured values described in the custom types below.

## Custom Types

### GruDirentry

Represents one entry returned by `gru.fs.read_dir`.

Properties:

- `name: string`: Name of the directory entry.
- `is_dir: boolean`: Whether the entry is a directory.
- `parent_path: string`: Absolute parent path of the directory entry.

### GruFileInfo

Information about a file or directory.

Properties:

- `name: string`: Name of the file or directory.
- `fullpath: string`: Full path of the file or directory.
- `is_dir: boolean`: Whether the entry is a directory.
- `last_modification_time: number`: Unix timestamp of the last modification.
- `size: number`: Size of the file or directory.

## Functions

### read_dir(dir: string) -> GruDirentry[], GruError

Reads the provided directory path and returns its contents.

Errors:

- The path does not exist.
- The path is not a directory.
- The process does not have permission to read the directory.

Example:

```lua
local entries, err = gru.fs.read_dir("./src")
if err then
  print(err)
  return
end
```

### read_file(file_path: string) -> string, GruError

Reads the file in the provided path and returns its contents as a string.

Errors:

- The file does not exist.
- The path points to a directory instead of a file.
- The process does not have permission to read the file.

Example:

```lua
local content, err = gru.fs.read_file("./README.md")
if err then
  print(err)
  return
end
```

### write_file(file_path: string, data: string, permissions: number?) -> GruError

Writes data to the provided file path, creating it if necessary with the given permissions.

Errors:

- The destination path is invalid.
- The parent directory does not exist or is not writable.
- The process does not have permission to create or write the file.

Example:

```lua
local err = gru.fs.write_file("./dist/output.txt", "hello", 0644)
if err then
  print(err)
end
```

### create(file_path: string) -> GruError

Creates or truncates the named file and ensures the parent directories exist.

Errors:

- The file path is invalid or empty.
- The parent directory cannot be created.
- The process does not have permission to create or truncate the file.

Example:

```lua
local err = gru.fs.create("./tmp/data.txt")
if err then
  print(err)
end
```

### mkdir(path: string, permissions: number?) -> GruError

Creates a directory along with any necessary parents.

Errors:

- The path is invalid.
- The process does not have permission to create the directory.
- A non-directory path component prevents parent creation.

Example:

```lua
local err = gru.fs.mkdir("./tmp/cache", 0755)
if err then
  print(err)
end
```

### remove(path: string) -> GruError

Removes the named file or directory.

Use `remove_all` when removing a directory with children.

Errors:

- The path does not exist.
- The target is a non-empty directory.
- The process does not have permission to remove the target.

Example:

```lua
local err = gru.fs.remove("./tmp/data.txt")
if err then
  print(err)
end
```

### remove_all(path: string) -> GruError

Removes a path and any children it contains.

Errors:

- The process does not have permission to remove part of the target tree.
- A file system error prevents traversal or removal.

Example:

```lua
local err = gru.fs.remove_all("./tmp")
if err then
  print(err)
end
```

### exists(path: string) -> boolean, GruError

Returns whether the file or directory at the given path exists.

Returns `false, nil` when the path does not exist.

Errors:

- The existence check fails for a reason other than a missing path.
- The process does not have permission to inspect the target.

Example:

```lua
local exists, err = gru.fs.exists("./tmp/data.txt")
if err then
  print(err)
  return
end
```

### stat(path: string) -> GruFileInfo, GruError

Returns information about the file or directory at the given path.

Errors:

- The path does not exist.
- The process does not have permission to inspect the target.
- The file metadata cannot be read.

Example:

```lua
local info, err = gru.fs.stat("./README.md")
if err then
  print(err)
  return
end
```
