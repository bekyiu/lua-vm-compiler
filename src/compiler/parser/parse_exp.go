package parser

import (
	. "write_lua/src/compiler/ast"
	. "write_lua/src/compiler/lexer"
	"write_lua/src/number"
)

// explist ::= exp {',' exp}
func parseExpList(lexer *Lexer) []Exp {
	exps := make([]Exp, 0, 4)
	exps = append(exps, parseExp(lexer))
	for lexer.LookAhead() == TOKEN_SEP_COMMA {
		lexer.NextToken()
		exps = append(exps, parseExp(lexer))
	}
	return exps
}

/*
运算符表达式:
exp   ::= exp12
exp12 ::= exp11 {or exp11}
exp11 ::= exp10 {and exp10}
exp10 ::= exp9 {(‘<’ | ‘>’ | ‘<=’ | ‘>=’ | ‘~=’ | ‘==’) exp9}
exp9  ::= exp8 {‘|’ exp8}
exp8  ::= exp7 {‘~’ exp7}
exp7  ::= exp6 {‘&’ exp6}
exp6  ::= exp5 {(‘<<’ | ‘>>’) exp5}
exp5  ::= exp4 {‘..’ exp4}
exp4  ::= exp3 {(‘+’ | ‘-’) exp3}
exp3  ::= exp2 {(‘*’ | ‘/’ | ‘//’ | ‘%’) exp2}
exp2  ::= {(‘not’ | ‘#’ | ‘-’ | ‘~’)} exp1
exp1  ::= exp0 {‘^’ exp2}
非运算符表达式:
exp0  ::= nil | false | true | Numeral | LiteralString
		| ‘...’ | functiondef | prefixexp | tableconstructor
*/
func parseExp(lexer *Lexer) Exp {
	return parseExp12(lexer)
}

// 左结合
// exp12 ::= exp11 {or exp11}
func parseExp12(lexer *Lexer) Exp {
	exp := parseExp11(lexer)
	// exp 不断地作为左子树
	for lexer.LookAhead() == TOKEN_OP_OR {
		line, op, _ := lexer.NextToken() // or
		exp = &BinopExp{
			Line: line,
			Op:   op,
			Exp1: exp,
			Exp2: parseExp11(lexer),
		}
	}
	return exp
}

// exp11 ::= exp10 {and exp10}
func parseExp11(lexer *Lexer) Exp {
	exp := parseExp10(lexer)
	// exp 不断地作为左子树
	for lexer.LookAhead() == TOKEN_OP_AND {
		line, op, _ := lexer.NextToken() // and
		exp = &BinopExp{
			Line: line,
			Op:   op,
			Exp1: exp,
			Exp2: parseExp10(lexer),
		}
	}
	return exp
}

// exp10 ::= exp9 {(‘<’ | ‘>’ | ‘<=’ | ‘>=’ | ‘~=’ | ‘==’) exp9}
func parseExp10(lexer *Lexer) Exp {
	exp := parseExp9(lexer)
	for {
		switch lexer.LookAhead() {
		case TOKEN_OP_LT, TOKEN_OP_GT, TOKEN_OP_NE,
			TOKEN_OP_LE, TOKEN_OP_GE, TOKEN_OP_EQ:
			line, op, _ := lexer.NextToken()
			exp = &BinopExp{Line: line, Op: op, Exp1: exp, Exp2: parseExp9(lexer)}
		default:
			return exp
		}
	}
}

// exp9  ::= exp8 {‘|’ exp8}
func parseExp9(lexer *Lexer) Exp {
	exp := parseExp8(lexer)
	for lexer.LookAhead() == TOKEN_OP_BOR {
		line, op, _ := lexer.NextToken()
		exp = &BinopExp{
			Line: line,
			Op:   op,
			Exp1: exp,
			Exp2: parseExp8(lexer)}
	}
	return exp
}

// exp8  ::= exp7 {‘~’ exp7}
func parseExp8(lexer *Lexer) Exp {
	exp := parseExp7(lexer)
	for lexer.LookAhead() == TOKEN_OP_BXOR {
		line, op, _ := lexer.NextToken()
		exp = &BinopExp{Line: line, Op: op, Exp1: exp, Exp2: parseExp7(lexer)}
	}
	return exp
}

// exp7  ::= exp6 {‘&’ exp6}
func parseExp7(lexer *Lexer) Exp {
	exp := parseExp6(lexer)
	for lexer.LookAhead() == TOKEN_OP_BAND {
		line, op, _ := lexer.NextToken()
		exp = &BinopExp{Line: line, Op: op, Exp1: exp, Exp2: parseExp6(lexer)}
	}
	return exp
}

