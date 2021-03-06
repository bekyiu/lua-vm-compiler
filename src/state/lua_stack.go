package state

import (
	"fmt"
	. "write_lua/src/api"
)

type luaStack struct {
	slots []luaValue // 存放值
	top   int        // 栈顶, 指向最顶层数据的高一个位置

	prev    *luaStack        // 指向上一个调用帧
	closure *closure         // 该调用帧所对应的函数
	varargs []luaValue       // 函数可变参数
	pc      int              // 当前函数Codes的索引
	state   *luaState        //
	openuvs map[int]*upvalue // 存放子闭包捕获当前闭包的局部变量, key是局部变量的寄存器索引
}

func newLuaStack(size int, state *luaState) *luaStack {
	return &luaStack{
		slots: make([]luaValue, size),
		top:   0,
		state: state,
	}
}

// 检查是否至少可以容纳n个值, 不满足则扩容至可以容纳n个
func (this *luaStack) check(n int) {
	free := len(this.slots) - this.top
	for i := 0; i < n-free; i++ {
		this.slots = append(this.slots, nil)
	}
}

func (this *luaStack) push(val luaValue) {
	if this.top == len(this.slots) {
		panic(fmt.Sprintf("when push value: %T lua stack overflow!\n", val))
	}
	this.slots[this.top] = val
	this.top++
}

func (this *luaStack) pop() luaValue {
	if this.top < 1 {
		panic("lua stack underflow!")
	}
	this.top--
	val := this.slots[this.top]
	this.slots[this.top] = nil
	return val
}

// 索引转换为绝对索引
func (this *luaStack) absIndex(idx int) int {
	// 说明是伪索引
	if idx >= 0 || idx <= LUA_REGISTRYINDEX {
		return idx
	}
	return idx + this.top + 1
}

// 判断索引是否有效, 栈的索引是从1/-1开始
func (this *luaStack) isValid(idx int) bool {
	// upvalue伪索引
	if idx < LUA_REGISTRYINDEX {
		// 从伪索引转换为真实索引, 从0开始
		uvIdx := LUA_REGISTRYINDEX - idx - 1
		c := this.closure
		return c != nil && uvIdx < len(c.upvals)
	}
	// 注册表伪索引
	if idx == LUA_REGISTRYINDEX {
		return true
	}
	absIdx := this.absIndex(idx)
	return absIdx > 0 && absIdx <= this.top
}

// 根据索引取值
func (this *luaStack) get(idx int) luaValue {
	if idx < LUA_REGISTRYINDEX {
		uvIdx := LUA_REGISTRYINDEX - idx - 1
		c := this.closure
		if c == nil || uvIdx >= len(c.upvals) {
			return nil
		}
		return *(c.upvals[uvIdx].val)
	}
	if idx == LUA_REGISTRYINDEX {
		return this.state.registry
	}
	absIdx := this.absIndex(idx)
	if absIdx > 0 && absIdx <= this.top {
		return this.slots[absIdx-1]
	}
	return nil
}

// 根据索引设置值
func (this *luaStack) set(idx int, val luaValue) {
	if idx < LUA_REGISTRYINDEX {
		uvIdx := LUA_REGISTRYINDEX - idx - 1
		c := this.closure
		if c != nil && uvIdx < len(c.upvals) {
			*(c.upvals[uvIdx].val) = val
		}
		return
	}
	if idx == LUA_REGISTRYINDEX {
		this.state.registry = val.(*luaTable)
		return
	}
	absIdx := this.absIndex(idx)
	if absIdx > 0 && absIdx <= this.top {
		this.slots[absIdx-1] = val
		return
	}
	panic(fmt.Sprintf("无效的索引: %d, 不能设置值: %T\n", idx, val))
}

// 逆序from到to元素的顺序
func (this *luaStack) reverse(from, to int) {
	slots := this.slots
	for from < to {
		slots[from], slots[to] = slots[to], slots[from]
		from++
		to--
	}
}

func (this *luaStack) popN(n int) []luaValue {
	vals := make([]luaValue, n)
	for i := n - 1; i >= 0; i-- {
		vals[i] = this.pop()
	}
	return vals
}

// push n个值, 多退少补
func (this *luaStack) pushN(vals []luaValue, n int) {
	nVals := len(vals)
	if n < 0 {
		n = nVals
	}
	for i := 0; i < n; i++ {
		if i < nVals {
			this.push(vals[i])
		} else {
			this.push(nil)
		}
	}
}
