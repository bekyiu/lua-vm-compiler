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
