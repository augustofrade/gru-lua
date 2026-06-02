package gru

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Shopify/go-lua"
	"github.com/augustofrade/gru-lua/gru/definitions"
	"github.com/augustofrade/gru-lua/gru/internal/luautil"
)

// TODO: use this alias
type GruHttpHeaders map[string]string

func NewHttpModule() definitions.GruModule {
	module := definitions.NewModule("http", "HTTP operations")

	module.HasCustomAlias("GruHttpHeaders", "HTTP headers of a response or request.", "table<string, string>")

	module.HasCustomType("GruHttpResponse", "Response of a HTTP request").
		Prop("headers", "GruHttpHeaders", "Headers of the HTTP response").
		Prop("body", "GruHttpBody", "Body of the HTTP response.")

	module.HasCustomType("GruHttpBody", "Body of a HTTP response").
		Prop("raw", "fun(): string, GruError", "Returns the raw body as a string.").
		Prop("json", "fun(): any, GruError", "Parses the body as JSON and returns a table.")

	module.FunctionBuilder("get", "Does a GET request at url", httpGet).
		StringParam("url", "URL of the HTTP request").
		ReturnsWithError("GruHttpResponse").
		Register()

	return module
}

// pushes a Lua table onto the stack with :raw() and :json() methods
func pushBodyTable(l *lua.State, raw []byte) {
	l.CreateTable(0, 2)

	// :raw() method
	l.PushGoFunction(func(l *lua.State) int {
		l.PushString(string(raw))
		l.PushNil()
		return 2
	})
	// body[raw] = function() ... end
	l.SetField(-2, "raw")

	// :json() method
	l.PushGoFunction(func(l *lua.State) int {
		var result any
		if err := json.Unmarshal(raw, &result); err != nil {
			return luautil.ErrorResult(l, fmt.Sprintf("JSON parse error: %s", err.Error()))
		}
		luautil.PushValue(l, result)
		l.PushNil()
		return 2
	})
	// body[json] = function() ... end
	l.SetField(-2, "json")
}

func httpGet(l *lua.State) int {
	if !luautil.IsString(l, 1) {
		return luautil.ErrorResult(l, "Expected string for 'url' parameter.")
	}

	url, _ := l.ToString(1)

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return luautil.ErrorResult(l, fmt.Sprintf("HTTP request error: %s", err.Error()))
	}

	resp, err := client.Do(req)
	if err != nil {
		return luautil.ErrorResult(l, fmt.Sprintf("HTTP request error: %s", err.Error()))
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return luautil.ErrorResult(l, fmt.Sprintf("HTTP read error: %s", err.Error()))
	}

	// Build the response table: { headers = {...}, body = { string=fn, json=fn } }
	l.CreateTable(0, 2)
	// headers
	luautil.PushTable(l, getHeaders(resp.Header))
	l.SetField(-2, "headers")

	// body
	pushBodyTable(l, raw)
	l.SetField(-2, "body")

	l.PushNil()
	return 2
}

func getHeaders(httpHeader http.Header) map[string]any {
	headers := make(map[string]any, len(httpHeader))
	for key, values := range httpHeader {
		headers[key] = values[0]
	}
	return headers

}
