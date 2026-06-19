package exec_test

import (
	"runtime"
	"strings"
	"testing"

	"github.com/neostd/go/exec"
	"github.com/stretchr/testify/assert"
)

func testEchoCommand(text string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.New("cmd", "/C", "echo", text)
	}

	return exec.New("echo", text)
}

func testEchoString(text string) string {
	if runtime.GOOS == "windows" {
		return "cmd /C echo " + text
	}

	return "echo '" + text + "'"
}

func testPipeCommand() (*exec.Pipeline, bool) {
	if runtime.GOOS == "windows" {
		_, hasCmd := exec.Which("cmd")
		_, hasFindstr := exec.Which("findstr")
		if !hasCmd || !hasFindstr {
			return nil, false
		}
		return exec.New("cmd", "/C", "echo", "Hello World").PipeCommand("findstr Hello"), true
	}

	_, hasGrep := exec.Which("grep")
	_, hasEcho := exec.Which("echo")
	if !hasEcho || !hasGrep {
		return nil, false
	}
	return exec.Command("echo 'Hello World'").PipeCommand("grep Hello"), true
}

func TestNewCommandOutput(t *testing.T) {
	check := "echo"
	if runtime.GOOS == "windows" {
		check = "cmd"
	}
	_, ok := exec.Which(check)
	if !ok {
		t.Skip(check + " not found")
	}

	o, err := testEchoCommand("hello").Output()
	assert.NoError(t, err)
	assert.Equal(t, 0, o.Code)
	assert.Equal(t, "hello", strings.TrimSpace(o.Text()))
}

func TestCommandOutput(t *testing.T) {
	check := "echo"
	if runtime.GOOS == "windows" {
		check = "cmd"
	}
	_, ok := exec.Which(check)
	if !ok {
		t.Skip(check + " not found")
	}

	cmd := testEchoString("hello world")

	o, err := exec.Command(cmd).Output()
	assert.NoError(t, err)
	assert.Equal(t, 0, o.Code)
	assert.Equal(t, "hello world", strings.TrimSpace(o.Text()))
}

func TestPipeCommand(t *testing.T) {
	pipeline, ok := testPipeCommand()
	if !ok {
		t.Skip("required pipe commands not found")
	}

	o, err := pipeline.Output()
	assert.NoError(t, err)
	assert.Equal(t, 0, o.Code)
	assert.Equal(t, "Hello World", strings.Trim(strings.TrimSpace(o.Text()), `"'`))
}
