# ordered list

[![Github Actions](https://github.com/Allenxuxu/orderedlist/workflows/CI/badge.svg)](https://github.com/Allenxuxu/orderedlist/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/Allenxuxu/orderedlist)](https://goreportcard.com/report/github.com/Allenxuxu/orderedlist)

支持并发操作的有序链表

让每个链表元素持有一个锁，从而缩小锁粒度，来支持并发操作。

```go
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
```
