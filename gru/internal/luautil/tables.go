package luautil

import (
	"fmt"

	"github.com/Shopify/go-lua"
)

// Transforms the map into a Lua table and pushes it onto the stack
func PushTable(l *lua.State, kvp map[string]any) {
	l.CreateTable(0, len(kvp))
	for key, value := range kvp {
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
			l.PushString(fmt.Sprint(v))
		}

		l.SetField(-2, key)
	}
}

func PushArrayTable(l *lua.State, kvp []any) {
	l.CreateTable(0, len(kvp))
	for i, value := range kvp {
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
			l.PushString(fmt.Sprint(v))
		}

		l.RawSetInt(-2, i+1)
	}
}

func PushStringArrayTable(l *lua.State, data *[]string) int {
	l.CreateTable(0, len(*data))
	for i, val := range *data {
		l.PushString(val)
		// table.insert(tbl, i+1, val)
		l.RawSetInt(-2, i+1)
	}
	return 1
}

// Gets the length of the table at index of stack
//
// ({ "a", "b" }, ...)
//
// GetTableLength(1) -> 2
func GetTableLength(l *lua.State, index int) int {
	l.Length(index)
	length, _ := l.ToInteger(-1)
	l.Pop(1)

	return length
}

func GetTableStringValues(l *lua.State) (*[]string, error) {
	values := make([]string, 0)

	i := 1
	for {
		l.RawGetInt(1, i)
		if l.IsNil(-1) {
			// end of the table
			l.Pop(1)
			break
		}

		if l.IsString(-1) {
			val, _ := l.ToString(-1)
			values = append(values, val)
		} else {
			return nil, fmt.Errorf("")
		}
		l.Pop(1)
		i++
	}

	return &values, nil
}
