package main

import (
	"fmt"
	. "write_lua/src/api"
	"write_lua/src/binchunk"
	"write_lua/src/state"
	"write_lua/src/vm"
)

func main() {
	arr := make([]int, 0, 0)
	//arr[0] = 1
	zz := append(arr, 10)
	fmt.Println(zz)
	fmt.Println(zz[0])
	fmt.Println(cap(zz))
	fmt.Println(len(zz))
	//data, _ := ioutil.ReadFile("D:\\lua\\lua_code\\ch06\\luac.out")
	//proto := binchunk.Undump(data)
	//luaMain(proto)
}

func luaMain(proto *binchunk.Prototype) {
	nRegs := int(proto.MaxStackSize)
	ls := state.New(nRegs+8, proto)
	ls.SetTop(nRegs)
	for {
		pc := ls.PC()
		ins := vm.Instruction(ls.Fetch())
		if ins.Opcode() != vm.OP_RETURN {
			ins.Execute(ls)
			fmt.Printf("[%02d] %s ", pc+1, ins.OpName())
			printStack(ls)
		} else {
			break
		}
	}
}

func printStack(ls LuaState) {
	top := ls.GetTop()
	for i := 1; i <= top; i++ {
		t := ls.Type(i)
		switch t {
		case LUA_TBOOLEAN:
			fmt.Printf("[%t]", ls.ToBoolean(i))
		case LUA_TNUMBER:
			fmt.Printf("[%g]", ls.ToNumber(i))
		case LUA_TSTRING:
			fmt.Printf("[%q]", ls.ToString(i))
		default: // other values
			fmt.Printf("[%s]", ls.TypeName(t))
		}
	}
	fmt.Println()
}
