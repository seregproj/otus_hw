package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListBase(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("PushFront", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		require.Equal(t, 1, l.Len())
		require.Equal(t, 10, l.Front().Value)
		require.Equal(t, 10, l.Back().Value)

		l.PushFront(30) // [30, 10]
		require.Equal(t, 2, l.Len())
		require.Equal(t, 30, l.Front().Value)
		require.Equal(t, 10, l.Front().Next.Value)
		require.Equal(t, 10, l.Back().Value)
		require.Equal(t, 30, l.Back().Prev.Value)

		l.PushFront(50) // [50, 30, 10]
		require.Equal(t, 3, l.Len())
		require.Equal(t, 50, l.Front().Value)
		require.Equal(t, 30, l.Front().Next.Value)
		require.Equal(t, 50, l.Front().Next.Prev.Value)
		require.Equal(t, 10, l.Front().Next.Next.Value)
		require.Equal(t, 10, l.Back().Value)
		require.Equal(t, 30, l.Back().Prev.Value)
	})

	t.Run("PushBack", func(t *testing.T) {
		l := NewList()

		l.PushBack(10) // [10]
		require.Equal(t, 1, l.Len())
		require.Equal(t, 10, l.Back().Value)
		require.Equal(t, 10, l.Front().Value)

		l.PushBack(30) // [10, 30]
		require.Equal(t, 2, l.Len())
		require.Equal(t, 30, l.Back().Value)
		require.Equal(t, 10, l.Back().Prev.Value)
		require.Equal(t, 10, l.Front().Value)
		require.Equal(t, 30, l.Front().Next.Value)

		l.PushBack(50) // [10, 30, 50]
		require.Equal(t, 3, l.Len())
		require.Equal(t, 50, l.Back().Value)
		require.Equal(t, 30, l.Back().Prev.Value)
		require.Equal(t, 50, l.Back().Prev.Next.Value)
		require.Equal(t, 10, l.Front().Value)
		require.Equal(t, 30, l.Front().Next.Value)
		require.Equal(t, 10, l.Back().Prev.Prev.Value)
	})

	t.Run("Remove", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.Remove(l.Front())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
		require.Equal(t, 0, l.Len())

		l.PushBack(10) // [10]
		l.Remove(l.Back())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
		require.Equal(t, 0, l.Len())

		l.PushFront(10)
		l.PushBack("c") // [10, "c"]
		l.Remove(l.Front())
		require.Equal(t, 1, l.Len())
		require.Equal(t, "c", l.Front().Value)
		require.Equal(t, "c", l.Back().Value)
		require.Nil(t, l.Back().Prev)
		require.Nil(t, l.Front().Next)

		l.PushFront(10) // [10, "c"]
		l.Remove(l.Back())
		require.Equal(t, 1, l.Len())
		require.Equal(t, 10, l.Back().Value)
		require.Equal(t, 10, l.Front().Value)
		require.Nil(t, l.Front().Next)
		require.Nil(t, l.Back().Prev)

		l.PushBack("c")
		l.PushFront("a") // ["a", 10, "c"]

		l.Remove(l.Front().Next) // ["a", "c"]
		require.Equal(t, "a", l.Front().Value)
		require.Equal(t, "c", l.Back().Value)
		require.Equal(t, "c", l.Front().Next.Value)
		require.Equal(t, "a", l.Back().Prev.Value)
		require.Equal(t, 2, l.Len())
	})

	t.Run("MoveToFront", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.MoveToFront(l.Front())
		require.Equal(t, 1, l.Len())
		require.Equal(t, 10, l.Front().Value)
		require.Equal(t, 10, l.Back().Value)

		l.PushBack(20) // [10, 20]

		l.MoveToFront(l.Back()) // [20, 10]
		require.Equal(t, 2, l.Len())
		require.Equal(t, 20, l.Front().Value)
		require.Equal(t, 10, l.Back().Value)
		require.Equal(t, 10, l.Front().Next.Value)
		require.Equal(t, 20, l.Back().Prev.Value)

		l.PushFront("abc") // ["abc", 20, 10]
		l.PushFront("new") // ["new", "abc", 20, 10]

		l.MoveToFront(l.Front().Next) // ["abc", "new", 20, 10]
		require.Equal(t, 4, l.Len())
		require.Equal(t, "abc", l.Front().Value)
		require.Equal(t, "new", l.Front().Next.Value)
		require.Nil(t, l.Front().Prev)
		require.Equal(t, "abc", l.Front().Next.Prev.Value)
		require.Equal(t, 20, l.Front().Next.Next.Value)
		require.Equal(t, 10, l.Back().Value)
		require.Equal(t, 20, l.Back().Prev.Value)

		l.MoveToFront(l.Front().Next) // ["new", "abc", 20, 10]

		require.Equal(t, 4, l.Len())
		require.Equal(t, "new", l.Front().Value)
		require.Equal(t, "abc", l.Front().Next.Value)
		require.Nil(t, l.Front().Prev)
		require.Equal(t, "new", l.Front().Next.Prev.Value)
		require.Equal(t, 20, l.Front().Next.Next.Value)
		require.Equal(t, 10, l.Back().Value)
		require.Equal(t, 20, l.Back().Prev.Value)
	})
}

func TestListComplex(t *testing.T) {
	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
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
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}
