package exec_test

import (
	"errors"
	"os"
	osexec "os/exec"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/seregproj/otus_hw/hw08_envdir_tool/exec"
	"github.com/seregproj/otus_hw/hw08_envdir_tool/reader"
	"github.com/stretchr/testify/require"
)

func TestRunCmdValid(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	command := []string{"ls", "-lah", "abc"}
	env := reader.Environment{"HELLO": reader.EnvValue{
		Value: "hello my friend!",
	}}
	envList := make([]string, len(os.Environ())+1)
	envList = append(envList, os.Environ()...)
	envList = append(envList, "HELLO=hello my friend!")

	mockClient := NewMockIExecClient(mockCtl)
	mockCmd := NewMockICmd(mockCtl)
	mockClient.EXPECT().Command("ls", "-lah", "abc").Return(mockCmd)
	mockCmd.EXPECT().SetEnv(envList)
	mockCmd.EXPECT().SetStdin(os.Stdin)
	mockCmd.EXPECT().SetStdout(os.Stdout)
	mockCmd.EXPECT().SetStderr(os.Stderr)
	mockCmd.EXPECT().Run().Return(nil)

	code := exec.RunCmd(mockClient, command, env)
	require.Equal(t, 0, code)
}

func TestRunCmdExitError(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	command := []string{"ls", "-lah", "abc"}
	env := reader.Environment{"HELLO": reader.EnvValue{
		Value: "hello my friend!",
	}}
	expErr := &osexec.ExitError{}
	envList := make([]string, len(os.Environ())+1)
	envList = append(envList, os.Environ()...)
	envList = append(envList, "HELLO=hello my friend!")

	mockClient := NewMockIExecClient(mockCtl)
	mockCmd := NewMockICmd(mockCtl)
	mockClient.EXPECT().Command("ls", "-lah", "abc").Return(mockCmd)
	mockCmd.EXPECT().SetEnv(envList)
	mockCmd.EXPECT().SetStdin(os.Stdin)
	mockCmd.EXPECT().SetStdout(os.Stdout)
	mockCmd.EXPECT().SetStderr(os.Stderr)
	mockCmd.EXPECT().Run().Return(expErr)

	code := exec.RunCmd(mockClient, command, env)
	require.Equal(t, -1, code)
}

func TestRunCmdUnexpectedError(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	command := []string{"ls", "-lah", "abc"}
	env := reader.Environment{"HELLO": reader.EnvValue{
		Value: "hello my friend!",
	}}
	expErr := errors.New("unexpected error")
	envList := make([]string, len(os.Environ())+1)
	envList = append(envList, os.Environ()...)
	envList = append(envList, "HELLO=hello my friend!")

	mockClient := NewMockIExecClient(mockCtl)
	mockCmd := NewMockICmd(mockCtl)
	mockClient.EXPECT().Command("ls", "-lah", "abc").Return(mockCmd)
	mockCmd.EXPECT().SetEnv(envList)
	mockCmd.EXPECT().SetStdin(os.Stdin)
	mockCmd.EXPECT().SetStdout(os.Stdout)
	mockCmd.EXPECT().SetStderr(os.Stderr)
	mockCmd.EXPECT().Run().Return(expErr)

	code := exec.RunCmd(mockClient, command, env)
	require.Equal(t, 1, code)
}
