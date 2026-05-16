package gru

import (
	"time"

	"github.com/Shopify/go-lua"
)

func NewTimeModule() GruModule {
	module := NewModule("time", "Time related operations")
	module.FunctionBuilder("sleep", "Sleeps for <duration> seconds", timeSleep).
		NumberParam("seconds", "Duration in seconds").
		Register()

	return module
}

func timeSleep(l *lua.State) int {
	secs, valid := l.ToNumber(1)
	if !valid {
		return LuaVoidError("Expected number string or string")
	}

	time.Sleep(time.Duration(secs) * time.Second)

	l.PushNil()
	return 1
}
