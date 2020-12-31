package api

type LuaVM interface {
	LuaState
	PC() int          // 返回当前pc
	AddPC(n int)      // 修改pc
	Fetch() uint32    // 取指令
	GetConst(idx int) // 将指定常量推入栈顶
	GetRK(rk int)     // 将指定常量或栈值推入栈顶
}
