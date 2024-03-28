package mathxf

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const eof = -1

// lexer holds the state of the scanner.
type lexer struct {
	input           string            // the string being scanned
	replaceKeywords map[string]string // todo replace keywords

	lastTokenType tokenType  // last Token type
	tokens        chan Token // channel of scanned tokens

	start int
	pos   int
	line  int
	col   int
	width int
}

func lex(input string) *lexer {
	l := &lexer{
		input:  input,
		line:   1,
		col:    0,
		tokens: make(chan Token),
	}
	return l
}

// run runs the state machine for the lexer.
func (l *lexer) run() {
	go func() {
		state := baseStateFn
		for state != nil {
			state = state(l)
		}
		l.emit(TokenEOF)
		close(l.tokens)
	}()
}

func (l *lexer) nextToken() Token {
	item := <-l.tokens
	return item
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	l.col += l.width
	return r
}
func (l *lexer) backup() {
	l.pos -= l.width
	l.col -= l.width
}
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// accept consumes the next rune if it'name from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
	//l.startLine = l.line
	//l.startCol = l.col
}

func (l *lexer) value() string {
	return l.input[l.start:l.pos]
}

// emit passes an item back to the client.
func (l *lexer) emit(t tokenType) {
	l.lastTokenType = t
	l.tokens <- Token{t, l.line, l.col, l.value()}
	l.start = l.pos
}

func (l *lexer) emitError(format string, args ...interface{}) stateFn {
	l.tokens <- Token{TokenError, l.line, l.col, fmt.Sprintf(format, args...)}
	return nil
}

func (l *lexer) scanNumber() (bool, bool) {
	// Is it hex?
	digits := "0123456789"
	if l.accept("0") && l.accept("xX") {
		digits = "0123456789abcdefABCDEF"
	}
	l.acceptRun(digits)
	if l.accept(".") {
		pos := l.pos
		l.acceptRun(digits)
		if pos == l.pos {
			l.next()
			return false, false
		}
	}
	if l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789")
	}
	//Is it imaginary?
	isComplex := l.accept("i")
	//Next thing mustn't be alphanumeric.
	if isAlphaNumeric(l.peek()) {
		l.next()
		return false, false
	}
	return true, isComplex
}
