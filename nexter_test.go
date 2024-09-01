package iterutils_test

import (
	"slices"
	"testing"

	"github.com/mackee/iterutils"
	"github.com/stretchr/testify/assert"
)

func TestFromNexter(t *testing.T) {
	ss := []string{"a", "b", "c", "d", "e"}
	n := &nexter{
		ss: ss,
	}
	it := iterutils.FromNexter(n, func(n *nexter) string {
		return n.Value()
	})
	ret := slices.Collect(it)
	assert.Equal(t, ss, ret)
}

type nexter struct {
	ss    []string
	index int
}

func (n *nexter) Next() bool {
	if n.index < len(n.ss) {
		n.index++
		return true
	}
	return false
}

func (n *nexter) Value() string {
	return n.ss[n.index-1]
}
