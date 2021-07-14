package main

import (
	"errors"
	"log"
	"os"

	"github.com/seregproj/otus_hw/hw08_envdir_tool/exec"
	"github.com/seregproj/otus_hw/hw08_envdir_tool/reader"
)

var ErrInvalidArgumentsCount = errors.New("invalid arguments count")

func parseArgs() ([]string, error) {
	args := os.Args
	if len(args) < 3 {
		return nil, ErrInvalidArgumentsCount
	}

	return args, nil
}

func main() {
	args, err := parseArgs()
	if err != nil {
		log.Fatal(err)
	}

	dir := args[1]
	env, err := reader.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(exec.RunCmd(exec.Client{}, args[2:], env))
}
