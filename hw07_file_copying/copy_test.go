package main

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("Unsupported file - /dev/random", func(t *testing.T) {
		err := Copy("/dev/random", "", 0, 0)
		require.ErrorIs(t, err, ErrUnsupportedFile)
	})

	t.Run("Unsupported file - any dir", func(t *testing.T) {
		err := Copy("./testdata", "", 0, 0)
		require.ErrorIs(t, err, ErrUnsupportedFile)
	})

	t.Run("Unexisting file", func(t *testing.T) {
		filename := fmt.Sprintf("/testdata/%s", strconv.FormatInt(time.Now().Unix(), 10))
		err := Copy(filename, "", 0, 0)
		require.Error(t, err)
	})

	t.Run("Empty file", func(t *testing.T) {
		tmpFilename, err := os.CreateTemp("/tmp", "test.")
		require.NoError(t, err)
		defer os.Remove(tmpFilename.Name())

		filename := "./testdata/empty.txt"
		err = Copy(filename, tmpFilename.Name(), 0, 0)
		require.NoError(t, err)
	})

	t.Run("Offset more than file size", func(t *testing.T) {
		fileName := "./testdata/input.txt"
		fi, _ := os.Stat(fileName)
		err := Copy(fileName, "", fi.Size()+1, 0)

		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("Offset equal file size", func(t *testing.T) {
		fileName := "./testdata/input.txt"
		fi, _ := os.Stat(fileName)
		err := Copy(fileName, "", fi.Size(), 0)

		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})
}
