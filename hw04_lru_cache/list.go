package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

func NewListItem(value interface{}, next *ListItem, prev *ListItem) *ListItem {
	return &ListItem{
		Value: value,
		Next:  next,
		Prev:  prev,
	}
}

type list struct {
	head *ListItem
	tail *ListItem
	len  int
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	listItem := NewListItem(v, l.head, nil)

	if l.head != nil {
		l.head.Prev = listItem
	} else {
		l.tail = listItem
	}

	l.head = listItem
	l.len++

	return l.head
}

func (l *list) PushBack(v interface{}) *ListItem {
	listItem := NewListItem(v, nil, l.tail)

	if l.tail != nil {
		l.tail.Next = listItem
	} else {
		l.head = listItem
	}

	l.tail = listItem
	l.len++

	return l.head
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.head = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.tail = i.Prev
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil || i == l.head {
		return
	}

	l.Remove(i)
	l.PushFront(i.Value)
}
