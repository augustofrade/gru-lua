package gru

import (
	"time"

	"github.com/Shopify/go-lua"
	"github.com/augustofrade/gru-lua/gru/internal/luautil"
)

func NewTimeModule() GruModule {
	module := NewModule("time", "Time related operations")
	module.FunctionBuilder("sleep", "Sleeps for <duration> seconds", timeSleep).
		NumberParam("seconds", "Duration in seconds").
		Register()
	module.FunctionBuilder("unix", "Current Unix epoch time", timeUnix).
		ReturnsNumber().
		Register()
	return module
}

func timeSleep(l *lua.State) int {
	secs, valid := l.ToNumber(1)
	if !valid {
		return luautil.PushError(l, "Expected number or a string convertible to a number")
	}

	time.Sleep(time.Duration(secs) * time.Second)
	l.PushNil()
	return 1
}

func timeUnix(l *lua.State) int {

	l.PushInteger(int(time.Now().Unix()))
	return 1
}
