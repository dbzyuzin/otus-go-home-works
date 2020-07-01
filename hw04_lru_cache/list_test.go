package hw04_lru_cache //nolint:golint,stylecheck

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, l.Len(), 0)
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Back().Next // 20
		l.Remove(middle)        // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Back(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{50, 30, 10, 40, 60, 80, 70}, elems)
	})

	t.Run("nil на границах", func(t *testing.T) {
		l := NewList()

		require.Nil(t, l.Front())
		require.Nil(t, l.Back())

		l.PushFront(1)

		require.NotNil(t, l.Front())
		require.Equal(t, l.Back(), l.Front())

		l.PushFront(2)
		l.Remove(l.Front())
		l.Remove(l.Front())

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Front())
	})

	t.Run("Проверка последовательности", func(t *testing.T) {
		l := NewList()

		seq := []int{1, 2, 3, 4}
		for _, elem := range seq {
			l.PushFront(elem)
		}

		elems := make([]int, 0, l.Len())
		for i := l.Back(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}

		for i := range seq {
			require.Equal(t, seq[i], elems[i])
		}
	})

	t.Run("Для покрытие проверка на пуш в конец в пустой лист", func(t *testing.T) {
		l := NewList()

		l.PushBack(2)

		require.Equal(t, 1, l.Len())
		require.Equal(t, l.Front(), l.Back())
	})
}
