package orderedlist

import (
	"sync"
	at "sync/atomic"

	"github.com/Allenxuxu/toolkit/sync/atomic"
)

type List interface {
	// Contains 检查一个元素是否存在，如果存在则返回 true，否则返回 false
	Contains(value int) bool
	// Insert 插入一个元素，如果此操作成功插入一个元素，则返回 true，否则返回 false
	Insert(value int) bool
	// Delete 删除一个元素，如果此操作成功删除一个元素，则返回 true，否则返回 false
	Delete(value int) bool
	// Range 遍历此有序链表的所有元素，如果 f 返回 false，则停止遍历
	Range(f func(value int) bool)
	// Len 返回有序链表的元素个数
	Len() int
}

var _ List = (*ConcurrentList)(nil)

type ConcurrentList struct {
	head *node
	len  atomic.Int64
}

func New() *ConcurrentList {
	return &ConcurrentList{head: &node{value: 0}}
}

type node struct {
	sync.Mutex

	marked atomic.Bool
	value  int
	next   at.Value
}

func (c *ConcurrentList) Contains(value int) bool {
	a := c.head
	var b *node
	if v := a.next.Load(); v != nil {
		b = v.(*node)
	}

	for b != nil && b.value < value {
		a = b
		b = b.next.Load().(*node)
	}
	// Check if b is not exists
	if b == nil || b.value != value {
		return false
	}

	return !b.marked.Get()
}

func (c *ConcurrentList) Insert(value int) bool {
AGAIN:
	a := c.head
	var b *node
	if v := a.next.Load(); v != nil {
		b = v.(*node)
	}
	for b != nil && b.value < value {
		a = b
		b = b.next.Load().(*node)
	}

	// Check if the node is exist.
	if b != nil && b.value == value {
		return false
	}

	x := &node{value: value}

	a.Lock()
	if (a.next.Load() != nil && a.next.Load().(*node) != b) || a.marked.Get() {
		a.Unlock()
		goto AGAIN
	}
	defer a.Unlock()

	x.next.Store(b)
	a.next.Store(x)
	c.len.Add(1)
	return true
}

func (c *ConcurrentList) Delete(value int) bool {
AGAIN:
	a := c.head
	var b *node
	if v := a.next.Load(); v != nil {
		b = v.(*node)
	}
	for b != nil && b.value < value {
		a = b
		b = b.next.Load().(*node)
	}
	// Check if b is not exists
	if b == nil || b.value != value {
		return false
	}

	b.Lock()
	if b.marked.Get() {
		b.Unlock()
		goto AGAIN
	}

	a.Lock()
	if a.next.Load().(*node) != b || a.marked.Get() {
		a.Unlock()
		b.Unlock()
		goto AGAIN
	}
	defer b.Unlock()
	defer a.Unlock()

	b.marked.Set(true)
	a.next.Store(b.next.Load().(*node))

	c.len.Add(-1)
	return true
}

func (c *ConcurrentList) Range(f func(value int) bool) {
	a := c.head
	b := a.next.Load().(*node)
	for b != nil {
		if !f(b.value) {
			break
		}

		a = b
		b = b.next.Load().(*node)
	}
}

func (c *ConcurrentList) Len() int {
	return int(c.len.Get())
}
