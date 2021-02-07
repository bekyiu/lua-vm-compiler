package state

import (
	"fmt"
	"math"
	. "write_lua/src/api"
	"write_lua/src/number"
)

// 定义lua算术/位运算符, 统一用两个参数, 一个返回值的函数来表示
var (
	iadd  = func(a, b int64) int64 { return a + b }
	fadd  = func(a, b float64) float64 { return a + b }
	isub  = func(a, b int64) int64 { return a - b }
	fsub  = func(a, b float64) float64 { return a - b }
	imul  = func(a, b int64) int64 { return a * b }
	fmul  = func(a, b float64) float64 { return a * b }
	imod  = number.IMod
	fmod  = number.FMod
	pow   = math.Pow
	div   = func(a, b float64) float64 { return a / b }
	iidiv = number.IFloorDiv
	fidiv = number.FFloorDiv
	band  = func(a, b int64) int64 { return a & b }
	bor   = func(a, b int64) int64 { return a | b }
	bxor  = func(a, b int64) int64 { return a ^ b }
	shl   = number.ShiftLeft
	shr   = number.ShiftRight
	iunm  = func(a, _ int64) int64 { return -a }
	funm  = func(a, _ float64) float64 { return -a }
	bnot  = func(a, _ int64) int64 { return ^a }
)

type operator struct {
	metamethod string
	integerFunc func(int64, int64) int64
	floatFunc   func(float64, float64) float64
}

var operators = []operator{
	operator{"__add", iadd, fadd},
	operator{"__sub", isub, fsub},
	operator{"__mul", imul, fmul},
	operator{"__mod", imod, fmod},
	operator{"__pow", nil, pow},
	operator{"__div", nil, div},
	operator{"__idiv", iidiv, fidiv},
	operator{"__band", band, nil},
	operator{"__bor", bor, nil},
	operator{"__bxor", bxor, nil},
	operator{"__shl", shl, nil},
	operator{"__shr", shr, nil},
	operator{"__unm", iunm, funm},
	operator{"__bnot", bnot, nil},
}

// 基于栈执行算术和按位运算
func (this *luaState) Arith(op ArithOp) {
	var a, b luaValue
	b = this.stack.pop()
	// 不是一元运算符
	if op != LUA_OPUNM && op != LUA_OPBNOT {
		a = this.stack.pop()
	} else {
		a = b
	}
	operator := operators[op]
	if result := _arith(a, b, operator); result != nil {
		this.stack.push(result)
		return
	}
	// 执行元方法
	mm := operator.metamethod
	if result, ok := callMetamethod(a, b, mm, this); ok {
		this.stack.push(result)
		return
	}
	panic(fmt.Sprintf("算术运算或位运算发生错误, 操作数a: %v, b: %v, 操作符op: %v\n",
		a, b, operator))
}

func _arith(a, b luaValue, op operator) luaValue {
	// 位运算, 则只能是整数参与
	if op.floatFunc == nil {
		if x, ok := convertToInteger(a); ok {
			if y, ok := convertToInteger(b); ok {
				return op.integerFunc(x, y)
			}
		}
	} else { // 算术运算
		if op.integerFunc != nil {
			// 如果操作数都是整数
			if x, ok := a.(int64); ok {
				if y, ok := b.(int64); ok {
					return op.integerFunc(x, y)
				}
			}
		}
		// 有一个操作数不是整数, 就全部提升为浮点数
		if x, ok := convertToFloat(a); ok {
			if y, ok := convertToFloat(b); ok {
				return op.floatFunc(x, y)
			}
		}
	}
	// 操作数不符合运算规定
	return nil
}