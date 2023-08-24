package config

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Type int

// TODO: Add Float type support.
const (
	// Typical boolean value like true and false.
	TypeBool Type = iota
	// Duration type. Same format as Go's time.ParseDuration() uses.
	// See more: https://pkg.go.dev/time#ParseDuration
	TypeDuration
	// Positive integer value like 1 or 100.
	// TODO: Support negative values.
	TypeInt
	// String value -- sequence of characters enclosed with double quotes.
	TypeString
	// List of strings.
	TypeStringList
)

// Property value custom parser function.
// Also can be used to validate parsed value.
type Parser func(any) (any, error)

// Specification descriptor for single property in configuration file.
// Proper is a key=value parir expression which assigns value for an
// identifier.
//
// Examples:
// property-name = "string value"
// property-name = 100
type PropertySpec struct {
	Type    Type
	Name    string
	Repeat  bool
	Require bool
	Parser  Parser
}

// Specification descriptor for block of properties.
// Block is a named group of properties. Nested blocks are allowed.
//
// Examples:
// block-a {
//     foo = 1
//     block-b {
//         bar = 2
//         baz = 3
//     }
// }
//
// There is a reserved block name which has special meaning -- *. Only single
// star-block per nested-level is allowed. Star-block allows to define expected
// serie of block with arbitrary names. It can be useful when you want to allow
// series of uniform blocks.
//
// Example:
// sda {
//     dev = "/dev/sda"
// }
// sdb {
//     dev = "/dev/sdb"
// }
type BlockSpec struct {
	Name       string
	Repeat     bool
	Require    bool
	Properties []*PropertySpec
	Blocks     []*BlockSpec
	Strict     bool
}

type Spec struct {
	Properties []*PropertySpec
	Blocks     []*BlockSpec
	Strict     bool
}

const (
	rootBlock = ""
)

func Parse(spec *Spec, s string) (*Config, error) {
	t := NewTokenizer(s)
	rs := &BlockSpec{
		Name:       rootBlock,
		Repeat:     false,
		Properties: spec.Properties,
		Blocks:     spec.Blocks,
		Strict:     spec.Strict,
	}
	b, err := parseBlock(t, rs.Name, rs)
	if err != nil {
		return nil, err
	}

	return &Config{b.Properties, b.Blocks}, nil
}

