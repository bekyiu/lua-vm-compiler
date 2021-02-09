package state

import (
	"fmt"
	. "write_lua/src/api"
)

func (this *luaState) RawEqual(idx1, idx2 int) bool {
	if !this.stack.isValid(idx1) || !this.stack.isValid(idx2) {
		return false
	}

	a := this.stack.get(idx1)
	b := this.stack.get(idx2)
	return _eq(a, b, nil)
}

// 比较索引处的两个值, 不改变栈的状态
// stack[idx1] op stack[idx2]
func (this *luaState) Compare(idx1, idx2 int, op CompareOp) bool {
	a := this.stack.get(idx1)
	b := this.stack.get(idx2)

	switch op {
	case LUA_OPEQ:
		return _eq(a, b, this)
	case LUA_OPLT:
		return _lt(a, b, this)
	case LUA_OPLE:
		return _le(a, b, this)
	default:
		panic(fmt.Sprintf("未知的运算符: %v", op))
	}
}

// 小于操作经仅对数字和字符串有意义
func _lt(a luaValue, b luaValue, ls *luaState) bool {
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x < y
		}
	case int64:
		switch y := b.(type) {
		case int64:
			return x < y
		case float64:
			return float64(x) < y
		}
	case float64:
		switch y := b.(type) {
		case int64:
			return x < float64(y)
		case float64:
			return x < y
		}
	}
	if result, ok := callMetamethod(a, b, "__lt", ls); ok {
		return convertToBoolean(result)
	} else {
		panic("_lt error!")
	}
}

func _le(a luaValue, b luaValue, ls *luaState) bool {
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x <= y
		}
	case int64:
		switch y := b.(type) {
		case int64:
			return x <= y
		case float64:
			return float64(x) <= y
		}
	case float64:
		switch y := b.(type) {
		case int64:
			return x <= float64(y)
		case float64:
			return x <= y
		}
	}
	if result, ok := callMetamethod(a, b, "__le", ls); ok {
		return convertToBoolean(result)
	} else if result, ok := callMetamethod(b, a, "__lt", ls); ok {
		return !convertToBoolean(result)
	} else {
		panic("_le error!")
	}
}
// 只有两个操作数在lua语言层面具有相同类型时, 等于运算才有可能返回true
func _eq(a luaValue, b luaValue, ls *luaState) bool {
	switch x := a.(type) {
	case nil:
		return b == nil
	case bool:
		y, ok := b.(bool)
		return ok && x == y
	case string:
		y, ok := b.(string)
		return ok && x == y
	case int64:
		switch y := b.(type) {
		case int64:
			return x == y
		case float64:
			return float64(x) == y
		default:
			return false
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x == y
		case int64:
			return x == float64(y)
		default:
			return false
		}
	case *luaTable:
		if y, ok := b.(*luaTable); ok && x != y && ls != nil {
			if result, ok := callMetamethod(x, y, "__eq", ls); ok {
				return convertToBoolean(result)
			}
		}
		return a == b
	default:
		return a == b
	}
}

