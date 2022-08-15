package sql

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPaging(t *testing.T) {
	type expected struct {
		args     []uint
		lastPage uint
	}

	tests := []struct {
		input Paging
		exp   expected
	}{
		{
			input: Paging{1, 10, 100},
			exp:   expected{[]uint{10, 0}, 10},
		},
		{
			input: Paging{3, 10, 100},
			exp:   expected{[]uint{10, 20}, 10},
		},
		{
			input: Paging{1, 3, 100},
			exp:   expected{[]uint{3, 0}, 34},
		},
		{
			input: Paging{6, 3, 100},
			exp:   expected{[]uint{3, 15}, 34},
		},
	}

	t.Run("Test NewPaging", func(t *testing.T) {
		p := NewPaging(10)
		assert.Equal(t, p, &Paging{1, 12, 10})
	})

	t.Run("Test Normalize", func(t *testing.T) {
		p := Paging{}
		p.Normalize()
		assert.Equal(t, p, Paging{1, 12, 0})
	})

	for _, tt := range tests {
		testname := fmt.Sprintf("%+v", tt.input)
		t.Run(testname, func(t *testing.T) {
			limitClause, limitClauseArgs := tt.input.BuildQuery()
			assert.Equal(t, limitClause, "LIMIT ? OFFSET ?")
			assert.Equal(t, tt.exp, expected{limitClauseArgs, tt.input.LastPage()})
		})
	}
}
