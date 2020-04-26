package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	Len() int                          // длина списка
	Front() *listItem                  // первый Item
	Back() *listItem                   // последний Item
	PushFront(v interface{}) *listItem // добавить значение в начало
	PushBack(v interface{}) *listItem  // добавить значение в конец
	Remove(i *listItem)                // удалить элемент
	MoveToFront(i *listItem)           // переместить элемент в начало
}

type listItem struct {
	Value interface{} // значение
	Next  *listItem   // следующий элемент
	Prev  *listItem   // предыдущий элемент

}

type list struct {
	front *listItem
	back  *listItem
	len   int
}

func (l list) Len() int {
	return l.len
}

func (l list) Front() *listItem {
	return l.front
}

func (l list) Back() *listItem {
	return l.back
}

func (l list) PushFront(v interface{}) *listItem {
	item := &listItem{
		Value: v,
		Next:  l.front,
	}
	if item.Next != nil {
		item.Next.Prev = item
	}

	l.front = item
	l.len++
	return l.front
}

func (l list) PushBack(v interface{}) *listItem {
	item := &listItem{
		Value: v,
		Prev:  l.back,
	}
	if item.Prev != nil {
		item.Prev.Next = item
	}

	l.back = item
	l.len++
	return l.back
}

func (l list) Remove(i *listItem) {
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	l.len--
}

func (l list) MoveToFront(i *listItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return &list{}
}
