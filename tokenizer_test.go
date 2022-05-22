package config

import (
	"testing"
)

func TestTokenizerString(t *testing.T) {
	testTokensSerie(t, `""`, &Token{NameString, ""})
	testTokensSerie(t, `" "`, &Token{NameString, " "})
	testTokensSerie(t, `"foo bar"`, &Token{NameString, "foo bar"})
	testTokensSerie(t, `"foo\"bar"`, &Token{NameString, "foo\"bar"})
	testTokensSerie(t, `"\"\'\\"`, &Token{NameString, "\"'\\"})
}

func TestTokenizer(t *testing.T) {
	testTokensSerie(t, "foo bar baz",
		&Token{NameIdent, "foo"},
		&Token{NameIdent, "bar"},
		&Token{NameIdent, "baz"})
	testTokensSerie(t, "foo = bar",
		&Token{NameIdent, "foo"},
		&Token{NameEq, "="},
		&Token{NameIdent, "bar"})
	testTokensSerie(t, "foo = 123xxx123",
		&Token{NameIdent, "foo"},
		&Token{NameEq, "="},
		&Token{NameIdent, "123xxx123"})
	testTokensSerie(t, "foo = 1, bar, \"baz\"",
		&Token{NameIdent, "foo"},
		&Token{NameEq, "="},
		&Token{NameIdent, "1"},
		&Token{NameComma, ","},
		&Token{NameIdent, "bar"},
		&Token{NameComma, ","},
		&Token{NameString, "baz"})
	testTokensSerie(t, "block {foo = 1; bar = 2;}",
		&Token{NameIdent, "block"},
		&Token{NameBlockStart, "{"},
		&Token{NameIdent, "foo"},
		&Token{NameEq, "="},
		&Token{NameIdent, "1"},
		&Token{NameIdent, "bar"},
		&Token{NameEq, "="},
		&Token{NameIdent, "2"},
		&Token{NameBlockEnd, "}"})
	testTokensSerie(t, "heartbeat-ttl = 3s\n\n# block {\n#}\n",
		&Token{NameIdent, "heartbeat-ttl"},
		&Token{NameEq, "="},
		&Token{NameIdent, "3s"})
	testTokensSerie(t, ""+
		"# Comment line.\n"+
		"param-a = 1\n"+
		"\n"+
		"param-b = 2;\n"+
		"block-a {\n"+
		"    param-c = 3\n"+
		"    param-d = 4;\n"+
		"}"+
		"block-b {\n"+
		"    param-d = \"value-d\"\n"+
		"    param-e = value-e;\n"+
		"}",
		&Token{NameIdent, "param-a"},
		&Token{NameEq, "="},
		&Token{NameIdent, "1"},
		&Token{NameIdent, "param-b"},
		&Token{NameEq, "="},
		&Token{NameIdent, "2"},
		&Token{NameIdent, "block-a"},
		&Token{NameBlockStart, "{"},
		&Token{NameIdent, "param-c"},
		&Token{NameEq, "="},
		&Token{NameIdent, "3"},
		&Token{NameIdent, "param-d"},
		&Token{NameEq, "="},
		&Token{NameIdent, "4"},
		&Token{NameBlockEnd, "}"},
		&Token{NameIdent, "block-b"},
		&Token{NameBlockStart, "{"},
		&Token{NameIdent, "param-d"},
		&Token{NameEq, "="},
		&Token{NameString, "value-d"},
		&Token{NameIdent, "param-e"},
		&Token{NameEq, "="},
		&Token{NameIdent, "value-e"},
		&Token{NameBlockEnd, "}"})
}

func testTokensSerie(t *testing.T, s string, toks ...*Token) {
	tk := NewTokenizer(s)

	for _, tok := range toks {
		if !tk.HasNext() {
			t.Fatal()
		}
		act, err := tk.Next()
		if err != nil {
			t.Fatal(err)
		}
		if *act != *tok {
			t.Fatalf("%v != %v", act, tok)
		}
		tk.Unread()
		if !tk.HasNext() {
			t.Fatal()
		}
		act, err = tk.Next()
		if err != nil {
			t.Fatal(err)
		}
		if *act != *tok {
			t.Fatalf("%v != %v", act, tok)
		}
	}
	if tk.HasNext() {
		t.Fatal()
	}
}
