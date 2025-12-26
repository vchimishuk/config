package config

import (
	"fmt"
	"time"
)

type Property struct {
	Type  Type
	Name  string
	Value any
}

type Block struct {
	Name       string
	Properties []*Property
	Blocks     []*Block
}

func (b *Block) Has(name string) bool {
	return property(b.Properties, name) != nil || b.Block(name) != nil
}

func (b *Block) Any(name string) any {
	p := property(b.Properties, name)
	if p == nil {
		panic(fmt.Sprintf("`%s` property is not defined", name))
	}

	return p.Value
}

func (b *Block) AnyOr(name string, defvalue any) any {
	p := property(b.Properties, name)
	if p == nil {
		return defvalue
	}

	return p.Value
}

func (b *Block) Anys(name string) []any {
	return values[any](properties(b.Properties, name))
}

func (b *Block) Bool(name string) bool {
	p := property(b.Properties, name)
	if p == nil {
		panic(fmt.Sprintf("`%s` property is not defined", name))
	}

	return p.Value.(bool)
}

func (b *Block) BoolOr(name string, defvalue bool) bool {
	p := property(b.Properties, name)
	if p == nil {
		return defvalue
	}

	return p.Value.(bool)
}

func (b *Block) Bools(name string) []bool {
	return values[bool](properties(b.Properties, name))
}

func (b *Block) Duration(name string) time.Duration {
	p := property(b.Properties, name)
	if p == nil {
		panic(fmt.Sprintf("`%s` property is not defined", name))
	}

	return p.Value.(time.Duration)
}

func (b *Block) DurationOr(name string, defvalue time.Duration) time.Duration {
	p := property(b.Properties, name)
	if p == nil {
		return defvalue
	}

	return p.Value.(time.Duration)
}

func (b *Block) Durations(name string) []time.Duration {
	return values[time.Duration](properties(b.Properties, name))
}

func (b *Block) Int(name string) int {
	p := property(b.Properties, name)
	if p == nil {
		panic(fmt.Sprintf("`%s` property is not defined", name))
	}

	return p.Value.(int)
}

func (b *Block) IntOr(name string, defvalue int) int {
	p := property(b.Properties, name)
	if p == nil {
		return defvalue
	}

	return p.Value.(int)
}

func (b *Block) Ints(name string) []int {
	return values[int](properties(b.Properties, name))
}

func (b *Block) String(name string) string {
	p := property(b.Properties, name)
	if p == nil {
		panic(fmt.Sprintf("`%s` property is not defined", name))
	}

	return p.Value.(string)
}

func (b *Block) StringOr(name string, defvalue string) string {
	p := property(b.Properties, name)
	if p == nil {
		return defvalue
	}

	return p.Value.(string)
}

func (b *Block) Strings(name string) []string {
	return values[string](properties(b.Properties, name))
}

func (b *Block) StringList(name string) []string {
	p := property(b.Properties, name)
	if p == nil {
		panic(fmt.Sprintf("`%s` property is not defined", name))
	}

	return p.Value.([]string)
}

func (b *Block) StringListOr(name string, defvalue []string) []string {
	p := property(b.Properties, name)
	if p == nil {
		return defvalue
	}

	return p.Value.([]string)
}

func (b *Block) StringLists(name string) [][]string {
	return values[[]string](properties(b.Properties, name))
}

// Block returns block by name or nil if no such block found.
func (b *Block) Block(name string) *Block {
	for _, b := range b.Blocks {
		if b.Name == name {
			return b
		}
	}

	return nil
}

type Config struct {
	Properties []*Property
	Blocks     []*Block
}

func (c *Config) Has(name string) bool {
	return property(c.Properties, name) != nil || c.Block(name) != nil
}

func (c *Config) Any(name string) any {
	p := property(c.Properties, name)
	if p == nil {
		panic(fmt.Sprintf("`%s` property is not defined", name))
	}

	return p.Value
}

func (c *Config) AnyOr(name string, defvalue any) any {
	p := property(c.Properties, name)
	if p == nil {
		return defvalue
	}

	return p.Value
}

func (c *Config) Anys(name string) []any {
	return values[any](properties(c.Properties, name))
}

func (c *Config) Bool(name string) bool {
	p := property(c.Properties, name)
	if p == nil {
		panic(fmt.Sprintf("`%s` property is not defined", name))
	}

	return p.Value.(bool)
}

func (c *Config) BoolOr(name string, defvalue bool) bool {
	p := property(c.Properties, name)
	if p == nil {
		return defvalue
	}

	return p.Value.(bool)
}

func (c *Config) Bools(name string) []bool {
	return values[bool](properties(c.Properties, name))
}

func (c *Config) Duration(name string) time.Duration {
	p := property(c.Properties, name)
	if p == nil {
		panic(fmt.Sprintf("`%s` property is not defined", name))
	}

	return p.Value.(time.Duration)
}

func (c *Config) DurationOr(name string, defvalue time.Duration) time.Duration {
	p := property(c.Properties, name)
	if p == nil {
		return defvalue
	}

	return p.Value.(time.Duration)
}

func (c *Config) Durations(name string) []time.Duration {
	return values[time.Duration](properties(c.Properties, name))
}

func (c *Config) Int(name string) int {
	p := property(c.Properties, name)
	if p == nil {
		panic(fmt.Sprintf("`%s` property is not defined", name))
	}

	return p.Value.(int)
}

func (c *Config) IntOr(name string, defvalue int) int {
	p := property(c.Properties, name)
	if p == nil {
		return defvalue
	}

	return p.Value.(int)
}

func (c *Config) Ints(name string) []int {
	return values[int](properties(c.Properties, name))
}

func (c *Config) String(name string) string {
	p := property(c.Properties, name)
	if p == nil {
		panic(fmt.Sprintf("`%s` property is not defined", name))
	}

	return p.Value.(string)
}

func (c *Config) StringOr(name string, defvalue string) string {
	p := property(c.Properties, name)
	if p == nil {
		return defvalue
	}

	return p.Value.(string)
}

func (c *Config) Strings(name string) []string {
	return values[string](properties(c.Properties, name))
}

func (c *Config) StringList(name string) []string {
	p := property(c.Properties, name)
	if p == nil {
		panic(fmt.Sprintf("`%s` property is not defined", name))
	}

	return p.Value.([]string)
}

func (c *Config) StringListOr(name string, defvalue []string) []string {
	p := property(c.Properties, name)
	if p == nil {
		return defvalue
	}

	return p.Value.([]string)
}

func (c *Config) StringLists(name string) [][]string {
	return values[[]string](properties(c.Properties, name))
}

// Block returns block by name or nil if no such block found.
func (c *Config) Block(name string) *Block {
	for _, b := range c.Blocks {
		if b.Name == name {
			return b
		}
	}

	return nil
}

func values[T any](props []*Property) []T {
	var vs []T
	for _, p := range props {
		vs = append(vs, p.Value.(T))
	}

	return vs
}

func property(props []*Property, name string) *Property {
	ps := properties(props, name)
	if len(ps) == 0 {
		return nil
	}

	return ps[0]
}

func properties(props []*Property, name string) []*Property {
	var ps []*Property

	for _, p := range props {
		if p.Name == name {
			ps = append(ps, p)
		}
	}

	return ps
}
