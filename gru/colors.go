package gru

import "github.com/Shopify/go-lua"

var colorCodes = map[string]string{
	"black":         "30",
	"red":           "31",
	"green":         "32",
	"yellow":        "33",
	"blue":          "34",
	"magenta":       "35",
	"cyan":          "36",
	"white":         "37",
	"brightBlack":   "90",
	"brightRed":     "91",
	"brightGreen":   "92",
	"brightYellow":  "93",
	"brightBlue":    "94",
	"brightMagenta": "95",
	"brightCyan":    "96",
	"brightWhite":   "97",
}

func NewColorsModule() GruModule {
	module := NewModule("colors", "Write colored text int the terminal")
	module.FunctionBuilder("red", "Red color", redColor).
		StringParam("text", "Text to be colored").
		Register()
	module.FunctionBuilder("black", "Black color", blackColor).
		StringParam("text", "Text to be colored").
		Register()
	module.FunctionBuilder("green", "Green color", greenColor).
		StringParam("text", "Text to be colored").
		Register()
	module.FunctionBuilder("yellow", "Yellow color", yellowColor).
		StringParam("text", "Text to be colored").
		Register()
	module.FunctionBuilder("blue", "Blue color", blueColor).
		StringParam("text", "Text to be colored").
		Register()
	module.FunctionBuilder("magenta", "Magenta color", magentaColor).
		StringParam("text", "Text to be colored").
		Register()
	module.FunctionBuilder("cyan", "Cyan color", cyanColor).
		StringParam("text", "Text to be colored").
		Register()
	module.FunctionBuilder("white", "White color", whiteColor).
		StringParam("text", "Text to be colored").
		Register()
	module.FunctionBuilder("brightBlack", "Bright black color", brightBlackColor).
		StringParam("text", "Text to be colored").
		Register()
	module.FunctionBuilder("brightRed", "Bright red color", brightRedColor).
		StringParam("text", "Text to be colored").
		Register()
	module.FunctionBuilder("brightGreen", "Bright green color", brightGreenColor).
		StringParam("text", "Text to be colored").
		Register()
	module.FunctionBuilder("brightYellow", "Bright yellow color", brightYellowColor).
		StringParam("text", "Text to be colored").
		Register()
	module.FunctionBuilder("brightBlue", "Bright blue color", brightBlueColor).
		StringParam("text", "Text to be colored").
		Register()
	module.FunctionBuilder("brightMagenta", "Bright magenta color", brightMagentaColor).
		StringParam("text", "Text to be colored").
		Register()
	module.FunctionBuilder("brightCyan", "Bright cyan color", brightCyanColor).
		StringParam("text", "Text to be colored").
		Register()
	module.FunctionBuilder("brightWhite", "Bright white color", brightWhiteColor).
		StringParam("text", "Text to be colored").
		Register()
	return module
}

func genericColor(l *lua.State, color string) int {
	text, _ := l.ToString(1)

	l.PushString("\x1b[" + color + "m" + text + "\x1b[0m")
	return 1
}

func blackColor(l *lua.State) int {
	return genericColor(l, colorCodes["black"])
}

func redColor(l *lua.State) int {
	return genericColor(l, colorCodes["red"])
}

func greenColor(l *lua.State) int {
	return genericColor(l, colorCodes["green"])
}

func yellowColor(l *lua.State) int {
	return genericColor(l, colorCodes["yellow"])
}

func blueColor(l *lua.State) int {
	return genericColor(l, colorCodes["blue"])
}

func magentaColor(l *lua.State) int {
	return genericColor(l, colorCodes["magenta"])
}

func cyanColor(l *lua.State) int {
	return genericColor(l, colorCodes["cyan"])
}

func whiteColor(l *lua.State) int {
	return genericColor(l, colorCodes["white"])
}

func brightBlackColor(l *lua.State) int {
	return genericColor(l, colorCodes["brightBlack"])
}

func brightRedColor(l *lua.State) int {
	return genericColor(l, colorCodes["brightRed"])
}

func brightGreenColor(l *lua.State) int {
	return genericColor(l, colorCodes["brightGreen"])
}

func brightYellowColor(l *lua.State) int {
	return genericColor(l, colorCodes["brightYellow"])
}

func brightBlueColor(l *lua.State) int {
	return genericColor(l, colorCodes["brightBlue"])
}

func brightMagentaColor(l *lua.State) int {
	return genericColor(l, colorCodes["brightMagenta"])
}

func brightCyanColor(l *lua.State) int {
	return genericColor(l, colorCodes["brightCyan"])
}

func brightWhiteColor(l *lua.State) int {
	return genericColor(l, colorCodes["brightWhite"])
}
