package mathxf

import (
	"strings"
	"unicode"
)

// stateFn represents the state of the scanner as a function that returns the NextToken state.
type stateFn func(l *lexer) stateFn

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
		n := l.next()
		if n == '=' {
			l.emit(TokenLessEquals)
		} else if n == '>' {
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
		l.ignore()
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
	case r == ']':
		l.emit(TokenRightBrackets)
	case r == '{':
		l.emit(TokenLeftBigBrackets)
	case r == '}':
		l.emit(TokenRightBigBrackets)
	case r == '(':
		l.emit(TokenLeftParen)
	case r == ')':
		l.emit(TokenRightParen)
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
		case eof:
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
			token, ok := TokenKeywords[word]
			if ok {
				l.emit(token)
			} else {
				l.emit(TokenIdentifier)
			}
			return baseStateFn
		}
	}
}
