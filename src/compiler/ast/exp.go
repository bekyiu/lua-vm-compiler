package ast

/*
exp ::=  nil | false | true | Numeral | LiteralString | ‘...’ | functiondef |
	 prefixexp | tableconstructor | exp binop exp | unop exp

prefixexp ::= var | functioncall | ‘(’ exp ‘)’
var ::=  Name | prefixexp ‘[’ exp ‘]’ | prefixexp ‘.’ Name
functioncall ::=  prefixexp args | prefixexp ‘:’ Name args
*/
type Exp interface{}

type NilExp struct{ Line int }    // nil
type TrueExp struct{ Line int }   // true
type FalseExp struct{ Line int }  // false
type VarargExp struct{ Line int } // ...

// Numeral
type IntegerExp struct {
	Line int
	Val  int64
}
type FloatExp struct {
	Line int
	Val  float64
}

// LiteralString
type StringExp struct {
	Line int
	Str  string
}

type NameExp struct {
	Line int
	Name string
}

// ==== 运算符表达式 =====
// unop exp
type UnopExp struct {
	Line int // line of operator
	Op   int // operator
	Exp  Exp
}

// exp1 op exp2
type BinopExp struct {
	Line int // line of operator
	Op   int // operator
	Exp1 Exp
	Exp2 Exp
}

type ConcatExp struct {
	Line int // line of last ..
	Exps []Exp
}


// ==== 函数定义表达式 =====
// functiondef ::= function funcbody
// funcbody ::= ‘(’ [parlist] ‘)’ block end
// parlist ::= namelist [‘,’ ‘...’] | ‘...’
// namelist ::= Name {‘,’ Name}
type FuncDefExp struct {
	Line     int
	LastLine int // line of `end`
	ParList  []string
	IsVararg bool
	Block    *Block
}