// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package list implements a doubly linked list.
//
// To iterate over a list (where l is a *List):
//
//	for e := l.Front(); e != nil; e = e.Next() {
//		// do something with e.Value
//	}
package list

// Element is an element of a linked list.
type Element struct {
	// Next and previous pointers in the doubly-linked list of elements.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next element of the last
	// list element (l.Back()) and the previous element of the first list
	// element (l.Front()).
	next, prev *Element

	// 此元素所属的列表。
	list *List

	// The value stored with this element.
	Value any
}

// Next returns the next list element or nil.
func (e *Element) Next() *Element {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Prev returns the previous list element or nil.
func (e *Element) Prev() *Element {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// List 表示双向链表。List 的零值是一个可供使用的空列表。.
type List struct {
	root Element // Sentinel 列表元素，仅使用 &root、root.prev 和 root.next
	len  int     // 当前列表长度，不包括 （this） Sentinel 元素
}

// Init initializes or clears list l.
func (l *List) Init() *List {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// New returns an initialized list.
func New() *List { return new(List).Init() }

// Len returns the number of elements of list l.
// The complexity is O(1).
func (l *List) Len() int { return l.len }

// Front returns the first element of list l or nil if the list is empty.
func (l *List) Front() *Element {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// Back returns the last element of list l or nil if the list is empty.
func (l *List) Back() *Element {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// lazyInit lazily initializes a zero List value.
func (l *List) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

// insert inserts e after at, increments l.len, and returns e.
func (l *List) insert(e, at *Element) *Element {
	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
	e.list = l
	l.len++
	return e
}

// insertValue is a convenience wrapper for insert(&Element{Value: v}, at).
func (l *List) insertValue(v any, at *Element) *Element {
	return l.insert(&Element{Value: v}, at)
}

// remove removes e from its list, decrements l.len
func (l *List) remove(e *Element) {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	e.list = nil
	l.len--
}

// move moves e to next to at.
func (l *List) move(e, at *Element) {
	if e == at {
		return
	}
	e.prev.next = e.next
	e.next.prev = e.prev

	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
}

// Remove 如果 e 是列表 l 的元素，则从 l 中删除 e。它返回元素值 e.Value。该元素不能为 nil。
func (l *List) Remove(e *Element) any {
	if e.list == l {
		// if e.list == l, l must have been initialized when e was inserted
		// in l or l == nil (e is a zero Element) and l.remove will crash
		l.remove(e)
	}
	return e.Value
}

// PushFront 在列表 L 的前面插入一个值为 v 的新元素 e，并返回 e。
func (l *List) PushFront(v any) *Element {
	l.lazyInit()
	return l.insertValue(v, &l.root)
}

// PushBack 在列表 L 的后面插入一个值为 v 的新元素 e，并返回 e。
func (l *List) PushBack(v any) *Element {
	l.lazyInit()
	return l.insertValue(v, l.root.prev)
}

// InsertBefore 在 mark 之前插入一个值为 v 的新元素 e，并返回 e。如果 mark 不是 l 的元素，则不会修改列表。标记不得为 nil
func (l *List) InsertBefore(v any, mark *Element) *Element {
	if mark.list != l {
		return nil
	}
	// see comment in List.Remove about initialization of l
	return l.insertValue(v, mark.prev)
}

// InsertAfter 在 mark 之后立即插入一个值为 v 的新元素 e，并返回 e。如果 mark 不是 l 的元素，则不会修改列表。标记不得为 nil。
func (l *List) InsertAfter(v any, mark *Element) *Element {
	if mark.list != l {
		return nil
	}
	// see comment in List.Remove about initialization of l
	return l.insertValue(v, mark)
}

// MoveToFront 将元素 E 移动到列表 L 的前面。如果 e 不是 l 的元素，则不会修改列表。该元素不能为 nil。
func (l *List) MoveToFront(e *Element) {
	if e.list != l || l.root.next == e {
		return
	}
	// see comment in List.Remove about initialization of l
	l.move(e, &l.root)
}

// MoveToBack 将元素 E 移动到列表 L 的后面。如果 e 不是 l 的元素，则不会修改列表。该元素不能为 nil。
func (l *List) MoveToBack(e *Element) {
	if e.list != l || l.root.prev == e {
		return
	}
	// see comment in List.Remove about initialization of l
	l.move(e, l.root.prev)
}

// MoveBefore 将元素 E 移动到 Mark 之前的新位置。如果 e 或 mark 不是 l 的元素，或者 e == mark，则不会修改列表。元素和标记不能为 nil。
func (l *List) MoveBefore(e, mark *Element) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.move(e, mark.prev)
}

// MoveAfter 将元素 E 移动到标记后的新位置。如果 e 或 mark 不是 l 的元素，或者 e == mark，则不会修改列表。元素和标记不能为 nil。
func (l *List) MoveAfter(e, mark *Element) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.move(e, mark)
}

// PushBackList 在列表 L 的后面插入另一个列表的副本。列表 l 和其他可能相同。它们不能为零。
func (l *List) PushBackList(other *List) {
	l.lazyInit()
	for i, e := other.Len(), other.Front(); i > 0; i, e = i-1, e.Next() {
		l.insertValue(e.Value, l.root.prev)
	}
}

// PushFrontList 在列表 L 的前面插入另一个列表的副本。列表 l 和其他可能相同。它们不能为零。
func (l *List) PushFrontList(other *List) {
	l.lazyInit()
	for i, e := other.Len(), other.Back(); i > 0; i, e = i-1, e.Prev() {
		l.insertValue(e.Value, &l.root)
	}
}
