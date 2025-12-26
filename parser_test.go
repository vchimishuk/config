package config

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		Input string
		Spec  *Spec
		Exp   *Config
	}{
		{
			"name = \"foo\", \"bar\", \"baz\"",
			&Spec{
				[]*PropertySpec{
					&PropertySpec{TypeStringList, "name", false, false, nil},
				},
				nil,
				true,
			},
			&Config{
				[]*Property{
					&Property{TypeStringList, "name", []string{"foo", "bar", "baz"}},
				},
				nil,
			},
		},
		{
			"foo = 1; bar = 2;",
			&Spec{
				[]*PropertySpec{
					&PropertySpec{TypeInt, "foo", false, true, nil},
					&PropertySpec{TypeInt, "bar", false, true, nil},
				},
				nil,
				true,
			},
			&Config{
				[]*Property{
					&Property{TypeInt, "foo", 1},
					&Property{TypeInt, "bar", 2},
				},
				nil,
			},
		},
		{
			"foo = 123\nbar = \"value\"",
			&Spec{
				[]*PropertySpec{
					&PropertySpec{TypeInt, "foo", false, true, nil},
					&PropertySpec{TypeString, "bar", false, true, nil},
				},
				nil,
				true,
			},
			&Config{
				[]*Property{
					&Property{TypeInt, "foo", 123},
					&Property{TypeString, "bar", "value"},
				},
				nil,
			},
		},
		{
			"foo = 1; bar { baz = 2; qux = 3; }",
			&Spec{
				[]*PropertySpec{
					&PropertySpec{TypeInt, "foo", false, true, nil},
				},
				[]*BlockSpec{
					&BlockSpec{
						"bar",
						false,
						false,
						[]*PropertySpec{
							&PropertySpec{TypeInt, "baz", false, true, nil},
							&PropertySpec{TypeInt, "qux", false, true, nil},
						},
						nil,
						true,
					},
				},
				true,
			},
			&Config{
				[]*Property{
					&Property{TypeInt, "foo", 1},
				},
				[]*Block{
					&Block{
						"bar",
						[]*Property{
							&Property{TypeInt, "baz", 2},
							&Property{TypeInt, "qux", 3},
						},
						nil,
					},
				},
			},
		},
		{
			"foo { foo-prop = 1; bar { bar-prop = 2; baz { baz-prop = 3; } qux { qux-prop = 4; } } }",
			&Spec{
				nil,
				[]*BlockSpec{
					&BlockSpec{
						"foo",
						false,
						false,
						[]*PropertySpec{
							&PropertySpec{TypeInt, "foo-prop", false, true, nil},
						},
						[]*BlockSpec{
							&BlockSpec{
								"bar",
								false,
								false,
								[]*PropertySpec{
									&PropertySpec{TypeInt, "bar-prop", false, true, nil},
								},
								[]*BlockSpec{
									&BlockSpec{
										"baz",
										false,
										false,
										[]*PropertySpec{
											&PropertySpec{TypeInt, "baz-prop", false, true, nil},
										},
										nil,
										true,
									},
									&BlockSpec{
										"qux",
										false,
										false,
										[]*PropertySpec{
											&PropertySpec{TypeInt, "qux-prop", false, true, nil},
										},
										nil,
										true,
									},
								},
								true,
							},
						},
						true,
					},
				},
				true,
			},
			&Config{
				nil,
				[]*Block{
					&Block{
						"foo",
						[]*Property{
							&Property{TypeInt, "foo-prop", 1},
						},
						[]*Block{
							&Block{
								"bar",
								[]*Property{
									&Property{TypeInt, "bar-prop", 2},
								},
								[]*Block{
									&Block{
										"baz",
										[]*Property{
											&Property{TypeInt, "baz-prop", 3},
										},
										nil,
									},
									&Block{
										"qux",
										[]*Property{
											&Property{TypeInt, "qux-prop", 4},
										},
										nil,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		testParse(t, tc.Input, tc.Spec, tc.Exp)
	}
}

func TestParseRepeatProperty(t *testing.T) {
	// Repeated property is allowed by the spec.
	testParse(t,
		"foo = 1; foo = 2;",
		&Spec{
			[]*PropertySpec{
				&PropertySpec{TypeInt, "foo", true, false, nil},
			},
			nil,
			true,
		},
		&Config{
			[]*Property{
				&Property{TypeInt, "foo", 1},
				&Property{TypeInt, "foo", 2},
			},
			nil,
		},
	)

	// Repeated property is not allowed.
	_, err := Parse(
		&Spec{
			[]*PropertySpec{
				&PropertySpec{TypeInt, "foo", false, false, nil},
			},
			nil,
			true,
		},
		"foo = 1; foo = 2;",
	)
	if err == nil || err.Error() != "1: property `foo` already defined" {
		t.Fatal(err)
	}
}

func TestParseRequireProperty(t *testing.T) {
	// Property is not required.
	cfg, err := Parse(
		&Spec{
			[]*PropertySpec{
				&PropertySpec{TypeInt, "foo", false, false, nil},
			},
			nil,
			true,
		},
		"",
	)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(&Config{}, cfg) {
		t.Fatal(cfg)
	}

	// Property is required.
	cfg, err = Parse(
		&Spec{
			[]*PropertySpec{
				&PropertySpec{TypeInt, "foo", false, true, nil},
			},
			nil,
			true,
		},
		"",
	)
	if err == nil || err.Error() != "1: missing required property `foo`" {
		t.Fatal(err)
	}
}

func TestParseStarProperty(t *testing.T) {
	testParse(t,
		"",
		&Spec{
			[]*PropertySpec{
				&PropertySpec{TypeInt, "foo",
					false, false, nil},
				&PropertySpec{TypeInt, "foo.*",
					false, false, nil},
				&PropertySpec{TypeInt, "foo.bar.*",
					false, false, nil},
			},
			nil,
			true,
		},
		&Config{
			nil,
			nil,
		},
	)

	testParse(t,
		"foo = 1; foo.baz = true; foo.bar.baz = \"str\";",
		&Spec{
			[]*PropertySpec{
				&PropertySpec{TypeInt, "foo",
					false, false, nil},
				&PropertySpec{TypeBool, "foo.*",
					false, false, nil},
				&PropertySpec{TypeString, "foo.bar*",
					false, false, nil},
			},
			nil,
			true,
		},
		&Config{
			[]*Property{
				&Property{TypeInt, "foo", 1},
				&Property{TypeBool, "foo.baz", true},
				&Property{TypeString, "foo.bar.baz", "str"},
			},
			nil,
		},
	)
}

func TestParseRepeatBlock(t *testing.T) {
	// Repeated block is allowed by the spec.
	testParse(t,
		"foo {}; foo {};",
		&Spec{
			nil,
			[]*BlockSpec{
				&BlockSpec{"foo", true, true, nil, nil, true},
			},
			true,
		},
		&Config{
			nil,
			[]*Block{
				&Block{"foo", nil, nil},
				&Block{"foo", nil, nil},
			},
		},
	)

	// Repeated block is not allowed.
	_, err := Parse(
		&Spec{
			nil,
			[]*BlockSpec{
				&BlockSpec{"foo", false, true, nil, nil, true},
			},
			true,
		},
		"foo {}; foo {};",
	)
	if err == nil || err.Error() != "1: block `foo` already defined" {
		t.Fatal(err)
	}
}

func TestParseRequireBlock(t *testing.T) {
	// Block is not required.
	cfg, err := Parse(
		&Spec{
			nil,
			[]*BlockSpec{
				&BlockSpec{"foo", false, false, nil, nil, true},
			},
			true,
		},
		"",
	)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(&Config{}, cfg) {
		t.Fatal(cfg)
	}

	// Property is required.
	cfg, err = Parse(
		&Spec{
			nil,
			[]*BlockSpec{
				&BlockSpec{"foo", false, true, nil, nil, true},
			},
			true,
		},
		"",
	)
	if err == nil || err.Error() != "1: missing required block `foo`" {
		t.Fatal(err)
	}
}

func TestParseStarBlock(t *testing.T) {
	cfg, err := Parse(
		&Spec{
			nil,
			[]*BlockSpec{
				&BlockSpec{
					"*",
					true,
					false,
					[]*PropertySpec{
						&PropertySpec{TypeInt, "prop", false, false, nil},
					},
					nil,
					true,
				},
			},
			true,
		},
		"foo { prop = 1 } bar { prop = 2 } baz { prop = 3 }",
	)
	if err != nil {
		t.Fatal(err)
	}
	exp := &Config{
		nil,
		[]*Block{
			&Block{
				"foo",
				[]*Property{
					&Property{TypeInt, "prop", 1},
				},
				nil,
			},
			&Block{
				"bar",
				[]*Property{
					&Property{TypeInt, "prop", 2},
				},
				nil,
			},
			&Block{
				"baz",
				[]*Property{
					&Property{TypeInt, "prop", 3},
				},
				nil,
			},
		},
	}
	if !reflect.DeepEqual(exp, cfg) {
		t.Fatal(cfg)
	}
}

func TestParsePropertyType(t *testing.T) {
	spec := &Spec{
		[]*PropertySpec{
			&PropertySpec{TypeString, "foo", false, false, nil},
		},
		nil,
		true,
	}
	_, err := Parse(spec, "foo = 1")
	if err == nil || err.Error() != "1: string value expected" {
		t.Fatal(err)
	}
	_, err = Parse(spec, "foo = false")
	if err == nil || err.Error() != "1: string value expected" {
		t.Fatal(err)
	}
	_, err = Parse(spec, "foo = \"bar\"")
	if err != nil {
		t.Fatal(err)
	}
}

func TestParseString(t *testing.T) {
	spec := &Spec{
		[]*PropertySpec{
			&PropertySpec{TypeString, "foo", false, false, nil},
			&PropertySpec{TypeString, "bar", true, false, nil},
		},
		nil,
		true,
	}

	cfg := testParse(t, `bar = "one"; foo = "two"; bar = "three"`, spec,
		&Config{[]*Property{
			&Property{TypeString, "bar", "one"},
			&Property{TypeString, "foo", "two"},
			&Property{TypeString, "bar", "three"},
		}, nil})
	assert(t, "two", cfg.String("foo"))
	assert(t, []string{"one", "three"}, cfg.Strings("bar"))
}

func TestParseBool(t *testing.T) {
	spec := &Spec{
		[]*PropertySpec{
			&PropertySpec{TypeBool, "foo", false, false, nil},
			&PropertySpec{TypeBool, "bar", true, false, nil},
		},
		nil,
		true,
	}
	testParse(t, "bar = true; foo = true; bar = false", spec,
		&Config{[]*Property{
			&Property{TypeBool, "bar", true},
			&Property{TypeBool, "foo", true},
			&Property{TypeBool, "bar", false},
		}, nil})
	_, err := Parse(spec, "foo = bar")
	if err == nil || err.Error() != "1: invalid boolean value" {
		t.Fatal(err)
	}
}

func TestParseDuration(t *testing.T) {
	spec := &Spec{
		[]*PropertySpec{
			&PropertySpec{TypeDuration, "foo", false, false, nil},
		},
		nil,
		true,
	}
	d, _ := time.ParseDuration("1s")
	testParse(t, "foo = 1s", spec,
		&Config{[]*Property{&Property{TypeDuration, "foo", d}}, nil})
	d, _ = time.ParseDuration("1h30m")
	testParse(t, "foo = 1h30m", spec,
		&Config{[]*Property{&Property{TypeDuration, "foo", d}}, nil})
	d, _ = time.ParseDuration("1.5m")
	testParse(t, "foo = 1.5m", spec,
		&Config{[]*Property{&Property{TypeDuration, "foo", d}}, nil})
}

func TestParseComment(t *testing.T) {
	spec := &Spec{
		[]*PropertySpec{
			&PropertySpec{TypeDuration, "heartbeat-ttl", true, true, nil},
		},
		nil,
		true,
	}
	s := "heartbeat-ttl = 3s\n\n# comment\nheartbeat-ttl = 6s # more comment\n"
	testParse(t, s, spec,
		&Config{[]*Property{
			&Property{TypeDuration, "heartbeat-ttl", time.Second * 3},
			&Property{TypeDuration, "heartbeat-ttl", time.Second * 6},
		},
			nil})

	s = "heartbeat-ttl = 3s\n\n# block {\n#}\n"
	testParse(t, s, spec,
		&Config{[]*Property{
			&Property{TypeDuration, "heartbeat-ttl", time.Second * 3},
		},
			nil})
}

func TestParseStrict(t *testing.T) {
	spec := &Spec{
		[]*PropertySpec{
			&PropertySpec{TypeInt, "foo", false, false, nil},
		},
		nil,
		true,
	}
	_, err := Parse(spec, "bar = 1\nbaz = 2")
	exp := "1: unsupported property: bar"
	if err.Error() != exp {
		t.Fatalf("%s != %s", err.Error(), exp)
	}
}

func TestParseNonStrict(t *testing.T) {
	spec := &Spec{
		[]*PropertySpec{
			&PropertySpec{TypeInt, "foo", false, false, nil},
		},
		nil,
		false,
	}
	_, err := Parse(spec, "bar = 1\nbaz = 2")
	if err != nil {
		t.Fatal(err)
	}
}

func TestParseParser(t *testing.T) {
	parser := func(v any) (any, error) {
		s := v.(string)

		switch s {
		case "true":
			return true, nil
		case "false":
			return false, nil
		default:
			return false, errors.New("unsupported value")
		}
	}

	spec := &Spec{
		[]*PropertySpec{
			&PropertySpec{
				Type:    TypeString,
				Name:    "prop",
				Repeat:  false,
				Require: false,
				Parser:  parser,
			},
		},
		nil,
		false,
	}

	c, err := Parse(spec, "prop = \"true\"")
	if err != nil {
		t.Fatal(err)
	}
	if c.Any("prop") != true {
		t.Fatal()
	}

	c, err = Parse(spec, "prop = \"false\"")
	if err != nil {
		t.Fatal(err)
	}
	if c.Any("prop") != false {
		t.Fatal()
	}

	_, err = Parse(spec, "prop = \"invalid\"")
	if err == nil || err.Error() != "1: unsupported value" {
		t.Fatal(err)
	}
}

func testParse(t *testing.T, s string, spec *Spec, exp *Config) *Config {
	act, err := Parse(spec, s)
	if err != nil {
		t.Fatal(err)
	}
	assert(t, exp, act)

	return act
}

func assert(t *testing.T, exp any, act any) {
	if !reflect.DeepEqual(exp, act) {
		t.Fatalf("%+v != %+v", exp, act)
	}
}
