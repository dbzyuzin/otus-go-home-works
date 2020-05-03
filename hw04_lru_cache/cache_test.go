package hw04_lru_cache //nolint:golint,stylecheck

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

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		cache := NewCache(2)

		for _, i := range []int{1, 2, 3, 4} {
			cache.Set(strconv.Itoa(i), i)
		} // ["4", "3"]

		_, ok := cache.Get("3") // ["3", "4"]
		require.True(t, ok)

		_, ok = cache.Get("2")
		require.False(t, ok)

		_, ok = cache.Get("1")
		require.False(t, ok)

		cache.Set("newkey", 1000) // ["newkey", "3"]
		_, ok = cache.Get("4")
		require.False(t, ok)
		_, ok = cache.Get("3")
		require.True(t, ok)
		_, ok = cache.Get("newkey")
		require.True(t, ok)
	})

	t.Run("clear", func(t *testing.T) {
		cache := NewCache(2)

		for _, i := range []int{1, 2} {
			cache.Set(strconv.Itoa(i), i)
		}

		cache.Clear()
		_, ok := cache.Get("1")
		require.False(t, ok)
		_, ok = cache.Get("2")
		require.False(t, ok)
	})
}

func TestCacheMultithreading(t *testing.T) {

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(strconv.Itoa(i), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(strconv.Itoa(rand.Intn(1_000_000)))
		}
	}()

	wg.Wait()
}
