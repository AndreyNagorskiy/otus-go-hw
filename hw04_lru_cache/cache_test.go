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

		c = NewCache(0)
		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		val, exist := c.Get("aaa")
		require.Nil(t, val)
		require.False(t, exist)

		c = NewCache(-2)
		wasInCache = c.Set("aaa", 100)
		require.False(t, wasInCache)

		val, exist = c.Get("aaa")
		require.Nil(t, val)
		require.False(t, exist)
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
		c := NewCache(10)

		var keyName Key = "purge"

		c.Set(keyName, "test-value")

		val, ok := c.Get(keyName)
		require.Equal(t, "test-value", val)
		require.True(t, ok)

		c.Clear()

		val, ok = c.Get(keyName)
		require.Nil(t, val)
		require.False(t, ok)
	})

	t.Run("push-out logic", func(t *testing.T) {
		c := NewCache(2)
		c.Set("a", 1)
		c.Set("b", 2)
		c.Set("c", 3)

		val, ok := c.Get("a")
		require.Nil(t, val)
		require.False(t, ok)

		c = NewCache(3)
		c.Set("a", 1)
		c.Set("b", 2)
		c.Set("c", 3)

		c.Get("b")
		c.Get("c")
		c.Get("a")
		c.Set("c", "test")
		c.Set("d", 4)

		val, ok = c.Get("b")
		require.Nil(t, val)
		require.False(t, ok)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

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

	wg.Wait()
}
