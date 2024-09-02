package async_test

import (
	"iter"
	"math/rand/v2"
	"slices"
	"testing"
	"time"

	"github.com/mackee/iterutils/async"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
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

func TestMapWithBreak(t *testing.T) {
	defer goleak.VerifyNone(t)

	ss := slices.Repeat([]string{"a", "b", "c", "d", "e"}, 1000)
	it := slices.Values(ss)
	it2 := async.Map(it, func(s string) string {
		time.Sleep(time.Duration(rand.IntN(11)) * time.Millisecond * 100)
		return s
	})
	t1 := time.Now()
	ret := make([]string, 0, 100)
	for v := range it2 {
		ret = append(ret, v)
		if len(ret) >= 100 {
			break
		}
	}
	assert.Equal(t, ss[:100], ret)
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

func TestMap2WithBreak(t *testing.T) {
	defer goleak.VerifyNone(t)

	ss := []string{"a", "b", "c", "d", "e"}
	it := slices.All(ss)
	it2 := async.Map2(it, func(i int, s string) (int, string) {
		time.Sleep(1 * time.Second)
		return i, s
	})
	t1 := time.Now()
	ret := make([]string, 0, 3)
	for _, v := range it2 {
		ret = append(ret, v)
		if len(ret) >= 3 {
			break
		}
	}
	assert.Equal(t, ss[:3], ret)
	assert.True(t, time.Since(t1) < 3100*time.Millisecond)
}

func collect2[U any, V any](it iter.Seq2[U, V]) []V {
	var ret []V
	for _, v := range it {
		ret = append(ret, v)
	}
	return ret
}
