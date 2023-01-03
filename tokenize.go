package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TokenKind int

const (
	TK_RESERVED TokenKind = iota
	TK_NUM
	TK_EOF
)

type Token struct {
	next *Token
	kind TokenKind
	val  int
	str  string
	len  int
}

func errorAt(loc string, format ...string) {
	fmt.Print(format)
	fmt.Println()
	os.Exit(1)
}

func NewToken(kind TokenKind, cur *Token, str string, len int) *Token {
	tok := &Token{
		kind: kind,
		str:  str,
		len:  len,
	}
	cur.next = tok
	return tok
}

func isWhiteSpace(r rune) bool {
	return r == ' ' || r == '\n' || r == '\t'
}

func isLetter(r rune) bool {
	return ('A' <= r && r <= 'Z') || ('a' <= r && r <= 'z')
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func isAlNum(r rune) bool {
	return isDigit(r) || isLetter(r)
}

func isPunct(r rune) bool {
	return strings.ContainsRune("+-*/=()<>!;:,.", r)
}

func (p *Parser) consume(op string) bool {
	if p.token.kind != TK_RESERVED || p.token.str != op {
		return false
	}
	p.token = p.token.next
	return true
}

func (p *Parser) expect(op string) {
	if p.token.kind != TK_RESERVED || p.token.str != op {
		errorAt(p.token.str, "expected '%s'", op)
	}
	p.token = p.token.next
}

func (p *Parser) expectNumber() int {
	if p.token.kind != TK_NUM {
		errorAt(p.token.str, "expect a number")
	}
	val := p.token.val
	p.token = p.token.next
	return val
}

func (t *Token) AtEOF() bool {
	return t.kind == TK_EOF
}

func Tokenize(input string) *Token {
	head := &Token{}
	cur := head
	i := 0
	for i = 0; i < len(input); {
		if isWhiteSpace(rune(input[i])) {
			i++
			continue
		}
		if strings.HasPrefix(input[i:], "return") && !isAlNum(rune(input[i+6])) {
			cur = NewToken(TK_RESERVED, cur, input[i:i+6], 6)
			i += 6
			continue
		}
		if strings.HasPrefix(input[i:], "==") ||
			strings.HasPrefix(input[i:], "!=") ||
			strings.HasPrefix(input[i:], "<=") ||
			strings.HasPrefix(input[i:], ">=") {
			cur = NewToken(TK_RESERVED, cur, input[i:i+2], 2)
			i += 2
			continue
		}
		if isPunct(rune(input[i])) {
			cur = NewToken(TK_RESERVED, cur, input[i:i+1], 1)
			i++
			continue
		}
		if isDigit(rune(input[i])) {
			pos := i
			for ; i < len(input) && isDigit(rune(input[i])); i++ {
			}
			cur = NewToken(TK_NUM, cur, input[pos:i], i-pos)
			val, _ := strconv.ParseInt(input[pos:i], 10, 32)
			cur.val = int(val)
			continue
		}

		errorAt(input[i:], "invalid token")
	}

	NewToken(TK_EOF, cur, input[i:], 0)
	return head.next
}
