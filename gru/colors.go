package gru

import "github.com/Shopify/go-lua"

var colorCodes = map[string]string{
	"black":        "30",
	"red":          "31",
	"green":        "32",
	"yellow":       "33",
	"blue":         "34",
	"magenta":      "35",
	"cyan":         "36",
	"white":        "37",
	"lightBlack":   "90",
	"lightRed":     "91",
	"lightGreen":   "92",
	"lightYellow":  "93",
	"lightBlue":    "94",
	"lightMagenta": "95",
	"lightCyan":    "96",
	"lightWhite":   "97",
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
	module.FunctionBuilder("light_black", "Light black color", lightBlackColor).
		StringParam("text", "Text to be colored").
		Register()
	module.FunctionBuilder("light_red", "Light red color", lightRedColor).
		StringParam("text", "Text to be colored").
		Register()
	module.FunctionBuilder("light_green", "Light green color", lightGreenColor).
		StringParam("text", "Text to be colored").
		Register()
	module.FunctionBuilder("light_yellow", "Light yellow color", lightYellowColor).
		StringParam("text", "Text to be colored").
		Register()
	module.FunctionBuilder("light_blue", "Light blue color", lightBlueColor).
		StringParam("text", "Text to be colored").
		Register()
	module.FunctionBuilder("light_magenta", "Light magenta color", lightMagentaColor).
		StringParam("text", "Text to be colored").
		Register()
	module.FunctionBuilder("light_cyan", "Light cyan color", lightCyanColor).
		StringParam("text", "Text to be colored").
		Register()
	module.FunctionBuilder("light_white", "Light white color", lightWhiteColor).
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

func lightBlackColor(l *lua.State) int {
	return genericColor(l, colorCodes["lightBlack"])
}

func lightRedColor(l *lua.State) int {
	return genericColor(l, colorCodes["lightRed"])
}

func lightGreenColor(l *lua.State) int {
	return genericColor(l, colorCodes["lightGreen"])
}

func lightYellowColor(l *lua.State) int {
	return genericColor(l, colorCodes["lightYellow"])
}

func lightBlueColor(l *lua.State) int {
	return genericColor(l, colorCodes["lightBlue"])
}

func lightMagentaColor(l *lua.State) int {
	return genericColor(l, colorCodes["lightMagenta"])
}

func lightCyanColor(l *lua.State) int {
	return genericColor(l, colorCodes["lightCyan"])
}

func lightWhiteColor(l *lua.State) int {
	return genericColor(l, colorCodes["lightWhite"])
}
