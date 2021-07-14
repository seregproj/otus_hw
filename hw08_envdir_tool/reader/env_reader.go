package reader

import (
	"bufio"
	"errors"
	"os"
	"path"
	"strings"
	"unicode"
)

var ErrInvalidDirPath = errors.New("invalid dir path")

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	fi, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}

	if !fi.IsDir() {
		return nil, ErrInvalidDirPath
	}

	de, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment, len(de))
	for _, fe := range de {
		if !fe.Type().IsRegular() {
			continue
		}

		fi, err := fe.Info()
		if err != nil {
			return nil, err
		}

		if fi.Size() == 0 {
			env[fi.Name()] = EnvValue{NeedRemove: true}

			continue
		}

		if strings.ContainsRune(fi.Name(), '=') {
			continue
		}

		f, err := os.Open(path.Join(dir, fi.Name()))
		if err != nil {
			return nil, err
		}

		scanner := bufio.NewScanner(f)
		scanner.Scan()
		text := scanner.Text()
		lastIndex := strings.LastIndexFunc(text, func(r rune) bool {
			return !unicode.IsSpace(r)
		})

		text = strings.ReplaceAll(text, "\u0000", "\n")

		if lastIndex == -1 {
			env[fi.Name()] = EnvValue{NeedRemove: true}
		} else {
			env[fi.Name()] = EnvValue{Value: text[0 : lastIndex+1]}
		}
	}

	return env, nil
}
