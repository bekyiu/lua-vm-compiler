package api

type LuaVM interface {
	LuaState
	PC() int             // 返回当前pc
	AddPC(n int)         // 修改pc
	Fetch() uint32       // 取指令
	GetConst(idx int)    // 将指定常量推入栈顶
	GetRK(rk int)        // 将指定常量或栈值推入栈顶
	RegisterCount() int  // 返回当前lua函数的寄存器数量
	LoadVararg(n int)    // 把传递给当前lua函数的变长参数推入栈顶 多退少补
	LoadProto(idx int)   // 把当前lua函数的子函数原型实例化为闭包推入栈顶
	CloseUpvalues(a int) // 闭合处于开启状态的upvalue
}
