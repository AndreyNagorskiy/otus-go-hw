package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	testEnv := Environment{
		"FOO":        {Value: "bar", NeedRemove: false},
		"EMPTY":      {Value: "", NeedRemove: false},
		"DELME":      {Value: "", NeedRemove: true},
		"WITH_SPACE": {Value: "some value", NeedRemove: false},
		"QUOTED":     {Value: "\"Hello World\"", NeedRemove: false},
	}

	t.Run("empty command error", func(t *testing.T) {
		var cmd []string
		code := RunCmd(cmd, testEnv)

		require.Equal(t, code, ErrorCode)
	})

	t.Run("wrong command error", func(t *testing.T) {
		cmd := []string{"ls", "-la", "sdgsgsdfdsfs"}
		code := RunCmd(cmd, testEnv)

		require.Equal(t, code, ErrorCode)
	})

	t.Run("success", func(t *testing.T) {
		cmd := []string{"echo", "Hello World"}
		code := RunCmd(cmd, testEnv)

		require.Equal(t, code, SuccessCode)
	})
}
