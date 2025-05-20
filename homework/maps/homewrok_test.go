package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Tree struct {
	node, value int
	left, right *Tree
}
type OrderedMap struct {
	tree Tree
	size int
	keys []int
}

func (tree *Tree) insertInTree(key, value int) *Tree {
	if key < tree.node {
		if tree.left == nil {
			tree.left = &Tree{
				node:  key,
				value: value,
				right: nil,
				left:  nil,
			}
		} else {
			tree.left = tree.left.insertInTree(key, value)
		}
	} else if key > tree.node {
		if tree.right == nil {
			tree.right = &Tree{
				node:  key,
				value: value,
				right: nil,
				left:  nil,
			}
		} else {
			tree.right = tree.right.insertInTree(key, value)
		}
	}
	return tree
}

func (tree *Tree) deleteNode(key int) *Tree {
	if tree == nil {
		return nil
	}
	switch {
	case key < tree.node:
		tree.left = tree.left.deleteNode(key)
	case key > tree.node:
		tree.right = tree.right.deleteNode(key)
	default:
		if tree.left == nil {
			return tree.right
		}
		if tree.right == nil {
			return tree.left
		}

		succ := tree.right
		for succ.left != nil {
			succ = succ.left
		}

		tree.node, tree.value = succ.node, succ.value

		tree.right = tree.right.deleteNode(succ.node)
	}
	return tree
}

func (tree *Tree) changeTree(action func(int, int)) {
	if tree == nil {
		return
	}
	tree.left.changeTree(action)
	action(tree.node, tree.value)
	tree.right.changeTree(action)
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{
		tree: Tree{
			node:  0,
			value: 0,
			left:  nil,
			right: nil,
		},
		size: 0,
		keys: make([]int, 0),
	}
}

func (m *OrderedMap) Insert(key, value int) {
	if m.size == 0 {
		m.tree.node = key
		m.tree.value = value
		m.size++
		m.keys = append(m.keys, key)
		return
	}
	if !m.Contains(key) {
		m.tree = *m.tree.insertInTree(key, value)
		m.keys = append(m.keys, key)
		m.size++
	}
}

func (m *OrderedMap) Erase(key int) {
	if m.Contains(key) {
		m.tree = *m.tree.deleteNode(key)
		var needID int
		for i := range m.keys {
			if m.keys[i] == key {
				needID = i
				break
			}
		}
		m.keys = append(m.keys[:needID], m.keys[needID+1:]...)
		m.size--
	}
}

func (m *OrderedMap) Contains(key int) bool {
	for i := range m.keys {
		if m.keys[i] == key {
			return true
		}
	}
	return false
}

func (m *OrderedMap) Size() int {
	return m.size
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	m.tree.changeTree(action)
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
