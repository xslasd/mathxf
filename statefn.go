package mathxf

import (
	"fmt"
	"strings"
	"unicode"
)

// stateFn represents the state of the scanner as a function that returns the nextToken state.
type stateFn func(l *lexer) stateFn

func textStateFn(l *lexer) stateFn {
	for l.state = baseStateFn; l.state != nil; {
		l.state = l.state(l)
	}
	if l.bracketDepth > 0 {
		l.emitError("unexpected left bracket %#U", '[')
	}
	if l.parenDepth > 0 {
		l.emitError("unexpected left paren %#U", '(')
	}
	if l.bigBracketDepth > 0 {
		l.emitError("unexpected left bracket %#U", '{')
	}
	l.emit(TokenEOF)
	return nil
}

func baseStateFn(l *lexer) stateFn {
	if strings.HasPrefix(l.input[l.pos:], comment) {
		return comment1StateFn(l)
	}
	if strings.HasPrefix(l.input[l.pos:], leftComment) {
		return comment2StateFn(l)
	}
	switch r := l.next(); {
	case r == eof:
		return nil
	case isSpace(r):
		if r == '\n' {
			l.line++
			l.col = 0
		}
		l.ignore()
		return baseStateFn
	case r == '"':
		fmt.Println("string==========")
		return stringStateFn
	case r == ',':
		l.emit(TokenComma)
	case r == ';':
		l.emit(TokenSemicolon)
	case r == '*':
		l.emit(TokenMul)
	case r == '/':
		l.emit(TokenDiv)
	case r == '%':
		l.emit(TokenMod)
	case r == '-':
		if pr := l.peek(); (('0' <= pr && pr <= '9') || pr == '.') &&
			TokenNumber != l.lastTokenType &&
			TokenIdentifier != l.lastTokenType &&
			TokenBool != l.lastTokenType &&
			TokenField != l.lastTokenType &&
			TokenChar != l.lastTokenType {
			return numberStateFn
		}
		l.emit(TokenSub)
	case r == '+':
		if r := l.peek(); '0' <= r && r <= '9' &&
			TokenAdd != l.lastTokenType &&
			TokenSub != l.lastTokenType &&
			TokenNumber != l.lastTokenType &&
			TokenIdentifier != l.lastTokenType &&
			TokenBool != l.lastTokenType &&
			TokenField != l.lastTokenType &&
			TokenChar != l.lastTokenType {
			return numberStateFn
		}
		l.emit(TokenAdd)
	case r == '?':
		l.emit(TokenTernary)
	case r == '^':
		l.emit(TokenPow)
	case r == '&':
		if l.next() == '&' {
			l.emit(TokenAnd)
		} else {
			l.backup()
		}
	case r == '<':
		if l.next() == '=' {
			l.emit(TokenLessEquals)
		} else if l.next() == '>' {
			l.emit(TokenNotEquals)
		} else {
			l.backup()
			l.emit(TokenLess)
		}
	case r == '>':
		if l.next() == '=' {
			l.emit(TokenGreatEquals)
		} else {
			l.backup()
			l.emit(TokenGreat)
		}
	case r == '!':
		if l.next() == '=' {
			l.emit(TokenNotEquals)
		} else {
			l.backup()
			l.emit(TokenNot)
		}
	case r == '=':
		if l.next() == '=' {
			l.emit(TokenEquals)
		} else {
			l.backup()
			l.emit(TokenAssign)
		}
	case r == ':':
		if l.next() == '=' {
			l.emit(TokenAssign)
		} else {
			l.backup()
			l.emit(TokenColon)
		}
	case r == '|':
		if l.next() == '|' {
			l.emit(TokenOr)
		} else {
			return l.emitError("unrecognized character in action: %#U", r)
		}
	case r == '.':
		if r := l.peek(); '0' <= r && r <= '9' {
			return numberStateFn
		}
		for {
			switch r := l.next(); {
			case isAlphaNumeric(r):
			default:
				l.backup()
				l.emit(TokenField)
				return baseStateFn
			}
		}

	case '0' <= r && r <= '9':
		return numberStateFn
	case isAlphaNumeric(r):
		return identifierStateFn
	case r == '[':
		l.emit(TokenLeftBrackets)
		l.bracketDepth++
	case r == ']':
		l.emit(TokenRightBrackets)
		l.bracketDepth--
		if l.bracketDepth < 0 {
			return l.emitError("unexpected right bracket %#U", r)
		}
	case r == '{':
		l.emit(TokenLeftBigBrackets)
		l.bigBracketDepth++
	case r == '}':
		l.emit(TokenRightBigBrackets)
		l.bigBracketDepth--
		if l.bigBracketDepth < 0 {
			return l.emitError("unexpected right bracket %#U", r)
		}
	case r == '(':
		l.emit(TokenLeftParen)
		l.parenDepth++
	case r == ')':
		l.emit(TokenRightParen)
		l.parenDepth--
		if l.parenDepth < 0 {
			return l.emitError("unexpected right paren %#U", r)
		}
	case r == '\'':
		return charStateFn
	case r <= unicode.MaxASCII && unicode.IsPrint(r):
		l.emit(TokenChar)
	default:
		return l.emitError("unrecognized character in action: %#U", r)
	}
	return baseStateFn
}
func charStateFn(l *lexer) stateFn {
Loop:
	for {
		switch l.next() {
		case '\\':
			if r := l.next(); r != eof && r != '\n' {
				break
			}
			fallthrough
		case eof, '\n':
			return l.emitError("unterminated character constant")
		case '\'':
			break Loop
		}
	}
	l.emit(TokenCharConstant)
	return baseStateFn
}
func stringStateFn(l *lexer) stateFn {
	l.ignore()
	for {
		r := l.next()
		if r == eof || r == '\n' {
			return l.emitError("unterminated string constant")
		}
		if r == '"' {
			l.backup()
			break
		}
	}
	l.emit(TokenString)
	l.next()
	return baseStateFn
}

