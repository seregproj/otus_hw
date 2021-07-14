package exec

//go:generate mockgen -destination=mock_test.go -package=exec_test . IExecClient,ICmd

import (
	"fmt"
	"github.com/seregproj/otus_hw/hw08_envdir_tool/reader"
	"io"
	"os"
	"os/exec"
	"strings"
)

type ICmd interface {
	Run() error
	SetStdout(writer io.Writer)
	SetStdin(reader io.Reader)
	SetStderr(writer io.Writer)
	SetEnv(env []string)
}

type Cmd struct {
	command *exec.Cmd
}

func (c Cmd) SetStdin(reader io.Reader) {
	c.command.Stdin = reader
}

func (c Cmd) SetStdout(writer io.Writer) {
	c.command.Stdout = writer
}

func (c Cmd) SetStderr(writer io.Writer) {
	c.command.Stderr = writer
}

func (c Cmd) SetEnv(env []string) {
	c.command.Env = env
}

func (c Cmd) Run() error {
	return c.command.Run()
}

type IExecClient interface {
	Command(name string, args ...string) ICmd
}

type Client struct{}

func (ec Client) Command(name string, args ...string) ICmd {
	return Cmd{command: exec.Command(name, args...)}
}

func prepareEnv(env reader.Environment) []string {
	systemEnv := os.Environ()

	envList := make([]string, len(systemEnv)+len(env))

	for _, v := range systemEnv {
		str := strings.Split(v, "=")

		if envVal, ok := env[str[0]]; ok {
			if !envVal.NeedRemove {
				envList = append(envList, fmt.Sprintf("%v=%v", str[0], envVal.Value))
			}

			delete(env, str[0])
		} else {
			envList = append(envList, fmt.Sprintf("%v=%v", str[0], str[1]))
		}
	}

	for k, v := range env {
		if !v.NeedRemove {
			envList = append(envList, fmt.Sprintf("%v=%v", k, v.Value))
		}
	}

	return envList
}

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(client IExecClient, cmd []string, env reader.Environment) (returnCode int) {
	command := client.Command(cmd[0], cmd[1:]...)
	command.SetEnv(prepareEnv(env))
	command.SetStdin(os.Stdin)
	command.SetStdout(os.Stdout)
	command.SetStderr(os.Stderr)

	if err := command.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			return exiterr.ExitCode()
		} else {
			return 1
		}
	}

	return 0
}
