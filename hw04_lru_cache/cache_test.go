package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("set one item", func(t *testing.T) {
		// new item
		c := NewCache(2)
		expKey := Key("el")
		expVal := 1
		exists := c.Set(expKey, expVal)
		require.False(t, exists)
		val, ok := c.Get(expKey)
		require.Equal(t, expVal, val)
		require.True(t, ok)

		// update existing item by the same val
		exists = c.Set(expKey, expVal)
		require.True(t, exists)
		val, ok = c.Get(expKey)
		require.Equal(t, expVal, val)
		require.True(t, ok)

		// update existing item by other val
		expVal2 := 2
		exists = c.Set(expKey, expVal2)
		require.True(t, exists)
		val, ok = c.Get(expKey)
		require.Equal(t, expVal2, val)
		require.True(t, ok)
	})

	t.Run("set some items", func(t *testing.T) {
		c := NewCache(3)
		expKey1 := Key("el1")
		expVal1 := 1
		expKey2 := Key("el2")
		expVal2 := "a"

		exists1 := c.Set(expKey1, expVal1)
		exists2 := c.Set(expKey2, expVal2)
		require.False(t, exists1)
		require.False(t, exists2)
		val, ok := c.Get(expKey1)
		require.Equal(t, expVal1, val)
		require.True(t, ok)
		val, ok = c.Get(expKey2)
		require.Equal(t, expVal2, val)
		require.True(t, ok)

		exists1 = c.Set(expKey1, expVal2)
		exists2 = c.Set(expKey2, expVal1)
		require.True(t, exists1)
		require.True(t, exists2)
		val, ok = c.Get(expKey1)
		require.Equal(t, expVal2, val)
		require.True(t, ok)
		val, ok = c.Get(expKey2)
		require.Equal(t, expVal1, val)
		require.True(t, ok)
	})

	t.Run("substitution", func(t *testing.T) {
		c := NewCache(3)
		key1 := Key("el1")
		val1 := 1
		key2 := Key("el2")
		val2 := 2
		key3 := Key("el3")
		val3 := "iii"

		c.Set(key1, val1) // [1]
		c.Set(key2, val2) // [2, 1]
		c.Set(key3, val3) // ["iii", 2, 1]

		c.Get(key1) // [1, "iii", 2]
		c.Get(key2) // [2, 1, "iii"]
		c.Get(key3) // ["iii", 2, 1]
		c.Get(key1) // [1, "iii", 2]
		c.Get(key2) // [2, 1, "iii"]

		// adding new item with substitution lru
		key4 := Key("el_insert")
		val4 := "new el"
		c.Set(key4, val4) // ["new el", 2, 1]

		// not exists
		val, ok := c.Get(key3)
		require.False(t, ok)
		require.Nil(t, val)

		// exists
		val, ok = c.Get(key1)
		require.True(t, ok)
		require.Equal(t, val1, val)

		val, ok = c.Get(key2)
		require.True(t, ok)
		require.Equal(t, val2, val)

		val, ok = c.Get(key4)
		require.True(t, ok)
		require.Equal(t, val4, val)

		// move to front of queue
		newVal2 := "🙂"
		c.Set(key2, newVal2) // [🙂, "new el", 1]

		// repeat adding new item with substitution lru
		key5 := Key("el_repeat_insert")
		val5 := "repeat new el"
		c.Set(key5, val5) // ["repeat new el", 🙂, "new el"]

		// not exists
		val, ok = c.Get(key1)
		require.False(t, ok)
		require.Nil(t, val)

		// exists
		val, ok = c.Get(key2)
		require.True(t, ok)
		require.Equal(t, newVal2, val)

		val, ok = c.Get(key4)
		require.True(t, ok)
		require.Equal(t, val4, val)

		val, ok = c.Get(key5)
		require.True(t, ok)
		require.Equal(t, val5, val)
	})

	t.Run("Clear items", func(t *testing.T) {
		c := NewCache(2)
		c.Set("key1", 100)
		c.Set("key2", 200)
		c.Clear()

		val, ok := c.Get("key1")
		require.False(t, ok)
		require.Nil(t, val)
		val, ok = c.Get("key2")
		require.False(t, ok)
		require.Nil(t, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()

		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()

		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	go func() {
		defer wg.Done()

		for i := 0; i < 1_000_000; i++ {
			c.Clear()
		}
	}()

	wg.Wait()
}
