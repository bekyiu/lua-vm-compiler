package main

import (
	"fmt"
	"io/ioutil"
	. "write_lua/src/api"
	. "write_lua/src/compiler/lexer"
	"write_lua/src/state"
)

func main() {
	data, _ := ioutil.ReadFile("/Users/bekyiu/dev/luaCode/ch09/hello.lua")
	testLexer(string(data), "hello.lua")
}

func testLexer(chunk, chunkName string) {
	lexer := NewLexer(chunk, chunkName)
	for {
		line, kind, token := lexer.NextToken()
		fmt.Printf("[%2d] [%-10s] %s\n",
			line, kindToCategory(kind), token)
		if kind == TOKEN_EOF {
			break
		}
	}
}

func kindToCategory(kind int) string {
	switch {
	case kind < TOKEN_SEP_SEMI:
		return "other"
	case kind <= TOKEN_SEP_RCURLY:
		return "separator"
	case kind <= TOKEN_OP_NOT:
		return "operator"
	case kind <= TOKEN_KW_WHILE:
		return "keyword"
	case kind == TOKEN_IDENTIFIER:
		return "identifier"
	case kind == TOKEN_NUMBER:
		return "number"
	case kind == TOKEN_STRING:
		return "string"
	default:
		return "other"
	}
}

// ---------- vm 测试代码 ----------------
func testVM() {
	data, _ := ioutil.ReadFile("/Users/bekyiu/dev/luaCode/ch13/luac.out")
	//proto := binchunk.Undump(data)
	//luaMain(proto)
	ls := state.New()
	ls.Register("print", luaPrint)
	ls.Register("getmetatable", getMetatable)
	ls.Register("setmetatable", setMetatable)
	ls.Register("next", next)
	ls.Register("pairs", pairs)
	ls.Register("ipairs", iPairs)
	ls.Register("error", error)
	ls.Register("pcall", pCall)
	ls.Load(data, "luac.out", "b")
	ls.Call(0, 0)
}

func error(ls LuaState) int {
	// 错误对象已在栈顶
	return ls.Error()
}

func pCall(ls LuaState) int {
	nArgs := ls.GetTop() - 1
	status := ls.PCall(nArgs, -1, 0)
	// 此时栈顶是返回值, 或错误对象
	// 再压入一个bool表示函数是否执行出错
	ls.PushBoolean(status == LUA_OK)
	ls.Insert(1)
	return ls.GetTop()
}

func pairs(ls LuaState) int {
	// 迭代器
	ls.PushGoFunction(next)
	// 表
	ls.PushValue(1)
	// 初始key
	ls.PushNil()
	return 3
}

func iPairs(ls LuaState) int {
	ls.PushGoFunction(_iPairsAux)
	ls.PushValue(1)
	ls.PushInteger(0)
	return 3
}

// 返回下一组键值对
func _iPairsAux(ls LuaState) int {
	i := ls.ToInteger(2) + 1
	ls.PushInteger(i)
	if ls.GetI(1, i) == LUA_TNIL {
		return 1
	} else {
		return 2
	}
}

func next(ls LuaState) int {
	// 确保栈顶有两个值, 第一个是表, 第二个是key
	ls.SetTop(2)
	if ls.Next(1) {
		return 2
	} else {
		ls.PushNil()
		return 1
	}
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
