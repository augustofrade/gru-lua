# Http

> Module path: `gru.http`

HTTP request helpers with structured request options and response objects.

Run `gru modules http` for documentation through the CLI.

## Summary

- All HTTP methods share the same function signature.
- Request bodies require a `Content-Type` header.
- Responses include headers, status and body helpers.

## Aliases

### GruHttpHeaders

`table<string, string>`

HTTP headers of a request or response.

### GruHttpRequestBody

`any`

Body of an HTTP request.

## Custom Types

### GruHttpRequestOptions

Options for an HTTP request.

Properties:

- `headers: GruHttpHeaders?`: Headers of the HTTP request.
- `body: GruHttpRequestBody?`: Body of the HTTP request.

### GruHttpResponseBody

Body helpers of an HTTP response.

Properties:

- `raw: fun(): string`: Returns the raw body as a string.
- `json: fun(): any, GruError`: Parses the body as JSON and returns a Lua value.

### GruHttpResponse

Response of an HTTP request.

Properties:

- `headers: GruHttpHeaders`: Headers of the HTTP response.
- `body: GruHttpResponseBody`: Body helpers of the HTTP response.
- `status: number`: HTTP status code.

## Functions

All methods have the same signature:

`<method>(url: string, options: GruHttpRequestOptions?) -> GruHttpResponse, GruError`

Available methods:

- `get`
- `post`
- `put`
- `patch`
- `delete`
- `head`
- `options`

Errors:

- The URL is invalid.
- `options.headers` is not a string-to-string table.
- A request body is provided without a `Content-Type` header.
- The body type is incompatible with the declared `Content-Type`.
- The request fails due to network, transport or response read errors.

Example:

```lua
local resp, err = gru.http.post("https://example.com/api", {
  headers = {
    ["Content-Type"] = "application/json",
  },
  body = {
    name = "Gru",
  }
})

if err then
  print(err)
  return
end

print(resp.status)
print(resp.body:raw())
```

### Response body helpers

`resp.body:raw()` returns the response body as a string.

`resp.body:json()` parses the body as JSON and returns `value, err`.

Example:

```lua
local data, err = resp.body:json()
if err then
  print(err)
  return
end

print(data.name)
```
