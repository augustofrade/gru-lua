package luautil

import (
	"github.com/Shopify/go-lua"
)

func LuaValueToGo(l *lua.State, index int) any {
	switch l.TypeOf(index) {
	case lua.TypeString:
		val, _ := l.ToString(index)
		if val == "nil" {
			return nil
		}
		return val
	case lua.TypeNumber:
		val, _ := l.ToNumber(index)
		return val
	case lua.TypeBoolean:
		return l.ToBoolean(index)
	case lua.TypeNil:
		return nil
	case lua.TypeTable:
		return LuaTableToGo(l, index)
	default:
		return nil
	}
}
