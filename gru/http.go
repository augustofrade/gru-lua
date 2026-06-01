package gru

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Shopify/go-lua"
	"github.com/augustofrade/gru-lua/gru/internal/luautil"
)

// TODO: use this alias
type GruHttpHeaders map[string]string

type GruHttpResponse struct {
	Headers map[string]any `lua:"headers"`
	Body    any            `lua:"body"`
}

func NewHttpModule() GruModule {
	module := NewModule("http", "HTTP operations")

	module.HasCustomAlias("GruHttpHeaders", "HTTP headers of a response or request.", "table<string, string>")

	module.HasCustomType("GruHttpResponse", "Response of a HTTP request").
		Prop("headers", "GruHttpHeaders", "Headers of the HTTP response").
		Prop("body", "any", "Body of the HTTP response. JSON responses are automatically parsed into tables.")

	module.FunctionBuilder("get", "Does a GET request at url", httpGet).
		StringParam("url", "URL of the HTTP request").
		ReturnsWithError("GruHttpResponse").
		Register()

	return module
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

	contentType := resp.Header.Get("Content-Type")

	gruResp := GruHttpResponse{
		Headers: make(map[string]any),
	}

	for key, values := range resp.Header {
		gruResp.Headers[key] = values[0]
	}
	if strings.HasPrefix(contentType, "application/json") {
		var respBody any
		if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
			return luautil.ErrorResult(l, fmt.Sprintf("HTTP response parsing error: %s", err.Error()))
		}

		gruResp.Body = respBody

		return luautil.PushValue(l, gruResp)
	}

	return luautil.PushValue(l, gruResp)
}
