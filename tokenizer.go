package config

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Name int

const (
	NameBlockEnd Name = iota
	NameBlockStart
	NameComma
	NameEq
	NameIdent
	NameString
)

type Token struct {
	Name  Name
	Value string
}

type Tokenizer struct {
	r      *strings.Reader
	line   int
	last   *Token
	unread bool
}

func NewTokenizer(s string) *Tokenizer {
	return &Tokenizer{
		r:    strings.NewReader(s),
		line: 1,
	}
}

func (t *Tokenizer) Line() int {
	return t.line
}

func (t *Tokenizer) HasNext() bool {
	t.eatWS()

	return t.unread || t.r.Len() > 0
}

func (t *Tokenizer) Unread() {
	t.unread = true
}

func (t *Tokenizer) Next() (*Token, error) {
	if t.unread {
		t.unread = false
		return t.last, nil
	}

	t.eatWS()

	var tok *Token
	var err error
	var r rune
	r, _, err = t.r.ReadRune()
	if err != nil {
		return nil, err
	}

	if r == '}' {
		tok, err = &Token{NameBlockEnd, "}"}, nil
	} else if r == '{' {
		tok, err = &Token{NameBlockStart, "{"}, nil
	} else if r == ',' {
		tok, err = &Token{NameComma, ","}, nil
	} else if r == '=' {
		tok, err = &Token{NameEq, "="}, nil
	} else if unicode.IsLetter(r) || unicode.IsDigit(r) {
		t.r.UnreadRune()
		id, e := t.readIdent()
		tok, err = &Token{NameIdent, id}, e
	} else if r == '"' {
		t.r.UnreadRune()
		s, e := t.readString()

		tok, err = &Token{NameString, s}, e
	} else {
		tok, err = nil, fmt.Errorf("unexpected `%c`", r)
	}

	if err == nil {
		t.last = tok
	}

	return tok, err
}

func (t *Tokenizer) readIdent() (string, error) {
	id := ""

	for {
		r, _, err := t.r.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		if t.lexemeEnd(r) {
			t.r.UnreadRune()
			break
		}
		id += string(r)
	}

	return id, nil
}

func (t *Tokenizer) readString() (string, error) {
	r, _, err := t.r.ReadRune()
	if r != '"' || err != nil {
		return "", errors.New("`\"` character expected")
	}

	s := ""
	for {
		r, _, err := t.r.ReadRune()
		if err != nil {
			return "", err
		}
		if r == '"' {
			break
		}
		if r == '\\' {
			r, _, err := t.r.ReadRune()
			if err != nil {
				return "", errors.New("unknown escape sequence: EOF")
			}
			s += string(r)
		} else {
			s += string(r)
		}
	}

	return s, nil
}

func (t *Tokenizer) lexemeEnd(r rune) bool {
	return unicode.IsSpace(r) ||
		r == ';' ||
		r == ',' ||
		r == '\n' ||
		r == '#'
}

// Eat all whitespaces, lexeme delimeters and comments.
func (t *Tokenizer) eatWS() {
	for {
		r, _, err := t.r.ReadRune()
		if err != nil {
			break
		}
		if r == '\n' {
			t.line++
		} else if r == '#' {
			for {
				r, _, err := t.r.ReadRune()
				if err != nil {
					break
				}
				if r == '\n' {
					t.r.UnreadRune()
					break
				}

			}
		} else if unicode.IsSpace(r) || r == ';' {
			// Just eat.
		} else {
			t.r.UnreadRune()
			break
		}
	}
}