func parseBlock(t *Tokenizer, name string, spec *BlockSpec) (*Block, error) {
	var props []*Property
	var blocks []*Block
	var closed bool = name == rootBlock

	for t.HasNext() {
		n, err := t.Next()
		if err != nil {
			return nil, newError(t.Line(), err.Error())
		}
		if name != rootBlock && n.Name == NameBlockEnd {
			closed = true
			break
		}
		if n.Name != NameIdent {
			return nil, newError(t.Line(), "identifier token expected")
		}

		op, err := t.Next()
		if err != nil {
			return nil, newError(t.Line(), err.Error())
		}

		switch op.Name {
		case NameEq:
			if !t.HasNext() {
				return nil, newError(t.Line(), "value expected")
			}
			v, err := t.Next()
			if err != nil {
				return nil, newError(t.Line(), err.Error())
			}
			s := findProperty(spec.Properties, n.Value)
			if s == nil {
				if spec.Strict {
					return nil, newError(t.Line(),
						"unsupported property: %s", n.Value)
				} else {
					continue
				}
			}
			i := contains(len(props), func(i int) bool {
				return props[i].Name == n.Value
			})
			if i != -1 {
				if !s.Repeat {
					return nil, newError(t.Line(),
						"property `%s` already defined",
						n.Value)

				}
			}

			var val any
			switch s.Type {
			case TypeBool:
				if v.Name == NameIdent && v.Value == "true" {
					val = true
				} else if v.Name == NameIdent && v.Value == "false" {
					val = false
				} else {
					return nil, newError(t.Line(),
						"invalid boolean value")
				}
			case TypeDuration:
				if v.Name != NameIdent {
					return nil, newError(t.Line(),
						"duration value expected")
				}
				d, err := time.ParseDuration(v.Value)
				if err != nil {
					return nil, newError(t.Line(),
						"invalid duration value")
				}
				val = d
			case TypeInt:
				if v.Name != NameIdent {
					return nil, newError(t.Line(),
						"integer value expected")
				}
				i, err := strconv.Atoi(v.Value)
				if err != nil {
					return nil, newError(t.Line(),
						"invalid integer value")
				}
				val = i
			case TypeString:
				if v.Name != NameString {
					return nil, newError(t.Line(),
						"string value expected")
				}
				val = v.Value
			case TypeStringList:
				// TODO: Add empty list support.
				if v.Name != NameString {
					return nil, newError(t.Line(),
						"strings list expected")
				}
				var lst []string

				lst = append(lst, v.Value)
				for {
					if !t.HasNext() {
						break
					}
					tk, err := t.Next()
					if err != nil {
						return nil, newError(t.Line(),
							err.Error())
					}
					if tk.Name != NameComma {
						t.Unread()
						break
					}
					if !t.HasNext() {
						return nil, newError(t.Line(),
							"unexpected EOF")
					}
					tk, err = t.Next()
					if err != nil {
						return nil, newError(t.Line(),
							err.Error())
					}
					if tk.Name != NameString {
						return nil, newError(t.Line(),
							"string value expected")
					}
					lst = append(lst, tk.Value)
				}
				val = lst
			default:
				panic("unsupported Type")
			}

			if s.Parser != nil {
				val, err = s.Parser(val)
				if err != nil {
					return nil, newError(t.Line(),
						err.Error())
				}
			}

			props = append(props, &Property{
				Type:  s.Type,
				Name:  n.Value,
				Value: val,
			})
		case NameBlockStart:
			s := findBlock(spec.Blocks, n.Value)
			if s == nil {
				return nil, newError(t.Line(),
					"unsupported block: %s", n.Value)
			}
			b, err := parseBlock(t, n.Value, s)
			if err != nil {
				return nil, err
			}
			i := contains(len(blocks), func(i int) bool {
				return blocks[i].Name == n.Value
			})
			if i != -1 {
				if !s.Repeat {
					return nil, newError(t.Line(),
						"block `%s` already defined",
						n.Value)
				}
			}
			blocks = append(blocks, b)
		default:
			ps := findProperty(spec.Properties, n.Value)
			if ps != nil {
				return nil, newError(t.Line(), "`=` expected")
			}
			bs := findBlock(spec.Blocks, n.Value)
			if bs != nil {
				return nil, newError(t.Line(), "`{` expected")
			}
			return nil, newError(t.Line(), "`=` or `{` expected")
		}
	}
	if !closed {
		return nil, newError(t.Line(), "`}` expected")
	}

	for _, s := range spec.Properties {
		if s.Require {
			i := contains(len(props), func(i int) bool {
				return props[i].Name == s.Name
			})
			if i == -1 {
				return nil, newError(t.Line(),
					"missing required property `%s`",
					s.Name)
			}
		}
	}
	for _, s := range spec.Blocks {
		if s.Require {
			i := contains(len(blocks), func(i int) bool {
				return blocks[i].Name == s.Name
			})
			if i == -1 {
				return nil, newError(t.Line(),
					"missing required block `%s`",
					s.Name)
			}
		}
	}

	return &Block{Name: name, Properties: props, Blocks: blocks}, nil
}

func contains(len int, f func(int) bool) int {
	for i := 0; i < len; i++ {
		if f(i) {
			return i
		}
	}

	return -1
}

func findProperty(specs []*PropertySpec, name string) *PropertySpec {
	var ps *PropertySpec

	for _, s := range specs {
		if matchName(name, s.Name) {
			if ps == nil || s.Name > ps.Name {
				ps = s
			}
		}
	}

	return ps
}

func findBlock(specs []*BlockSpec, name string) *BlockSpec {
	var bs *BlockSpec

	for _, s := range specs {
		if matchName(name, s.Name) {
			if bs == nil || s.Name > bs.Name {
				bs = s
			}
		}
	}

	return bs
}

func matchName(s string, pattern string) bool {
	pattern = "^" + strings.ReplaceAll(pattern, "*", ".*") + "$"
	p := regexp.MustCompile(pattern)

	return p.MatchString(s)
}
