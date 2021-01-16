package main

import (
	"fmt"
	"io/ioutil"
	. "write_lua/src/api"
	"write_lua/src/state"
)

func main() {
	data, _ := ioutil.ReadFile("D:\\lua\\lua_code\\ch08\\luac.out")
	//proto := binchunk.Undump(data)
	//luaMain(proto)
	ls := state.New()
	ls.Load(data, "luac.out", "b")
	ls.Call(0, 0)
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
	fmt.Printf("\t栈顶: %d\n", top)
}
