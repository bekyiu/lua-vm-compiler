package state

import (
	"write_lua/src/api"
	"write_lua/src/number"
)

// 不同类型的lua值
type luaValue interface {
}

// 获取值的类型
func typeOf(val luaValue) api.LuaType {
	switch val.(type) {
	case nil:
		return api.LUA_TNIL
	case bool:
		return api.LUA_TBOOLEAN
	case int64, float64:
		return api.LUA_TNUMBER
	case string:
		return api.LUA_TSTRING
	case *luaTable:
		return api.LUA_TTABLE
	case *closure:
		return api.LUA_TFUNCTION
	default:
		panic("todo, 未知的值类型")
	}
}

// lua中, 只有false和nil表示假
func convertToBoolean(val luaValue) bool {
	switch x := val.(type) {
	case nil:
		return false
	case bool:
		return x
	default:
		return true
	}
}

func convertToFloat(val luaValue) (float64, bool) {
	switch x := val.(type) {
	case float64:
		return x, true
	case int64:
		return float64(x), true
	case string:
		return number.ParseFloat(x)
	default:
		return 0, false
	}
}

func convertToInteger(val luaValue) (int64, bool) {
	switch x := val.(type) {
	case int64:
		return x, true
	case float64:
		return number.FloatToInteger(x)
	case string:
		return _stringToInteger(x)
	default:
		return 0, false
	}
}

func _stringToInteger(str string) (int64, bool) {
	// 如果可以直接转换为整数
	if i, ok := number.ParseInteger(str); ok {
		return i, true
	}
	// 如果可以转换为浮点数
	if f, ok := number.ParseFloat(str); ok {
		return number.FloatToInteger(f)
	}
	return 0, false
}