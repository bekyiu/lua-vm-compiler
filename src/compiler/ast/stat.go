package ast

/*
stat ::=  ‘;’ |
	 varlist ‘=’ explist |
	 functioncall |
	 label |
	 break |
	 goto Name |
	 do block end |
	 while exp do block end |
	 repeat block until exp |
	 if exp then block {elseif exp then block} [else block] end |
	 for Name ‘=’ exp ‘,’ exp [‘,’ exp] do block end |
	 for namelist in explist do block end |
	 function funcname funcbody |
	 local function Name funcbody |
	 local namelist [‘=’ explist]
*/
type Stat interface{}

// === 简单语句 ===
// ‘;’
// 无任何作用
type EmptyStat struct{}

// break
type BreakStat struct {
	Line int // 记录行号
}

// ‘::’ Name ‘::’
type LabelStat struct {
	Name string // 标签名
}

// goto Name
type GotoStat struct {
	Name string // 标签名
}

// do block end
// 引入新的作用域
type DoStat struct {
	Block *Block
}

// functioncall
// 既可以是语句也可以是表达式
// type FuncCallStat = FuncCallExp

// ==== while和repeat用于实现同条件循环 ====
// while exp do block end
type WhileStat struct {
	Exp   Exp
	Block *Block
}

// repeat block until exp
type RepeatStat struct {
	Block *Block
	Exp   Exp
}

// 数值for循环
// for Name '=' exp ',' exp [',' exp] do block end
type ForNumStat struct {
	LineOfFor int // for的行号
	LineOfDo  int // do的行号
	VarName   string
	InitExp   Exp
	LimitExp  Exp
	StepExp   Exp
	Block     *Block
}

// 通用for循环
// for namelist in explist do block end
// namelist ::= Name {',' Name}
// explist ::= exp {',' exp}
type ForInStat struct {
	LineOfDo int
	NameList []string
	ExpList  []Exp
	Block    *Block
}

// ==== if ====
// if exp then block {elseif exp then block} end
type IfStat struct {
	// Exps[0] 对应于 Blocks[0] 对应于 if then的表达式和代码块
	Exps   []Exp
	Blocks []*Block
}

// ==== 局部变量声明语句 =====
// local namelist ['=' explist]
// namelist ::= Name {',' Name}
// explist ::= exp {',' exp}
type LocalVarDeclStat struct {
	LastLine int
	NameList []string
	ExpList  []Exp
}

// ==== 赋值语句 ====
// varlist ‘=’ explist
// varlist ::= var {‘,’ var}
// var ::=  Name | prefixexp ‘[’ exp ‘]’ | prefixexp ‘.’ Name
// explist ::= exp {',' exp}
type AssignStat struct {
	LastLine int
	VarList  []Exp
	ExpList  []Exp
}

// ==== 非局部函数定义语句, 会转换为赋值语句 =====
// function funcname funcbody
// funcname ::= Name {'.' Name} [':' Name]
// funcbody ::= '(' [parlist] ')' block end
// parlist ::= namelist [',' '...'] | '...'
// namelist ::= Name {',' Name}


// ==== 局部函数定义语句 ====
// local function Name funcbody
type LocalFuncDefStat struct {
	Name string
	Exp  *FuncDefExp
}
