package btree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Set_Get(t *testing.T) {
	tree := NewBTreeGWithOption[float64](lessFloat64, Options{Degree: 2})
	tree.Set(1.0)
	tree.Set(2.0)
	tree.Set(3.0)

	assert.Equal(t, []float64{1.0, 2.0, 3.0}, tree.root.items)

	tree.Set(4.0)
	assert.Equal(t, []float64{2.0}, tree.root.items)
	assert.Equal(t, []float64{1.0}, (*tree.root.children)[0].items)
	assert.Equal(t, []float64{3.0, 4.0}, (*tree.root.children)[1].items)

	ret, ok := tree.Get(5.0)
	assert.False(t, ok)
	assert.Equal(t, 0.0, ret)

	ret, ok = tree.Get(4.0)
	assert.True(t, ok)
	assert.Equal(t, 4.0, ret)

	tree.Set(5.0)
	tree.Set(6.0)
	tree.Set(2.1)
	tree.Set(2.2)
	assert.Equal(t, []float64{2.0, 4.0}, tree.root.items)
	assert.Equal(t, []float64{1.0}, (*tree.root.children)[0].items)
	assert.Equal(t, []float64{2.1, 2.2, 3.0}, (*tree.root.children)[1].items)
	assert.Equal(t, []float64{5.0, 6.0}, (*tree.root.children)[2].items)

	tree.Set(3.1)
	assert.Equal(t, []float64{2.0, 2.2, 4.0}, tree.root.items)
	assert.Equal(t, []float64{1.0}, (*tree.root.children)[0].items)
	assert.Equal(t, []float64{2.1}, (*tree.root.children)[1].items)
	assert.Equal(t, []float64{3.0, 3.1}, (*tree.root.children)[2].items)
	assert.Equal(t, []float64{5.0, 6.0}, (*tree.root.children)[3].items)

	tree.Set(7.0)
	tree.Set(8.0)
	assert.Equal(t, []float64{2.2}, tree.root.items)
	assert.Equal(t, []float64{2.0}, (*tree.root.children)[0].items)
	assert.Equal(t, []float64{4.0, 6.0}, (*tree.root.children)[1].items)
	assert.Equal(t, []float64{7.0, 8.0}, (*(*tree.root.children)[1].children)[2].items)
}
