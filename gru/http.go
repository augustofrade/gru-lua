package gru

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Shopify/go-lua"
	"github.com/augustofrade/gru-lua/gru/definitions"
	"github.com/augustofrade/gru-lua/gru/internal/luautil"
)

type GruHttpHeaders map[string]string

type GruHttpRequest struct {
	httpReq *http.Request
	options map[string]any
	l       *lua.State
}

func NewHttpModule() definitions.GruModule {
	module := definitions.NewModule("http", "HTTP operations")

	module.HasCustomAlias("GruHttpHeaders", "HTTP headers of a response or request.", "table<string, string>")

	module.HasCustomAlias("GruHttpRequestBody", "Body of a HTTP request", "table<string, unknown>")

	module.HasCustomType("GruHttpRequestOptions", "Options for a HTTP request").
		Prop("headers", "GruHttpHeaders", "Headers of the HTTP request").
		Prop("body", "GruHttpRequestBody", "Body of the HTTP request.")

	module.HasCustomType("GruHttpResponseBody", "Body of a HTTP response").
		Prop("raw", "fun(): string", "Returns the raw body as a string.").
		Prop("json", "fun(): any, GruError", "Parses the body as JSON and returns a table.")

	module.HasCustomType("GruHttpResponse", "Response of a HTTP request").
		Prop("headers", "GruHttpHeaders", "Headers of the HTTP response").
		Prop("body", "GruHttpResponseBody", "Body of the HTTP response.").
		NumberProp("status", "Status code of the HTTP response")

	module.FunctionBuilder("get", "Does a GET request at url", httpGet).
		StringParam("url", "URL of the HTTP request").
		Param("options", "GruHttpRequestOptions?", "Settings of the request. The body property is ignored.").
		ReturnsWithError("GruHttpResponse").
		Register()

	module.FunctionBuilder("post", "Does a POST request at url", httpPost).
		StringParam("url", "URL of the HTTP request").
		Param("options", "GruHttpRequestOptions?", "Settings of the request.").
		ReturnsWithError("GruHttpResponse").
		Register()

	return module
}

func httpGet(l *lua.State) int {
	gruReq, err := newGruHttpRequest(l, http.MethodGet)
	if err != nil {
		return httpRequestErrorResult(l, err)
	}
	gruReq.SetRequestHeaders()

	return gruReq.DoRequest()
}

func httpPost(l *lua.State) int {
	gruReq, err := newGruHttpRequest(l, http.MethodGet)
	if err != nil {
		return httpRequestErrorResult(l, err)
	}
	gruReq.SetRequestHeaders()
	err = gruReq.SetRequestBody()
	if err != nil {
		return httpRequestErrorResult(l, err)
	}

	return gruReq.DoRequest()
}

func newGruHttpRequest(l *lua.State, httpMethod string) (*GruHttpRequest, error) {
	if !luautil.IsString(l, 1) {
		return nil, fmt.Errorf("Expected string for 'url' parameter.")
	}

	err := validateRequestOptionsTable(l)
	if err != nil {
		return nil, err
	}

	url, _ := l.ToString(1)

	httpReq, err := http.NewRequest(httpMethod, url, nil)

	return &GruHttpRequest{
			httpReq: httpReq,
			options: luautil.LuaTableToGo(l, 2).(map[string]any),
			l:       l},
		err
}

func validateRequestOptionsTable(l *lua.State) error {
	if l.IsNil(2) == false && l.IsTable(2) && luautil.IsArrayTable(l, 2) {
		return fmt.Errorf("Expected GruHttpRequestOptions for 'options' parameter.")
	}
	return nil
}

func (gruReq *GruHttpRequest) SetRequestHeaders() {
	headers, _ := gruReq.options["headers"].(map[string]any)
	for k, v := range headers {
		gruReq.httpReq.Header.Add(k, v.(string))
	}
}

func (gruReq *GruHttpRequest) SetRequestBody() error {
	body, _ := gruReq.options["body"].(map[string]any)

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	gruReq.httpReq.Body = io.NopCloser(bytes.NewReader(jsonBytes))
	gruReq.httpReq.ContentLength = int64(len(jsonBytes))
	return err
}

func (gruReq *GruHttpRequest) DoRequest() int {
	client := &http.Client{}

	resp, err := client.Do(gruReq.httpReq)
	if err != nil {
		return httpRequestErrorResult(gruReq.l, err)
	}

	err = buildResponseTable(gruReq.l, resp)
	if err != nil {
		return httpResponseErrorResult(gruReq.l, err)
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
		return 1
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

func httpRequestErrorResult(l *lua.State, err error) int {
	return luautil.ErrorResult(l, fmt.Sprintf("HTTP Request error: %s", err.Error()))
}

func httpResponseErrorResult(l *lua.State, err error) int {
	return luautil.ErrorResult(l, fmt.Sprintf("HTTP Response error: %s", err.Error()))
}
