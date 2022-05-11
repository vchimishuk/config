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
	NameEq
	NameIdent
	NameString
)

type Token struct {
	Name  Name
	Value string
}

type Tokenizer struct {
	r    *strings.Reader
	line int
}

func NewTokenizer(s string) *Tokenizer {
	return &Tokenizer{strings.NewReader(s), 1}
}

func (t *Tokenizer) Line() int {
	return t.line
}

func (t *Tokenizer) HasNext() bool {
	t.eatWS()

	return t.r.Len() > 0
}

func (t *Tokenizer) Next() (*Token, error) {
	t.eatWS()

	r, _, err := t.r.ReadRune()
	if err != nil {
		return nil, err
	}

	if r == '}' {
		return &Token{NameBlockEnd, "}"}, nil
	} else if r == '{' {
		return &Token{NameBlockStart, "{"}, nil
	} else if r == '=' {
		return &Token{NameEq, "="}, nil
	} else if unicode.IsLetter(r) || unicode.IsDigit(r) {
		t.r.UnreadRune()
		id, err := t.readIdent()

		return &Token{NameIdent, id}, err
	} else if r == '"' {
		t.r.UnreadRune()
		s, err := t.readString()

		return &Token{NameString, s}, err
	} else {
		return nil, fmt.Errorf("unexpected `%c`", r)
	}
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
	return unicode.IsSpace(r) || r == ';' || r == '\n' || r == '#'
}

// Eat all whitespaces, lexeme delimeters and comments.
func (t *Tokenizer) eatWS() {
	for {
		r, _, err := t.r.ReadRune()
		if err != nil {
			break
		}
		if !t.lexemeEnd(r) {
			t.r.UnreadRune()
			if r == '\n' {
				t.line++
			}
			break
		}
		if r == '#' {
			for {
				r, _, err := t.r.ReadRune()
				if err != nil {
					break
				}
				if r == '\n' {
					t.line++
					break
				}
			}
			break
		}
	}
}
