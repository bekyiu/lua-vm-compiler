package parser

import (
	. "write_lua/src/compiler/ast"
	. "write_lua/src/compiler/lexer"
)

// prefixexp ::= var | functioncall | ‘(’ exp ‘)’
// var ::=  Name | prefixexp ‘[’ exp ‘]’ | prefixexp ‘.’ Name
// functioncall ::=  prefixexp args | prefixexp ‘:’ Name args

/*
prefixexp ::= Name
	| ‘(’ exp ‘)’
	| prefixexp ‘[’ exp ‘]’
	| prefixexp ‘.’ Name
	| prefixexp [‘:’ Name] args
*/
func parsePrefixExp(lexer *Lexer) Exp {
	var exp Exp
	if lexer.LookAhead() == TOKEN_IDENTIFIER {
		line, name := lexer.NextIdentifier() // name
		exp = &NameExp{line, name}
	} else {
		exp = parseParensExp(lexer)
	}
	return _finishPrefixExp(lexer, exp)
}

func parseParensExp(lexer *Lexer) Exp {
	lexer.NextTokenOfKind(TOKEN_SEP_LPAREN) // (
	exp := parseExp(lexer)                  // exp
	lexer.NextTokenOfKind(TOKEN_SEP_RPAREN) // )
	switch exp.(type) {
	case *VarargExp, *FuncCallExp, *NameExp, *TableAccessExp:
		return &ParensExp{exp}
	}
	// 无需保留圆括号
	return exp
}

func _finishPrefixExp(lexer *Lexer, exp Exp) Exp {
	for {
		switch lexer.LookAhead() {
		case TOKEN_SEP_LBRACK:
			lexer.NextToken()                       // [
			keyExp := parseExp(lexer)               // exp
			lexer.NextTokenOfKind(TOKEN_SEP_RBRACK) // ]
			exp = &TableAccessExp{lexer.Line(), exp, keyExp}
		case TOKEN_SEP_DOT:
			lexer.NextToken()                    // .
			line, name := lexer.NextIdentifier() // name
			keyExp := &StringExp{line, name}
			exp = &TableAccessExp{line, exp, keyExp}
		case TOKEN_SEP_COLON, TOKEN_SEP_LPAREN, TOKEN_SEP_LCURLY, TOKEN_STRING: // [‘:’ Name] args
			exp = _finishFuncCallExp(lexer, exp)
		default:
			return exp
		}
	}
}

// [‘:’ Name] args
func _finishFuncCallExp(lexer *Lexer, prefixExp Exp) *FuncCallExp {
	nameExp := _parseNameExp(lexer) // [':' name]
	line := lexer.Line()            // todo
	args := _parseArgs(lexer)       // args
	lastLine := lexer.Line()
	return &FuncCallExp{line, lastLine, prefixExp, nameExp, args}
}

func _parseNameExp(lexer *Lexer) *StringExp {
	if lexer.LookAhead() == TOKEN_SEP_COLON {
		lexer.NextToken() // :
		line, name := lexer.NextIdentifier()
		return &StringExp{line, name}
	}
	return nil
}

// args ::=  ‘(’ [explist] ‘)’ | tableconstructor | LiteralString
func _parseArgs(lexer *Lexer) (args []Exp) {
	switch lexer.LookAhead() {
	case TOKEN_SEP_LPAREN: // ‘(’ [explist] ‘)’
		lexer.NextToken()
		if lexer.LookAhead() != TOKEN_SEP_RPAREN {
			args = parseExpList(lexer)
		}
		lexer.NextTokenOfKind(TOKEN_SEP_RPAREN)
	case TOKEN_SEP_LCURLY: // ‘{’ [fieldlist] ‘}’
		args = []Exp{parseTableConstructorExp(lexer)}
	default: // LiteralString
		line, str := lexer.NextTokenOfKind(TOKEN_STRING)
		args = []Exp{&StringExp{line, str}}
	}
	return
}
