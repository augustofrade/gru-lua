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

// Gets table values. Expect all values to be of the string type.
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

// Checks whether the table at index is an array table with consistent numerical indexes.
// Throws a runtime error if encounters an absent index
//
// { "orange", "apple", 1, true } -> true
//
// { [1] = 2, [2] = 3, [3] = 4 } -> true
//
// { [1] = 2, [3] = 4, [5] = 6 } -> false
//
// { 1, 2, nil, 3 } -> runtime error
func IsArrayTable(l *lua.State, index int) bool {
	length := GetTableLength(l, index)
	if length == 0 {
		return false
	}

	count := 0

	// Prevents errors from relative indexes (ex: -1) with the following required PushNil()
	absIndex := l.AbsIndex(index)

	l.PushNil()
	for l.Next(absIndex) {
		if l.TypeOf(-2) != lua.TypeNumber {
			// key is not is a number
			l.Pop(1)
			return false
		}

		key, _ := l.ToInteger(-2)
		if key < 1 || key > length {
			l.Pop(1)
			return false
		}

		count++
		l.Pop(1)
	}

	return count == length
}

// Returns table at index either as []string for "arrays" or map[string]any for keyed tables
func LuaTableToGo(l *lua.State, index int) any {
	if IsArrayTable(l, index) {
		length := GetTableLength(l, index)
		result := make([]any, 0, length)
		fmt.Println("a")

		for i := 1; i <= length; i++ {
			l.RawGetInt(index, i)
			result = append(result, LuaValueToGo(l, -1))
			l.Pop(1)
		}

		return result
	}

	// Prevents errors from relative indexes (ex: -1) with the following required PushNil()
	absIndex := l.AbsIndex(index)

	result := map[string]any{}
	l.PushNil()
	for l.Next(absIndex) {
		key, _ := l.ToString(-2)
		result[key] = LuaValueToGo(l, -1)
		l.Pop(1)
	}

	return result
}
