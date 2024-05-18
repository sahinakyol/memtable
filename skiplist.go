package main

import (
	"fmt"
	"math/rand"
)

const MaxLevel = 16

type Node struct {
	key   int
	value interface{}
	level []*Node
}

type SkipList struct {
	head  *Node
	level int
}

func NewNode(key int, value interface{}, level int) *Node {
	return &Node{
		key:   key,
		value: value,
		level: make([]*Node, level),
	}
}

func NewSkipList() *SkipList {
	head := NewNode(-1, nil, MaxLevel)
	return &SkipList{
		head:  head,
		level: 1,
	}
}

func randomLevel() int {
	level := 1
	for rand.Float32() < 0.5 && level < MaxLevel {
		level++
	}
	return level
}

func (sl *SkipList) Insert(key int, value interface{}) {
	update := make([]*Node, MaxLevel)
	current := sl.head

	for i := sl.level - 1; i >= 0; i-- {
		for current.level[i] != nil && current.level[i].key < key {
			current = current.level[i]
		}
		update[i] = current
	}

	current = current.level[0]

	if current == nil || current.key != key {
		newLevel := randomLevel()

		if newLevel > sl.level {
			for i := sl.level; i < newLevel; i++ {
				update[i] = sl.head
			}
			sl.level = newLevel
		}

		newNode := NewNode(key, value, newLevel)

		for i := 0; i < newLevel; i++ {
			newNode.level[i] = update[i].level[i]
			update[i].level[i] = newNode
		}
	} else {
		current.value = value
	}
}

func (sl *SkipList) Search(key int) *Node {
	current := sl.head

	for i := sl.level - 1; i >= 0; i-- {
		for current.level[i] != nil && current.level[i].key < key {
			current = current.level[i]
		}
	}

	current = current.level[0]

	if current != nil && current.key == key {
		return current
	}
	return nil
}

func (sl *SkipList) Delete(key int) {
	update := make([]*Node, MaxLevel)
	current := sl.head

	for i := sl.level - 1; i >= 0; i-- {
		for current.level[i] != nil && current.level[i].key < key {
			current = current.level[i]
		}
		update[i] = current
	}

	current = current.level[0]

	if current != nil && current.key == key {
		for i := 0; i < sl.level; i++ {
			if update[i].level[i] != current {
				break
			}
			update[i].level[i] = current.level[i]
		}

		for sl.level > 1 && sl.head.level[sl.level-1] == nil {
			sl.level--
		}
	}
}

func (sl *SkipList) Display() {
	for i := 0; i < sl.level; i++ {
		current := sl.head.level[i]
		fmt.Printf("Level %d: ", i)
		for current != nil {
			fmt.Printf("%d:%v ", current.key, current.value)
			current = current.level[i]
		}
		fmt.Println()
	}
}
