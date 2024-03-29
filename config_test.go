package config

import (
	"reflect"
	"testing"
)

func TestConfigProperty(t *testing.T) {
	cfg := &Config{
		[]*Property{
			&Property{TypeInt, "foo", 123},
			&Property{TypeString, "bar", "value"},
			&Property{TypeStringList, "baz", []string{"foo", "bar"}},
		},
		nil,
	}

	if !cfg.Has("foo") {
		t.Fatal()
	}
	if cfg.Int("foo") != 123 {
		t.Fatal()
	}
	if cfg.IntOr("not-foo", 321) != 321 {
		t.Fatal()
	}
	if !cfg.Has("bar") {
		t.Fatal()
	}
	if cfg.String("bar") != "value" {
		t.Fatal()
	}
	if cfg.StringOr("not-bar", "default") != "default" {
		t.Fatal()
	}
	if !cfg.Has("baz") {
		t.Fatal()
	}
	if !reflect.DeepEqual(cfg.StringList("baz"), []string{"foo", "bar"}) {
		t.Fatal()
	}
	if !reflect.DeepEqual(cfg.StringListOr("not-baz", []string{"qux"}),
		[]string{"qux"}) {
		t.Fatal()
	}
}

func TestConfigBlock(t *testing.T) {
	cfg := &Config{
		nil,
		[]*Block{
			&Block{
				"foo",
				[]*Property{
					&Property{TypeInt, "foo-foo", 1},
					&Property{TypeInt, "foo-bar", 2},
				},
				[]*Block{
					&Block{
						"bar",
						[]*Property{
							&Property{TypeInt, "bar-foo", 3},
							&Property{TypeInt, "bar-bar", 4},
						},
						nil,
					},
				},
			},
		},
	}

	if !cfg.Has("foo") {
		t.Fatal()
	}
	foo := cfg.Block("foo")
	if foo == nil {
		t.Fatal()
	}
	if !foo.Has("bar") {
		t.Fatal()
	}
	bar := foo.Block("bar")
	if bar == nil {
		t.Fatal()
	}
	if bar.Int("bar-foo") != 3 {
		t.Fatal()
	}
	baz := bar.Block("baz")
	if baz != nil {
		t.Fatal()
	}
}
