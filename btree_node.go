package btree

type node[T any] struct {
	tree *BTreeG[T]

	// len(children) = len(items) + 1
	items    []T
	children *[]*node[T]
}

func (n *node[T]) leaf() bool {
	return n.children == nil
}

func (n *node[T]) binarySearch(key T) (int, bool) {
	low, high := 0, len(n.items)
	for low < high {
		mid := (low + high) / 2
		if !n.tree.less(key, n.items[mid]) {
			low = mid + 1
		} else {
			high = mid
		}
	}
	// now, node.items[low] is the first element that greater than key
	if low > 0 && !n.tree.less(n.items[low-1], key) {
		return low - 1, true
	}
	return low, false
}

func (n *node[T]) set(item T) (prev T, replaced, split bool) {
	i, found := n.binarySearch(item)
	if found {
		prev = n.items[i]
		n.items[i] = item
		return prev, true, false
	}
	if n.leaf() {
		if len(n.items) == n.tree.max {
			return n.tree.empty, false, true
		}
		n.items = append(n.items, n.tree.empty)
		copy(n.items[i+1:], n.items[i:])
		n.items[i] = item
		return n.tree.empty, false, false
	}
	prev, replaced, split = (*n.children)[i].set(item)
	if split {
		if len(n.items) == n.tree.max {
			return n.tree.empty, false, true
		}
		right, median := (*n.children)[i].split()
		n.items = append(n.items, n.tree.empty)
		copy(n.items[i+1:], n.items[i:])
		n.items[i] = median
		*n.children = append(*n.children, nil)
		copy((*n.children)[i+1:], (*n.children)[i:])
		(*n.children)[i+1] = right
		return n.set(item)
	}
	return prev, replaced, false
}

func (n *node[T]) split() (*node[T], T) {
	i := n.tree.max / 2 // degree - 1
	median := n.items[i]

	right := n.tree.newNode(n.leaf())
	right.items = n.items[i+1:] // degree, ... 2*degree-2

	n.items[i] = n.tree.empty
	n.items = n.items[:i:i] // resize node.items, 0, ... degree-2

	if !n.leaf() {
		*right.children = (*n.children)[i+1:]    // degree, ... 2*dgree-1
		*n.children = (*n.children)[: i+1 : i+1] // 0, ... degree-1
	}
	return right, median
}

func (n *node[T]) delete(key T) (T, bool) {
	var prev T
	var deleted bool
	i, found := n.binarySearch(key)
	if n.leaf() {
		if found {
			prev = n.items[i]
			copy(n.items[i:], n.items[i+1:])
			n.items[len(n.items)-1] = n.tree.empty
			n.items = n.items[:len(n.items)-1]
			return prev, true
		}
		return n.tree.empty, false
	}
	if found {
		prev = n.items[i]
		maxLeftItem, _ := (*n.children)[i].deleteMax()
		deleted = true
		n.items[i] = maxLeftItem

	} else {
		prev, deleted = (*n.children)[i].delete(key)
	}
	if !deleted {
		return n.tree.empty, false
	}
	if len((*n.children)[i].items) < n.tree.min {
		n.rebalance(i)
	}
	if len(n.items) == 0 && !n.leaf() {
		n.items = (*n.children)[0].items
		n.children = (*n.children)[0].children
	}
	return prev, true
}

func (n *node[T]) deleteMax() (prev T, deleted bool) {
	i := len(n.items) - 1
	if n.leaf() {
		prev = n.items[i]
		n.items[i] = n.tree.empty
		n.items = n.items[:i]
		return prev, true
	}
	return (*n.children)[i+1].deleteMax()
}

func (n *node[T]) rebalance(i int) {
	if i == len(n.items) {
		i--
	}
	left, right := (*n.children)[i], (*n.children)[i+1]
	if len(left.items)+len(right.items) < n.tree.max {
		left.items = append(left.items, n.items[i])
		left.items = append(left.items, right.items...)
		if !left.leaf() {
			*left.children = append(*left.children, *right.children...)
		}

		copy(n.items[i:], n.items[i+1:])
		n.items[len(n.items)-1] = n.tree.empty
		n.items = n.items[:len(n.items)-1]

		copy((*n.children)[i+1:], (*n.children)[i+2:])
		(*n.children)[len(*n.children)-1] = nil
		*n.children = (*n.children)[:len(*n.children)-1]
	} else if len(left.items) > len(right.items) {
		// move item from left to right
		right.items = append(right.items, n.tree.empty)
		copy(right.items[1:], right.items)
		right.items[0] = n.items[i]
		n.items[i] = left.items[len(left.items)-1]
		left.items[len(left.items)-1] = n.tree.empty
		left.items = left.items[:len(left.items)-1]

		if !left.leaf() {
			*right.children = append(*right.children, nil)
			copy((*right.children)[1:], *right.children)
			(*right.children)[0] = (*left.children)[len(*left.children)-1]
			(*left.children)[len(*left.children)-1] = nil
			*left.children = (*left.children)[:len(*left.children)-1]
		}
	} else {
		left.items = append(left.items, n.items[i])
		n.items[i] = right.items[0]
		copy(right.items, right.items[1:])
		right.items[len(right.items)-1] = n.tree.empty
		right.items = right.items[:len(right.items)-1]

		if !left.leaf() {
			*left.children = append(*left.children, (*right.children)[0])
			copy(*right.children, (*right.children)[1:])
			(*right.children)[len(*right.children)-1] = nil
			*right.children = (*right.children)[:len(*right.children)-1]
		}
	}
}
