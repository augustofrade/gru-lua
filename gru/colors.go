package gru

import "github.com/Shopify/go-lua"

func NewColorsModule() GruModule {
	module := NewModule("colors", "Write colored text int the terminal")
	module.Register("red", "Red color", redColor)
	return module
}

func redColor(l *lua.State) int {
	text, _ := l.ToString(1)

	l.PushString("\x1b[31m" + text + "\x1b[0m")
	return 1
}
