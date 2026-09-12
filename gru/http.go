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

type GruHttpHeaders map[string]string

type GruHttpRequestOptions struct {
	Body    []byte
	Headers GruHttpHeaders
}

func NewHttpModule() definitions.GruModule {
	module := definitions.NewModule("http", "HTTP operations")

	module.HasCustomAlias("GruHttpHeaders", "HTTP headers of a response or request.", "table<string, string>")

	module.HasCustomAlias("GruHttpRequestBody", "Body of a HTTP request", "table<string, unknown>")

	module.HasCustomType("GruHttpRequestOptions", "Options for a HTTP request").
		Prop("headers", "GruHttpHeaders", "Headers of the HTTP request").
		Prop("body", "GruHttpRequestBody", "Body of the HTTP request.")

	module.HasCustomType("GruHttpResponseBody", "Body of a HTTP response").
		Prop("raw", "fun(): string, GruError", "Returns the raw body as a string.").
		Prop("json", "fun(): any, GruError", "Parses the body as JSON and returns a table.")

	module.HasCustomType("GruHttpResponse", "Response of a HTTP request").
		Prop("headers", "GruHttpHeaders", "Headers of the HTTP response").
		Prop("body", "GruHttpBody", "Body of the HTTP response.").
		NumberProp("status", "Status code of the HTTP response")

	module.FunctionBuilder("get", "Does a GET request at url", httpGet).
		StringParam("url", "URL of the HTTP request").
		ReturnsWithError("GruHttpResponse").
		Register()

	return module
}

func httpGet(l *lua.State) int {
	req, err := newRequest(l, http.MethodGet)
	if err != nil {
		return httpErrorResult(l, err)
	}
	setRequestHeaders(l, req)

	return doRequestAndHandleResponse(l, req)
}

func newRequest(l *lua.State, httpMethod string) (*http.Request, error) {
	if !luautil.IsString(l, 1) {
		return nil, fmt.Errorf("Expected string for 'url' parameter.")
	}

	err := validateRequestOptionsTable(l)
	if err != nil {
		return nil, err
	}

	url, _ := l.ToString(1)

	req, err := http.NewRequest(httpMethod, url, nil)

	return req, err
}

func validateRequestOptionsTable(l *lua.State) error {
	if l.IsNil(2) == false && l.IsTable(2) && luautil.IsArrayTable(l, 2) {
		return fmt.Errorf("Expected keyed table for 'options' parameter.")
	}
	return nil
}

func setRequestHeaders(l *lua.State, req *http.Request) {
	options, _ := luautil.LuaTableToGo(l, 2).(map[string]any)

	headers, _ := options["headers"].(map[string]any)
	for k, v := range headers {
		req.Header.Add(k, v.(string))
	}
}

func doRequestAndHandleResponse(l *lua.State, req *http.Request) int {
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return luautil.ErrorResult(l, fmt.Sprintf("HTTP request error: %s", err.Error()))
	}

	err = buildResponseTable(l, resp)
	if err != nil {
		return luautil.ErrorResult(l, fmt.Sprintf("HTTP read error: %s", err.Error()))
	}
	return 2
}

// Converts the response headers into a GruHttpHeaders (map[string]string)
func getResponseHeaders(httpHeader http.Header) GruHttpHeaders {
	headers := make(GruHttpHeaders, len(httpHeader))
	for key, values := range httpHeader {
		headers[key] = values[0]
	}
	return headers
}

// Builds the response table:
//
//	{
//	  headers = {...},
//	  body = {
//	    string=fn,
//	    json=fn
//	  },
//	  status = number
//	}
func buildResponseTable(l *lua.State, resp *http.Response) error {
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	l.CreateTable(0, 2)
	// headers
	luautil.PushTable(l, getResponseHeaders(resp.Header))
	l.SetField(-2, "headers")

	// body
	pushResponseBodyTable(l, raw)
	l.SetField(-2, "body")

	// status
	l.PushNumber(float64(resp.StatusCode))
	l.SetField(-2, "status")

	l.PushNil()

	return nil
}

// Pushes a Lua table onto the stack with :raw() and :json() methods
func pushResponseBodyTable(l *lua.State, raw []byte) {
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

func httpErrorResult(l *lua.State, err error) int {
	return luautil.ErrorResult(l, fmt.Sprintf("HTTP Request error: %s", err.Error()))
}
