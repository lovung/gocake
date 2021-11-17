package gocake

import (
	"container/list"
	"fmt"
	"strings"
	"sync"
)

// lfu saves the frequency of each entry.
type lfu struct {
	lock sync.Mutex

	// Linked list of the freqNode values.
	counterList *list.List // list.Element.Value is freqNode

	keyMap map[uint64]*list.Element // list.Element.Value is freqNode
}

type freqNode struct {
	// The number of times this freqNode has been accessed.
	value uint64

	keys map[uint64]struct{}
}

func newLFU() *lfu {
	return &lfu{
		counterList: list.New(),
		keyMap:      make(map[uint64]*list.Element),
	}
}

func (l *lfu) String() string {
	l.lock.Lock()
	defer l.lock.Unlock()
	b := strings.Builder{}
	for e := l.counterList.Front(); e != nil; e = e.Next() {
		b.WriteString(fmt.Sprintf("%d: ", e.Value.(*freqNode).value))
		for k := range e.Value.(*freqNode).keys {
			b.WriteString(fmt.Sprintf("%d ", k))
		}
		if e.Next() != nil {
			b.WriteString("- ")
		}
	}
	return b.String()
}

func (l *lfu) newItem(key uint64) *list.Element {
	l.lock.Lock()
	defer l.lock.Unlock()
	firstNode := l.counterList.Front()
	if firstNode == nil {
		// Create a new freqNode.
		node := &freqNode{
			value: 1,
			keys:  map[uint64]struct{}{key: {}},
		}
		e := l.counterList.PushFront(node)
		l.keyMap[key] = e
		return e
	}
	if firstNode.Value.(*freqNode).value == 1 {
		// The first freqNode is the only one with this value.
		firstNode.Value.(*freqNode).keys[key] = struct{}{}
		l.keyMap[key] = firstNode
		return firstNode
	}
	// Create a new freqNode.
	node := &freqNode{
		value: 1,
		keys:  map[uint64]struct{}{key: {}},
	}
	e := l.counterList.PushFront(node)
	l.keyMap[key] = e
	return e
}

func (l *lfu) touch(key uint64) *list.Element {
	l.lock.Lock()
	e := l.keyMap[key]
	if e == nil {
		l.lock.Unlock()
		return l.newItem(key)
	}
	defer l.lock.Unlock()
	node := e.Value.(*freqNode)
	if len(node.keys) == 1 {
		// Increase the freqNode.
		node.value++
		return e
	}
	delete(node.keys, key)
	nextNode := e.Next()
	if nextNode == nil {
		// Create a new freqNode.
		newNode := &freqNode{
			value: node.value + 1,
			keys:  map[uint64]struct{}{key: {}},
		}
		e := l.counterList.PushBack(newNode)
		l.keyMap[key] = e
		return e
	} else if nextNode.Value.(*freqNode).value == node.value+1 {
		// Move up
		nextNode.Value.(*freqNode).keys[key] = struct{}{}
		l.keyMap[key] = nextNode
		return nextNode
	}
	// Create a new freqNode.
	newNode := &freqNode{
		value: node.value + 1,
		keys:  map[uint64]struct{}{key: {}},
	}
	e = l.counterList.InsertBefore(newNode, nextNode)
	l.keyMap[key] = e
	return e
}

// clean removes the least frequently used items.
func (l *lfu) clean(quantity int) []uint64 {
	l.lock.Lock()
	defer l.lock.Unlock()
	var keys []uint64
	for i := 0; i < quantity; {
		node := l.counterList.Front()
		if node == nil {
			break
		}
		for k := range node.Value.(*freqNode).keys {
			delete(node.Value.(*freqNode).keys, k)
			delete(l.keyMap, k)
			keys = append(keys, k)
			i++
			if i == quantity {
				if len(node.Value.(*freqNode).keys) == 0 {
					// Remove the freqNode.
					l.counterList.Remove(node)
				}
				return keys
			}
		}
		// Remove the freqNode.
		l.counterList.Remove(node)
	}
	return keys
}

// del a key in lfu
func (l *lfu) del(key uint64) {
	l.lock.Lock()
	defer l.lock.Unlock()
	e := l.keyMap[key]
	if e == nil {
		return
	}
	node := e.Value.(*freqNode)
	delete(node.keys, key)
	if len(node.keys) == 0 {
		// Remove the freqNode.
		l.counterList.Remove(e)
	}
	delete(l.keyMap, key)
}
