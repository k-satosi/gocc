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
	TK_IDENT
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

func errorAt(loc string, format string, a ...string) {
	fmt.Println(loc)
	fmt.Printf(format, a)
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
	return ('A' <= r && r <= 'Z') || ('a' <= r && r <= 'z') || r == '_'
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func isAlNum(r rune) bool {
	return isDigit(r) || isLetter(r)
}

func isPunct(r rune) bool {
	return strings.ContainsRune("+-*/=(){}<>!;:,.", r)
}

func (p *Parser) consume(op string) bool {
	if p.token.kind != TK_RESERVED || p.token.str != op {
		return false
	}
	p.token = p.token.next
	return true
}

func (p *Parser) consumeIdent() *Token {
	if p.token.kind != TK_IDENT {
		return nil
	}
	token := p.token
	p.token = token.next
	return token
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

func (p *Parser) expectIdent() string {
	if p.token.kind != TK_IDENT {
		errorAt(p.token.str, "expected an identifier")
	}
	s := p.token.str
	p.token = p.token.next
	return s
}

func (t *Token) AtEOF() bool {
	return t.kind == TK_EOF
}

func startsWithReserved(s string) (string, bool) {
	keywords := []string{"return", "if", "else", "while", "for"}
	for _, v := range keywords {
		if strings.HasPrefix(s, v) {
			return v, true
		}
	}

	ops := []string{"==", "!=", "<=", ">="}
	for _, v := range ops {
		if strings.HasPrefix(s, v) {
			return v, true
		}
	}

	return "", false
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
		if keyword, ok := startsWithReserved(input[i:]); ok {
			cur = NewToken(TK_RESERVED, cur, keyword, len(keyword))
			i += len(keyword)
			continue
		}
		if isLetter(rune(input[i])) {
			pos := i
			i++
			for ; i < len(input) && isAlNum(rune(input[i])); i++ {
			}
			cur = NewToken(TK_IDENT, cur, input[pos:i], i-pos)
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
