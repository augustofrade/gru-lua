# Assert

> Module path: `gru.assert`

Assertion helpers that raise Lua errors when validations fail.

Run `gru modules assert` for documentation through the CLI.

## Summary

- Assertion failures raise Lua errors.
- Most functions return the validated value so they can be used inline.
- The optional `label` parameter customizes the error message.

## Functions

### type(value: any, expected_type: string, label: string?) -> any

Asserts that a value matches the expected Lua type.

Example:

```lua
gru.assert.type(username, "string", "username")
```

### truthy(value: any, label: string?) -> any

Asserts that a value is truthy in Lua.

Example:

```lua
gru.assert.truthy(config.enabled, "config.enabled")
```

### equals(value: any, expected: any, label: string?) -> any

Asserts that two Lua values are equal.

Example:

```lua
gru.assert.equals(response.status, 200, "response.status")
```

### between(value: number, minimum: number, maximum: number, label: string?) -> number

Asserts that a number is between minimum and maximum, inclusive.

Example:

```lua
gru.assert.between(age, 18, 65, "age")
```

### is_string(value: any, label: string?) -> any

Asserts that a value is a string.

Example:

```lua
gru.assert.is_string(name, "name")
```

### is_number(value: any, label: string?) -> number

Asserts that a value is a number.

Example:

```lua
gru.assert.is_number(timeout, "timeout")
```

### is_boolean(value: any, label: string?) -> any

Asserts that a value is a boolean.

Example:

```lua
gru.assert.is_boolean(debug, "debug")
```

### is_table(value: any, label: string?) -> any

Asserts that a value is a table.

Example:

```lua
gru.assert.is_table(options, "options")
```

### is_array(value: any, label: string?) -> any

Asserts that a value is an array-like Lua table.

Example:

```lua
gru.assert.is_array(users, "users")
```

### is_integer(value: any, label: string?) -> number

Asserts that a value is an integer number.

Example:

```lua
gru.assert.is_integer(port, "port")
```

### is_function(value: any, label: string?) -> any

Asserts that a value is a function.

Example:

```lua
gru.assert.is_function(callback, "callback")
```

### greater_than(value: number, minimum: number, label: string?) -> number

Asserts that a number is greater than another number.

Example:

```lua
gru.assert.greater_than(retries, 0, "retries")
```

### lesser_than(value: number, maximum: number, label: string?) -> number

Asserts that a number is less than another number.

Example:

```lua
gru.assert.lesser_than(progress, 100, "progress")
```

### not_nil(value: any, label: string?) -> any

Asserts that a value is not nil.

Example:

```lua
gru.assert.not_nil(token, "token")
```

### not_empty(value: any, label: string?) -> any

Asserts that a string or table is not empty.

Example:

```lua
gru.assert.not_empty(items, "items")
```

### has_keys(value: table, keys: table, label: string?) -> any

Asserts that a table contains all provided string keys.

Example:

```lua
gru.assert.has_keys(user, { "id", "name", "email" }, "user")
```
