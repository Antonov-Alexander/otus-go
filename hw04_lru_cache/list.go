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
	length int
	front  *ListItem
	back   *ListItem
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	return l.pushFrontListItem(&ListItem{Value: v})
}

func (l *list) PushBack(v interface{}) *ListItem {
	return l.pushBackListItem(&ListItem{Value: v})
}

func (l *list) Remove(i *ListItem) {
	l.length--

	if i == l.back {
		l.back = i.Prev
		if l.back != nil {
			l.back.Next = nil
		}

		return
	}

	if i == l.front {
		l.front = i.Next
		if l.front != nil {
			l.front.Prev = nil
		}

		return
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
}

func (l *list) MoveToFront(i *ListItem) {
	if i != l.front {
		l.Remove(i)
		l.pushFrontListItem(i)
	}
}

func (l *list) pushFrontListItem(i *ListItem) *ListItem {
	i.Prev = nil
	i.Next = l.front

	if l.front != nil {
		l.front.Prev = i
	}

	if l.length == 0 {
		l.back = i
	}

	l.front = i
	l.length++
	return i
}

func (l *list) pushBackListItem(i *ListItem) *ListItem {
	i.Next = nil
	i.Prev = l.back

	if l.back != nil {
		l.back.Next = i
	}

	if l.length == 0 {
		l.front = i
	}

	l.back = i
	l.length++
	return i
}

func NewList() List {
	return new(list)
}
