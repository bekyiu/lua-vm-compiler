package state

// access系列方法用于通过索引从栈里获取信息
import (
	"fmt"
	. "write_lua/src/api"
)

// 把给定的lua类型转换为字符串表示
func (this *luaState) TypeName(tp LuaType) string {
	switch tp {
	case LUA_TNONE:
		return "no value"
	case LUA_TNIL:
		return "nil"
	case LUA_TBOOLEAN:
		return "boolean"
	case LUA_TNUMBER:
		return "number"
	case LUA_TSTRING:
		return "string"
	case LUA_TTABLE:
		return "table"
	case LUA_TFUNCTION:
		return "function"
	case LUA_TTHREAD:
		return "thread"
	default:
		return "userdata"
	}
}

// 根据索引 返回类型
func (this *luaState) Type(idx int) LuaType {
	if this.stack.isValid(idx) {
		val := this.stack.get(idx)
		return typeOf(val)
	}
	return LUA_TNONE
}

func (this *luaState) IsNone(idx int) bool {
	return this.Type(idx) == LUA_TNONE
}
func (this *luaState) IsNil(idx int) bool {
	return this.Type(idx) == LUA_TNIL
}
func (this *luaState) IsNoneOrNil(idx int) bool {
	return this.IsNone(idx) || this.IsNil(idx)
}
func (this *luaState) IsBoolean(idx int) bool {
	return this.Type(idx) == LUA_TBOOLEAN
}
func (this *luaState) IsString(idx int) bool {
	t := this.Type(idx)
	return t == LUA_TSTRING || t == LUA_TNUMBER
}
func (this *luaState) IsNumber(idx int) bool {
	_, ok := this.ToNumberX(idx)
	return ok
}
func (this *luaState) IsInteger(idx int) bool {
	val := this.stack.get(idx)
	_, ok := val.(int64)
	return ok
}

// 从索引处取一个布尔值, 若不是 进行转换
func (this *luaState) ToBoolean(idx int) bool {
	val := this.stack.get(idx)
	return convertToBoolean(val)
}

// 将索引处的值取出并转化为数字类型, 若无法转换, 返回0
func (this *luaState) ToNumber(idx int) float64 {
	n, _ := this.ToNumberX(idx)
	return n
}
func (this *luaState) ToNumberX(idx int) (float64, bool) {
	val := this.stack.get(idx)
	return convertToFloat(val)
}

func (this *luaState) ToInteger(idx int) int64 {
	i, _ := this.ToIntegerX(idx)
	return i
}
func (this *luaState) ToIntegerX(idx int) (int64, bool) {
	val := this.stack.get(idx)
	return convertToInteger(val)
}

func (this *luaState) ToStringX(idx int) (string, bool) {
	val := this.stack.get(idx)
	switch x := val.(type) {
	case string:
		return x, true
	case int64, float64:
		s := fmt.Sprintf("%v", x)
		// 修改堆栈
		this.stack.set(idx, s)
		return s, true
	default:
		return "", false
	}
}

func (this *luaState) ToString(idx int) string {
	s, _ := this.ToStringX(idx)
	return s
}

// 判断索引处的值是否为go函数
func (this *luaState) IsGoFunction(idx int) bool {
	val := this.stack.get(idx)
	if c, ok := val.(*closure); ok {
		return c.goFunc != nil
	}
	return false
}

// 取出索引处的值, 转换为go函数返回
func (this *luaState) ToGoFunction(idx int) GoFunction {
	val := this.stack.get(idx)
	if c, ok := val.(*closure); ok {
		return c.goFunc
	}
	return nil
}

func (this *luaState) RawLen(idx int) uint {
	val := this.stack.get(idx)
	switch x := val.(type) {
	case string:
		return uint(len(x))
	case *luaTable:
		return uint(x.len())
	default:
		return 0
	}
}