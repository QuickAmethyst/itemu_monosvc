package sql

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestColumnName(t *testing.T) {
	tests := []struct {
		input string
		expected string
	} {
		{"Foo", "foo"},
		{"FooBar", "foo_bar"},
		{"ID", "id"},
		{"FooID", "foo_id"},
		{"FooIDBar", "foo_id_bar"},
		{"IDFoo", "id_foo"},
		{"FooiD", "fooi_d"},
		{"FooNotEQ", "foo"},
		{"FooLike", "foo"},
		{"FooGTE", "foo"},
		{"FooLTE", "foo"},
		{"FooNotIN", "foo"},
		{"FooIN", "foo"},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.input)
		t.Run(testname, func(t *testing.T) {
			f := FieldStrategy(tt.input)
			assert.Equal(t, tt.expected, f.ColumnName())
		})
	}
}

func TestIsEqualStatement(t *testing.T) {
	tests := []struct {
		input string
		expected bool
	} {
		{"Foo", true},
		{"FooNotEQ", false},
		{"FooLike", false},
		{"FooGTE", false},
		{"FooLTE", false},
		{"FooNotIN", false},
		{"FooIN", false},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.input)
		t.Run(testname, func(t *testing.T) {
			f := FieldStrategy(tt.input)
			assert.Equal(t, tt.expected, f.IsEqualStatement())
		})
	}
}

func TestIsNotEqualStatement(t *testing.T) {
	tests := []struct {
		input string
		expected bool
	} {
		{"Foo", false},
		{"FooNotEQ", true},
		{"FooLike", false},
		{"FooGTE", false},
		{"FooLTE", false},
		{"FooNotIN", false},
		{"FooIN", false},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.input)
		t.Run(testname, func(t *testing.T) {
			f := FieldStrategy(tt.input)
			assert.Equal(t, tt.expected, f.IsNotEqualStatement())
		})
	}
}

func TestIsLikeStatement(t *testing.T) {
	tests := []struct {
		input string
		expected bool
	} {
		{"Foo", false},
		{"FooNotEQ", false},
		{"FooLike", true},
		{"FooGTE", false},
		{"FooLTE", false},
		{"FooNotIN", false},
		{"FooIN", false},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.input)
		t.Run(testname, func(t *testing.T) {
			f := FieldStrategy(tt.input)
			assert.Equal(t, tt.expected, f.IsLikeStatement())
		})
	}
}

func TestIsGreaterThanEqualStatement(t *testing.T) {
	tests := []struct {
		input string
		expected bool
	} {
		{"Foo", false},
		{"FooNotEQ", false},
		{"FooLike", false},
		{"FooGTE", true},
		{"FooLTE", false},
		{"FooNotIN", false},
		{"FooIN", false},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.input)
		t.Run(testname, func(t *testing.T) {
			f := FieldStrategy(tt.input)
			assert.Equal(t, tt.expected, f.IsGreaterThanEqualStatement())
		})
	}
}

func TestIsLessThanEqualStatement(t *testing.T) {
	tests := []struct {
		input string
		expected bool
	} {
		{"Foo", false},
		{"FooNotEQ", false},
		{"FooLike", false},
		{"FooGTE", false},
		{"FooLTE", true},
		{"FooNotIN", false},
		{"FooIN", false},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.input)
		t.Run(testname, func(t *testing.T) {
			f := FieldStrategy(tt.input)
			assert.Equal(t, tt.expected, f.IsLessThanEqualStatement())
		})
	}
}

func TestIsInStatement(t *testing.T) {
	tests := []struct {
		input string
		expected bool
	} {
		{"Foo", false},
		{"FooNotEQ", false},
		{"FooLike", false},
		{"FooGTE", false},
		{"FooLTE", false},
		{"FooNotIN", false},
		{"FooIN", true},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.input)
		t.Run(testname, func(t *testing.T) {
			f := FieldStrategy(tt.input)
			assert.Equal(t, tt.expected, f.IsInStatement())
		})
	}
}

func TestIsNotInStatement(t *testing.T) {
	tests := []struct {
		input string
		expected bool
	} {
		{"Foo", false},
		{"FooNotEQ", false},
		{"FooLike", false},
		{"FooGTE", false},
		{"FooLTE", false},
		{"FooNotIN", true},
		{"FooIN", false},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.input)
		t.Run(testname, func(t *testing.T) {
			f := FieldStrategy(tt.input)
			assert.Equal(t, tt.expected, f.IsNotInStatement())
		})
	}
}
