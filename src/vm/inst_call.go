package vm

import . "write_lua/src/api"

// 把当前函数的子函数原型实例化为闭包, 放入指定寄存器
// R(A) := closure(KPROTO[Bx])
func closure(ins Instruction, vm LuaVM) {
	a, bx := ins.ABx()
	a += 1
	vm.LoadProto(bx)
	vm.Replace(a)
}

// R(A), ... ,R(A+C-2) := R(A)(R(A+1), ... ,R(A+B-1))
func call(ins Instruction, vm LuaVM) {
	a, b, c := ins.ABC()
	a += 1

	nArgs := _pushFuncAndArgs(a, b, vm)
	vm.Call(nArgs, c-1)
	// Call完之后栈顶都是返回值
	_popResults(a, c, vm)
}

// 压入函数和参数
func _pushFuncAndArgs(a, b int, vm LuaVM) (nArgs int) {
	// b - 1个参数需要传递
	if b >= 1 { // b == 1 压入函数
		vm.CheckStack(b)
		for i := a; i < a+b; i++ {
			vm.PushValue(i)
		}
		return b - 1
	} else { // 上一轮函数调用的返回值, 全部作为这一轮函数调用的参数
		_fixStack(a, vm)
		// 栈底 -> 栈顶
		// [寄存器..., 这一轮的函数, 这一轮函数的参数..., 上一轮调用的返回值]
		return vm.GetTop() - vm.RegisterCount() - 1
	}
}

// 把返回值放到指定寄存器
func _popResults(a, c int, vm LuaVM) {
	if c == 1 {
		// 无返回
	} else if c > 1 { // c-1个返回
		for i := a + c - 2; i >= a; i-- {
			vm.Replace(i)
		}
	} else { // 把被调函数的返回值全部返回
		vm.CheckStack(1)
		vm.PushInteger(int64(a))
	}
}

func _fixStack(a int, vm LuaVM) {
	x := int(vm.ToInteger(-1))
	vm.Pop(1)

	// x - a: 当前函数一个 + 上半部分参数的数量
	vm.CheckStack(x - a)
	// [a, x), 当前函数和上半部分参数压入到堆栈
	for i := a; i < x; i++ {
		vm.PushValue(i)
	}
	// 栈底 -> 栈顶
	// [寄存器..., 上一轮调用的返回值, 这一轮的函数, 这一轮函数的参数...]
	vm.Rotate(vm.RegisterCount()+1, x-a)

}

// return R(A), ..., R(A+B-2)
func _return(ins Instruction, vm LuaVM) {
	a, b, _ := ins.ABC()
	a += 1
	if b == 1 {
		// 无返回
	} else if b > 1 {
		vm.CheckStack(b - 1)
		for i := a; i <= a+b-2; i++ {
			vm.PushValue(i)
		}
	} else {
		// return 1, 2, f()
		// 此时f()返回值已在栈中, 把另一半也压栈, rotate
		_fixStack(a, vm)
	}
}

// R(A), R(A+1), ..., R(A+B-2) = vararg
func vararg(ins Instruction, vm LuaVM) {
	a, b, _ := ins.ABC()
	a += 1

	// b > 1 or b == 0
	// 把b - 1个vararg复制到寄存器, 或是把所有vararg复制到寄存器
	if b != 1 {
		vm.LoadVararg(b - 1)
		_popResults(a, b, vm)
	}
}

// return R(A)(R(A+1), ... ,R(A+B-1))
func tailCall(ins Instruction, vm LuaVM) {
	a, b, _ := ins.ABC()
	a += 1
	nArgs := _pushFuncAndArgs(a, b, vm)
	// 返回所有值
	vm.Call(nArgs, -1)
	_popResults(a, 0, vm)
}

// 把对象和方法拷贝到两个相邻的目标寄存器中
// R(A+1) := R(B); R(A) := R(B)[RK(C)]
func self(ins Instruction, vm LuaVM) {
	a, b, c := ins.ABC()
	a += 1
	b += 1
	// 对象
	vm.Copy(b, a+1)
	// 方法
	vm.GetRK(c)
	vm.GetTable(b)
	vm.Replace(a)
}

// R(A+3), ... , R(A+2+C) := R(A)(R(A+1), R(A+2))
func tForCall(ins Instruction, vm LuaVM) {
	a, _, c := ins.ABC()
	a += 1
	_pushFuncAndArgs(a, 3, vm)
	vm.Call(2, c)
	_popResults(a+3, c+1, vm)
}
