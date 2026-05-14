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
	module.Register("red", "Red color", redColor)
	module.Register("black", "Black color", blackColor)
	module.Register("green", "Green color", greenColor)
	module.Register("yellow", "Yellow color", yellowColor)
	module.Register("blue", "Blue color", blueColor)
	module.Register("magenta", "Magenta color", magentaColor)
	module.Register("cyan", "Cyan color", cyanColor)
	module.Register("white", "White color", whiteColor)
	module.Register("brightBlack", "Bright black color", brightBlackColor)
	module.Register("brightRed", "Bright red color", brightRedColor)
	module.Register("brightGreen", "Bright green color", brightGreenColor)
	module.Register("brightYellow", "Bright yellow color", brightYellowColor)
	module.Register("brightBlue", "Bright blue color", brightBlueColor)
	module.Register("brightMagenta", "Bright magenta color", brightMagentaColor)
	module.Register("brightCyan", "Bright cyan color", brightCyanColor)
	module.Register("brightWhite", "Bright white color", brightWhiteColor)
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