// exp6  ::= exp5 {(‘<<’ | ‘>>’) exp5}
func parseExp6(lexer *Lexer) Exp {
	exp := parseExp5(lexer)
	for {
		switch lexer.LookAhead() {
		case TOKEN_OP_SHL, TOKEN_OP_SHR:
			line, op, _ := lexer.NextToken()
			exp = &BinopExp{Line: line, Op: op, Exp1: exp, Exp2: parseExp5(lexer)}
		default:
			return exp
		}
	}
}

// 右结合 但是特殊处理 生成多叉树
// exp5  ::= exp4 {‘..’ exp4}
func parseExp5(lexer *Lexer) Exp {
	exp := parseExp4(lexer)
	if lexer.LookAhead() != TOKEN_OP_CONCAT {
		return exp
	}
	line := 0
	exps := []Exp{exp}
	for lexer.LookAhead() == TOKEN_OP_CONCAT {
		line, _, _ = lexer.NextToken() // ..
		exps = append(exps, parseExp4(lexer))
	}
	return &ConcatExp{
		Line: line,
		Exps: exps,
	}
}

// exp4  ::= exp3 {(‘+’ | ‘-’) exp3}
func parseExp4(lexer *Lexer) Exp {
	exp := parseExp3(lexer)
	for {
		switch lexer.LookAhead() {
		case TOKEN_OP_ADD, TOKEN_OP_SUB:
			line, op, _ := lexer.NextToken()
			exp = &BinopExp{line, op, exp, parseExp3(lexer)}
		default:
			return exp
		}
	}
}

// exp3  ::= exp2 {(‘*’ | ‘/’ | ‘//’ | ‘%’) exp2}
func parseExp3(lexer *Lexer) Exp {
	exp := parseExp2(lexer)
	for {
		switch lexer.LookAhead() {
		case TOKEN_OP_MUL, TOKEN_OP_MOD, TOKEN_OP_DIV, TOKEN_OP_IDIV:
			line, op, _ := lexer.NextToken()
			exp = &BinopExp{line, op, exp, parseExp2(lexer)}
		default:
			return exp
		}
	}
}

// 右结合
// exp2  ::= {(‘not’ | ‘#’ | ‘-’ | ‘~’)} exp1
func parseExp2(lexer *Lexer) Exp {
	switch lexer.LookAhead() {
	case TOKEN_OP_UNM, TOKEN_OP_BNOT, TOKEN_OP_LEN, TOKEN_OP_NOT:
		line, op, _ := lexer.NextToken()
		return &UnopExp{
			Line: line,
			Op:   op,
			// not (not 1)
			Exp: parseExp2(lexer),
		}
	}
	return parseExp1(lexer)

}

// 右结合
// exp1  ::= exp0 {‘^’ exp2}
func parseExp1(lexer *Lexer) Exp {
	// 左子树只有一个exp
	exp := parseExp0(lexer)
	if lexer.LookAhead() == TOKEN_OP_POW {
		line, op, _ := lexer.NextToken()
		exp = &BinopExp{
			Line: line,
			Op:   op,
			Exp1: exp,
			// 构造整个右子树
			Exp2: parseExp2(lexer),
		}
	}
	return exp
}

// exp0  ::= nil | false | true | Numeral | LiteralString
// | ‘...’ | functiondef | prefixexp | tableconstructor
func parseExp0(lexer *Lexer) Exp {
	switch lexer.LookAhead() {
	case TOKEN_VARARG: // ...
		line, _, _ := lexer.NextToken()
		return &VarargExp{line}
	case TOKEN_KW_NIL: // nil
		line, _, _ := lexer.NextToken()
		return &NilExp{line}
	case TOKEN_KW_TRUE: // true
		line, _, _ := lexer.NextToken()
		return &TrueExp{line}
	case TOKEN_KW_FALSE: // false
		line, _, _ := lexer.NextToken()
		return &FalseExp{line}
	case TOKEN_STRING: // LiteralString
		line, _, token := lexer.NextToken()
		return &StringExp{line, token}
	case TOKEN_NUMBER: // Numeral
		return parseNumberExp(lexer)
	case TOKEN_SEP_LCURLY: // tableconstructor
		return parseTableConstructorExp(lexer)
	case TOKEN_KW_FUNCTION: // functiondef
		lexer.NextToken()
		return parseFuncDefExp(lexer)
	default: // prefixexp
		return parsePrefixExp(lexer)
	}
}

func parseNumberExp(lexer *Lexer) Exp {
	line, _, token := lexer.NextToken()
	if i, ok := number.ParseInteger(token); ok {
		return &IntegerExp{
			Line: line,
			Val:  i,
		}
	} else if f, ok := number.ParseFloat(token); ok {
		return &FloatExp{
			Line: line,
			Val:  f,
		}
	}
	panic("not a number: " + token)
}