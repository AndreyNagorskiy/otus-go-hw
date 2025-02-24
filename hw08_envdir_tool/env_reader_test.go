package main

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReadDir(t *testing.T) {
	t.Run("dir not exists", func(t *testing.T) {
		env, err := ReadDir("testdata/test-test")
		require.ErrorIs(t, err, ErrDirectoryNotExists)
		require.Nil(t, env)
	})

	t.Run("path is not directory", func(t *testing.T) {
		env, err := ReadDir("testdata/echo.sh")
		require.ErrorIs(t, err, ErrPathIsNotDir)
		require.Nil(t, env)
	})

	t.Run("dir include file with equal in file name", func(t *testing.T) {
		env, err := ReadDir("testdata/equal_in_file_name")
		require.Error(t, err)
		require.Nil(t, env)
	})

	t.Run("success", func(t *testing.T) {
		env, err := ReadDir("testdata/env")
		require.NoError(t, err)

		// Ожидаемые значения в формате: ключ -> {значение, флаг NeedRemove}
		expected := map[string]struct {
			value      string
			needRemove bool
		}{
			"BAR":   {"bar", false},
			"EMPTY": {"", false},
			"FOO":   {"   foo\nwith new line", false},
			"HELLO": {"\"hello\"", false},
			"UNSET": {"", true},
		}

		for key, expectedValue := range expected {
			t.Run(fmt.Sprintf("checking %s", key), func(t *testing.T) {
				actual, exists := env[key]
				require.True(t, exists, "expected key %s not found in env", key)
				require.Equal(t, expectedValue.value, actual.Value, "unexpected value for key %s", key)
				require.Equal(t, expectedValue.needRemove, actual.NeedRemove, "unexpected NeedRemove for key %s", key)
			})
		}
	})
}
