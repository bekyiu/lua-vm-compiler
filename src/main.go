package main

import (
	"fmt"
	"io/ioutil"
	. "write_lua/src/api"
	"write_lua/src/state"
)

func main() {
	data, _ := ioutil.ReadFile("/Users/bekyiu/dev/luaCode/ch11/luac.out")
	//proto := binchunk.Undump(data)
	//luaMain(proto)
	ls := state.New()
	ls.Register("print", luaPrint)
	ls.Register("getmetatable", getMetatable)
	ls.Register("setmetatable", setMetatable)
	ls.Load(data, "luac.out", "b")
	ls.Call(0, 0)
}

func getMetatable(ls LuaState) int {
	// 栈顶有一个参数, 想要被获取元表的值
	if !ls.GetMetatable(1) {
		ls.PushNil()
	}
	return 1
}

func setMetatable(ls LuaState) int {
	// 栈顶两个参数
	// 1: 被设置元表的值
	// 2: 准备好的元表
	ls.SetMetatable(1)
	return 1
}

// go函数实现lua中的print
func luaPrint(ls LuaState) int {
	// lua给go传递的参数个数
	nArgs := ls.GetTop()
	for i := 1; i <= nArgs ; i++ {
		if ls.IsBoolean(i) {
			fmt.Printf("%t", ls.ToBoolean(i))
		} else if ls.IsString(i) {
			fmt.Print(ls.ToString(i))
		} else {
			fmt.Print(ls.TypeName(ls.Type(i)))
		}
		if i < nArgs {
			fmt.Print("\t")
		}
	}
	fmt.Println()
	return 0
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
