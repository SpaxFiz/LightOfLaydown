// @program: LightOfLaydown
// @author: Fizzy
// @created: 2021-11-25

package storage

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

type testType struct {
	Age *int
}

var incrNum int
var once sync.Once

func incr() {
	incrNum++
}

func testFn() (interface{}, error) {
	once.Do(incr)

	a := 4
	return &testType{Age: &a}, nil
}

func TestCacheGetWithFallback(t *testing.T) {
	var a, b *testType

	err := GetCache().LoadOrDo("test", &a, testFn)
	assert.NoError(t, err)
	assert.Equal(t, *a.Age, 4)

	err = GetCache().LoadOrDo("test", &b, testFn)
	assert.NoError(t, err)
	assert.Equal(t, *b.Age, 4)
	assert.Equal(t, 1, incrNum)
}
