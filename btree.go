package btree

type BTreeG[T any] struct {
	degree int
	max    int // 2 * degree - 1
	min    int // dgress - 1
	root   *node[T]
	empty  T
	less   func(a, b T) bool
}

type Options struct {
	Degree int
}

func NewBTreeGWithOption[T any](less func(a, b T) bool, option Options) *BTreeG[T] {
	degree := option.Degree
	if degree <= 0 {
		degree = 32
	}
	return &BTreeG[T]{
		degree: degree,
		max:    2*degree - 1,
		min:    degree - 1,
		less:   less,
	}
}

func (tree *BTreeG[T]) newNode(leaf bool) *node[T] {
	n := &node[T]{
		tree: tree,
	}
	if !leaf {
		n.children = new([]*node[T])
	}
	return n
}

func (tree *BTreeG[T]) Set(item T) (prev T, replaced bool) {
	if tree.root == nil {
		tree.root = tree.newNode(true)
		tree.root.items = append([]T{}, item)
		return tree.empty, false
	}
	prev, replaced, split := tree.root.set(item)
	if split {
		left := tree.root
		right, median := left.split()
		tree.root = tree.newNode(false)
		tree.root.items = append([]T{}, median)
		*tree.root.children = make([]*node[T], 0, tree.max+1)
		*tree.root.children = append(*tree.root.children, left, right)
		return tree.Set(item)
	}
	if replaced {
		return prev, true
	}
	return tree.empty, false
}

func (tree *BTreeG[T]) Get(key T) (T, bool) {
	n := tree.root
	for {
		i, found := n.binarySearch(key)
		if found {
			return n.items[i], true
		}
		if n.children == nil {
			return tree.empty, false
		}
		n = (*n.children)[i]
	}
}

func (tree *BTreeG[T]) Delete(key T) (T, bool) {
	if tree.root == nil {
		return tree.empty, false
	}
	prev, deleted := tree.root.delete(key)
	if !deleted {
		return tree.empty, false
	}
	// image a 2-3-4 tree
	//   2                  2                               ?
	//  / \  delete 1 -->  / \  --> trigger rebalance -->  / \  only root --> 2,3
	// 1   3              ?   3                          2,3  ?
	if len(tree.root.items) == 0 && !tree.root.leaf() {
		tree.root = (*tree.root.children)[0]
	}
	return prev, deleted
}
