package mathxf

import (
	"fmt"
	"unicode"
)

// Token represents a Token or text string returned from the scanner.
type Token struct {
	typ  tokenType // The type of this Token.
	line int       // The line number of this Token.
	col  int       // The column number of this Token.
	val  string    // The value of this Token.
}

func (t Token) String() string {
	return fmt.Sprintf("Token{typ: %v, line: %d,col: %d, Val: %s}", t.typ, t.line, t.col, t.val)
}

// itemType identifies the type of lex tokens.
type tokenType int

const (
	TokenError tokenType = iota
	TokenEOF
	TokenComplex
	TokenNumber
	TokenCharConstant // 'a'
	TokenComma        // ,
	TokenSemicolon    // ;
	TokenAdd          // +
	TokenSub          // -
	TokenMul          // *
	TokenDiv          // /
	TokenMod          // %
	TokenBool         // true or false
	TokenTernary      // ?
	TokenLessEquals   //<=
	TokenLess         // <
	TokenGreatEquals  //>=
	TokenGreat        // >
	TokenNotEquals    // != or <>
	TokenEquals       // ==
	TokenAssign       // := or =
	TokenColon        // :
	TokenPow          // ^

	TokenString
	TokenIdentifier // alphanumeric identifier not starting with '.'
	TokenField      // alphanumeric identifier starting with '.'

	TokenChar             // character constant
	TokenLeftParen        // (
	TokenRightParen       // )
	TokenLeftBrackets     // [
	TokenRightBrackets    // ]
	TokenLeftBigBrackets  // {
	TokenRightBigBrackets // {

	TokenAnd // && or and
	TokenOr  // || or or
	TokenNot // ! or not
	TokenIn
	TokenNil // nil
)

const (
	comment      = "//"
	leftComment  = "/*"
	rightComment = "*/"
)

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r' || r == '\n'
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
