package state

import "fmt"

type luaStack struct {
	slots []luaValue // 存放值
	top   int        // 栈顶, 指向最顶层数据的高一个位置
}

func newLuaStack(size int) *luaStack {
	return &luaStack{
		slots: make([]luaValue, size),
		top:   0,
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
	if idx >= 0 {
		return idx
	}
	return idx + this.top + 1
}

// 判断索引是否有效, 栈的索引是从1/-1开始
func (this *luaStack) isValid(idx int) bool {
	absIdx := this.absIndex(idx)
	return absIdx > 0 && absIdx <= this.top
}

// 根据索引取值
func (this *luaStack) get(idx int) luaValue {
	if this.isValid(idx) {
		return this.slots[idx-1]
	}
	return nil
}

// 根据索引设置值
func (this *luaStack) set(idx int, val luaValue) {
	if this.isValid(idx) {
		this.slots[idx-1] = val
		return
	}
	panic(fmt.Sprintf("无效的索引: %d, 不能设置值: %T\n", idx, val))
}
