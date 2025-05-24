package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopySystemFiles(t *testing.T) {
	tests := []struct {
		name     string
		fromPath string
		wantErr  error
	}{
		{
			name:     "Copy from /dev/null",
			fromPath: "/dev/null",
			wantErr:  ErrUnsupportedFile,
		},
		{
			name:     "Copy from /proc/cpuinfo",
			fromPath: "/proc/cpuinfo",
			wantErr:  ErrUnsupportedFile,
		},
		{
			name:     "Copy from /dev/urandom",
			fromPath: "/dev/urandom",
			wantErr:  ErrUnsupportedFile,
		},
		{
			name:     "Copy from /run/systemd",
			fromPath: "/run/systemd",
			wantErr:  ErrUnsupportedFile,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Copy(tt.fromPath, "out.txt", 0, 0)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestCopy(t *testing.T) {
	t.Run("from path equal to path", func(t *testing.T) {
		err := Copy("testdata/input.txt", "testdata/input.txt", 0, 0)
		require.ErrorIs(t, err, ErrFromPathEqualToPath)
	})

	t.Run("offset exceeds file size", func(t *testing.T) {
		err := Copy("testdata/input.txt", "out.txt", 1000000, 0)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})

	tests := []struct {
		name   string
		offset int64
		limit  int64
	}{
		{
			name:   "Copy without offset and limit",
			offset: 0,
			limit:  0,
		},
		{
			name:   "Copy with offset",
			offset: 10,
			limit:  0,
		},
		{
			name:   "Copy with limit",
			offset: 0,
			limit:  10,
		},
		{
			name:   "Copy with limit and offset",
			offset: 10,
			limit:  10,
		},
	}

	t.Run("without errors", func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := Copy("testdata/input.txt", "out.txt", 0, 0)
				require.NoError(t, err)
				err = os.Remove("out.txt")
				if err != nil {
					t.Fatal(err)
				}
			})
		}
	})
}
