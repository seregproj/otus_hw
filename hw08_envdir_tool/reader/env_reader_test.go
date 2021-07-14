package reader

import (
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("No such dir", func(t *testing.T) {
		env, err := ReadDir("invalid_path")
		require.Error(t, err)
		require.Nil(t, env)
	})

	t.Run("Invalid dir", func(t *testing.T) {
		env, err := ReadDir(path.Join("testdata", "not_dir"))
		require.ErrorIs(t, err, ErrInvalidDirPath)
		require.Nil(t, env)
	})

	t.Run("Valid case", func(t *testing.T) {
		expected := Environment{"BAR": EnvValue{
			Value:      "bar",
			NeedRemove: false,
		}, "EMPTY": EnvValue{
			Value:      "",
			NeedRemove: true,
		}, "FOO": EnvValue{
			Value:      "   foo\nwith new line",
			NeedRemove: false,
		}, "HELLO": EnvValue{
			Value:      `"hello"`,
			NeedRemove: false,
		}, "UNSET": EnvValue{
			Value:      "",
			NeedRemove: true,
		}, "ПРИВЕТ": EnvValue{
			Value:      "привет!",
			NeedRemove: false,
		}}

		env, err := ReadDir(path.Join("testdata", "env"))
		require.NoError(t, err)
		require.Equal(t, expected, env)
	})
}
