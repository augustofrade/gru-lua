package luautil

import (
	"fmt"
	"reflect"

	"github.com/Shopify/go-lua"
	"github.com/augustofrade/gru-lua/gru/internal/luamapper"
)

// Converts a Lua value at the specified index of the lua stack.
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

// Pushes value onto the lua stack
func PushValue(l *lua.State, value any) int {
	switch v := value.(type) {
	case string:
		l.PushString(v)
	case int:
		l.PushInteger(v)
	case int64:
		l.PushInteger(int(v))
	case float32:
		l.PushNumber(float64(v))
	case float64:
		l.PushNumber(v)
	case bool:
		l.PushBoolean(v)
	case []any:
		PushArrayTable(l, v)
	case map[string]any:
		PushTable(l, v)
	default:
		if reflect.TypeOf(v).Kind() == reflect.Struct {
			mapped := luamapper.MapStructToSlice(v)
			return PushValue(l, *mapped)
		}
		l.PushString(fmt.Sprint(v))
	}
	return 1
}

// Pushes a kvp table onto the lua stack
func PushTable(l *lua.State, kvp map[string]any) int {
	l.CreateTable(0, len(kvp))
	for key, value := range kvp {
		PushValue(l, value)
		// tbl[key] = value
		l.SetField(-2, key)
	}
	return 1
}

// Pushes an array table onto the lua stack.
func PushArrayTable[T any](l *lua.State, tbl []T) int {
	l.CreateTable(0, len(tbl))
	for i, value := range tbl {
		PushValue(l, value)
		// table.insert(tbl, i+1, value)
		l.RawSetInt(-2, i+1)
	}
	return 1
}

// Pushes an array table only made of strings onto the lua stack.
func PushStringArrayTable(l *lua.State, data *[]string) int {
	l.CreateTable(0, len(*data))
	for i, value := range *data {
		l.PushString(value)
		// table.insert(tbl, i+1, value)
		l.RawSetInt(-2, i+1)
	}
	return 1
}
