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
	return nil
}
