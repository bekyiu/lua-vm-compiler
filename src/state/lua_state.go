package state

// 封装了lua解释器的状态
type luaState struct {
	stack *luaStack
}

func New() *luaState {
	return &luaState{
		stack: newLuaStack(20),
	}
}

// 压入一个调用栈帧, this.stack 永远指向栈顶的栈帧, 是链表的头部
func (this *luaState) pushLuaStack(stack *luaStack) {
	stack.prev = this.stack
	this.stack = stack
}

// 弹出一个调用帧
func (this *luaState) popLuaStack() {
	stack := this.stack
	this.stack = stack.prev
	stack.prev = nil
}