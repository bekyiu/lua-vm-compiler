package parser

import (
	. "write_lua/src/compiler/ast"
	. "write_lua/src/compiler/lexer"
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
