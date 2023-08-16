package btree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func lessInt(a, b int) bool {
	return a < b
}

func lessFloat64(a, b float64) bool {
	return a < b
}

func Test_binarySearch(t *testing.T) {
	tree := NewBTreeGWithOption[int](lessInt, Options{Degree: 4}) // 3-4-5-6-7
	assert.Equal(t, 7, tree.max)
	assert.Equal(t, 3, tree.min)
	node := tree.newNode(true)
	node.items = []int{2, 6, 9, 9, 11, 11, 18}
	index, found := node.binarySearch(10)
	assert.False(t, found)
	assert.Equal(t, 4, index)

	index, found = node.binarySearch(9)
	assert.True(t, found)
	assert.Equal(t, 3, index)

	index, found = node.binarySearch(19)
	assert.False(t, found)
	assert.Equal(t, 7, index)

	index, found = node.binarySearch(1)
	assert.False(t, found)
	assert.Equal(t, 0, index)

	index, found = node.binarySearch(2)
	assert.True(t, found)
	assert.Equal(t, 0, index)
}

func Test_nodeSplit(t *testing.T) {
	tree := NewBTreeGWithOption[float64](lessFloat64, Options{Degree: 4}) // 3-4-5-6-7
	assert.Equal(t, 7, tree.max)
	assert.Equal(t, 3, tree.min)
	n := tree.newNode(false)
	n.items = []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0}
	*n.children = make([]*node[float64], 0, tree.max+1)
	for i := 0; i < 8; i++ {
		(*n.children) = append((*n.children), &node[float64]{items: []float64{float64(i) + 0.1}})
	}

	right, median := n.split()
	assert.Equal(t, 4.0, median)
	assert.Equal(t, []float64{1.0, 2.0, 3.0}, n.items)
	assert.Equal(t, []float64{5.0, 6.0, 7.0}, right.items)
	assert.Equal(t, []float64{5.1}, (*right.children)[1].items)
}