func numberStateFn(l *lexer) stateFn {
	isNumber, isComplex := l.scanNumber()
	if !isNumber {
		return l.emitError("bad number syntax: %q", l.value())
	}
	if isComplex {
		l.emit(TokenComplex)
	} else {
		l.emit(TokenNumber)
	}
	return baseStateFn
}
func comment1StateFn(l *lexer) stateFn {
	l.pos += len(comment)
	i := strings.IndexRune(l.input[l.pos:], '\n')
	if i < 0 {
		return nil
	}
	l.pos += i + 1
	l.line++
	l.col = 0
	l.ignore()
	return baseStateFn
}
func comment2StateFn(l *lexer) stateFn {
	l.pos += len(leftComment)
	i := strings.Index(l.input[l.pos:], rightComment)
	if i < 0 {
		return l.emitError("unclosed comment")
	}
	i2 := strings.Count(l.input[l.pos:i+2], string('\n'))
	l.line += i2
	l.pos += i + len(rightComment)
	l.col = 0
	l.ignore()
	return baseStateFn
}

func identifierStateFn(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r):
		default:
			l.backup()
			word := l.input[l.start:l.pos]
			//token, ok := TokenKeywords[word]
			//if ok {
			//	l.emit(token)
			//} else {
			switch word {
			case "true", "false":
				l.emit(TokenBool)
			case "in":
				l.emit(TokenIn)
			case "and":
				l.emit(TokenAnd)
			case "or":
				l.emit(TokenOr)
			case "nil":
				l.emit(TokenNil)
			default:
				l.emit(TokenIdentifier)
			}

			//}
			return baseStateFn
		}
	}
}

func fieldStateFn(l *lexer) stateFn {
	for {
		r := l.next()
		if !isAlphaNumeric(r) {
			l.backup()
			break
		}
	}
	if l.pos-l.start <= 1 {
		l.next()
		return l.emitError("bad character: %q", l.value())
	}
	l.emit(TokenField)
	return baseStateFn
}
