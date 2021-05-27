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

type list struct {
	len       int
	firstItem *ListItem
	lastItem  *ListItem
}

func (l list) Len() int {
	return l.len
}

func (l list) Front() *ListItem {
	return l.firstItem
}

func (l list) Back() *ListItem {
	return l.lastItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	newListItem := ListItem{Value: v}

	if l.firstItem == nil {
		l.firstItem = &newListItem
		l.lastItem = &newListItem
	} else {
		l.firstItem.Prev = &newListItem
		newListItem.Next = l.firstItem
		l.firstItem = &newListItem
	}

	l.len++
	return &newListItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newListItem := ListItem{Value: v}

	if l.lastItem == nil {
		l.lastItem = &newListItem
		l.firstItem = &newListItem
	} else {
		l.lastItem.Next = &newListItem
		newListItem.Prev = l.lastItem
		l.lastItem = &newListItem
	}

	l.len++
	return &newListItem
}

func (l *list) Remove(i *ListItem) {
	switch {
	case i.Prev == nil && i.Next == nil:
		l.firstItem = nil
		l.lastItem = nil
	case i.Prev == nil:
		if i.Next.Next == nil {
			l.firstItem = l.lastItem
		}

		i.Next.Prev = nil
	case i.Next == nil:
		if i.Prev.Prev == nil {
			l.lastItem = l.firstItem
		}

		i.Prev.Next = nil
	default:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i.Prev == nil {
		return
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.lastItem = i.Prev
	}

	i.Prev.Next = i.Next
	l.firstItem.Prev = i
	i.Next = l.firstItem
	l.firstItem = i
	l.firstItem.Prev = nil
}

func NewList() List {
	return new(list)
}
