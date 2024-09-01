package async_test

import (
	"iter"
	"math/rand/v2"
	"slices"
	"testing"
	"time"

	"github.com/mackee/iterutils/async"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	ss := slices.Repeat([]string{"a", "b", "c", "d", "e"}, 1000)
	it := slices.Values(ss)
	it2 := async.Map(it, func(s string) string {
		time.Sleep(time.Duration(rand.IntN(11)) * time.Millisecond * 100)
		return s
	})
	t1 := time.Now()
	ret := slices.Collect(it2)
	assert.Equal(t, ss, ret)
	assert.True(t, time.Since(t1) < 1100*time.Millisecond)
}

func TestMap2(t *testing.T) {
	ss := []string{"a", "b", "c", "d", "e"}
	it := slices.All(ss)
	it2 := async.Map2(it, func(i int, s string) (int, string) {
		time.Sleep(1 * time.Second)
		return i, s
	})
	t1 := time.Now()
	ret := collect2(it2)
	assert.Equal(t, ss, ret)
	assert.True(t, time.Since(t1) < 1100*time.Millisecond)
}

func collect2[U any, V any](it iter.Seq2[U, V]) []V {
	var ret []V
	for _, v := range it {
		ret = append(ret, v)
	}
	return ret
}
