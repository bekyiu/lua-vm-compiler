package state

// 不同类型的lua值
type luaValue interface {
}
// 获取值的类型
//func typeOf(val luaValue) LuaType {
//	switch val.(type) {
//	case nil:
//		return api.LUA_TNIL
//	case bool:
//		return api.LUA_TBOOLEAN
//	case int64, float64:
//		return api.LUA_TNUMBER
//	case string:
//		return api.LUA_TSTRING
//	default:
//		panic("todo, 未知的值类型")
//	}
//}
