// Package exec provides higher-level helpers around os/exec.
package exec

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/neostd/go/cmdargs"
)

const (
	// STDIO_INHERIT inherits stdio from the current process.
	STDIO_INHERIT = 0
	// STDIO_PIPED uses an in-memory pipe or buffer for stdio.
	STDIO_PIPED = 1
	// STDIO_NULL discards stdio.
	STDIO_NULL = 2
)

var (
	logger func(cmd *Cmd)
)

// Cmd wraps exec.Cmd with additional helpers.
type Cmd struct {
	*exec.Cmd
	ctx           *context.Context // if true, the command is a context command
	logger        func(cmd *Cmd)
	disableLogger bool
	TempFile      *string
}

// New creates a new command from an executable name and args.
func New(name string, args ...string) *Cmd {
	cmd := exec.Command(name, args...)
	return &Cmd{Cmd: cmd}
}

// NewContext creates a new command bound to a context.
func NewContext(ctx context.Context, name string, args ...string) *Cmd {
	cmd := exec.CommandContext(ctx, name, args...)
	return &Cmd{Cmd: cmd, ctx: &ctx}
}

// SetLogger sets the package-wide command logger.
func SetLogger(f func(cmd *Cmd)) {
	logger = f
}

// SetLogger sets a command-specific logger.
func (c *Cmd) SetLogger(f func(cmd *Cmd)) {
	c.logger = f
}

// DisableLogger disables logging for this command.
func (c *Cmd) DisableLogger() {
	c.disableLogger = true
}

// CommandContext parses a command string and binds it to a context.
func CommandContext(ctx context.Context, command string) *Cmd {
	exe := ""
	args := cmdargs.Split(command).ToArray()
	if len(args) > 0 {
		exe = args[0]
		args = args[1:]
	}
	return NewContext(ctx, exe, args...)
}

// Command parses a command string and returns a new Cmd.
func Command(command string) *Cmd {
	exe := ""
	args := cmdargs.Split(command).ToArray()
	if len(args) > 0 {
		exe = args[0]
		args = args[1:]
	}

	return New(exe, args...)
}

// Run parses and runs a command string with inherited stdio.
func Run(command string) (*Result, error) {
	args := cmdargs.Split(command).ToArray()
	if len(args) == 0 {
		return nil, errors.New("command cannot be empty")
	}
	if len(args) == 1 {
		return New(args[0]).Run()
	}

	command = args[0]
	args = args[1:]
	c := New(command, args...)
	return c.Run()
}

// Output parses and runs a command string while capturing stdio.
func Output(command string) (*Result, error) {
	args := cmdargs.Split(command).ToArray()
	if len(args) == 0 {
		return nil, errors.New("command cannot be empty")
	}
	if len(args) == 1 {
		return New(args[0]).Output()
	}

	command = args[0]
	args = args[1:]
	c := New(command, args...)
	return c.Output()
}

// AppendArgs appends arguments to the command.
func (c *Cmd) AppendArgs(args ...string) *Cmd {
	c.Args = append(c.Args, args...)
	return c
}

// PrependArgs prepends arguments to the command.
func (c *Cmd) PrependArgs(args ...string) *Cmd {
	c.Args = append([]string{args[0]}, c.Args...)
	return c
}

// WithArgs replaces the command arguments.
func (c *Cmd) WithArgs(args ...string) *Cmd {
	c.Args = args
	return c
}

// AppendEnv appends environment entries to the command.
func (c *Cmd) AppendEnv(env ...string) *Cmd {
	c.Env = append(c.Env, env...)
	return c
}

// PrependEnv prepends environment entries to the command.
func (c *Cmd) PrependEnv(env ...string) *Cmd {
	c.Env = append([]string{env[0]}, c.Env...)
	return c
}

// WithEnvMap replaces the environment using a map.
func (c *Cmd) WithEnvMap(env map[string]string) *Cmd {
	data := make([]string, 0)
	for k, v := range env {
		data = append(data, k+"="+v)
	}
	return c.WithEnv(data...)
}

// WithEnv replaces the command environment.
func (c *Cmd) WithEnv(env ...string) *Cmd {
	c.Env = env
	return c
}

// WithCwd sets the working directory.
func (c *Cmd) WithCwd(dir string) *Cmd {
	c.Dir = dir
	return c
}

// WithStdin sets the command stdin reader.
func (c *Cmd) WithStdin(stdin io.Reader) *Cmd {
	c.Stdin = stdin
	return c
}

// WithStdout sets the command stdout writer.
func (c *Cmd) WithStdout(stdout io.Writer) *Cmd {
	c.Stdout = stdout
	return c
}

// WithStderr sets the command stderr writer.
func (c *Cmd) WithStderr(stderr io.Writer) *Cmd {
	c.Stderr = stderr
	return c
}

// WithStdio configures stdin, stdout, and stderr behavior.
func (c *Cmd) WithStdio(stdin, stdout, stderr int) *Cmd {
	switch stdin {
	case STDIO_INHERIT:
		c.Stdin = os.Stdin
	case STDIO_PIPED:
		c.Stdin = bytes.NewBuffer(nil)
	case STDIO_NULL:
		c.Stdin = nil
	}

	switch stdout {
	case STDIO_INHERIT:
		c.Stdout = os.Stdout
	case STDIO_PIPED:
		c.Stdout = bytes.NewBuffer(nil)
	case STDIO_NULL:
		c.Stdout = nil
	}

	switch stderr {
	case STDIO_INHERIT:
		c.Stderr = os.Stderr
	case STDIO_PIPED:
		c.Stderr = bytes.NewBuffer(nil)
	case STDIO_NULL:
		c.Stderr = nil
	}

	return c
}

// Quiet runs the command without inheriting or capturing output.
func (c *Cmd) Quiet() (*Result, error) {
	c.Stdout = nil
	c.Stderr = nil
	var out Result
	out.FileName = c.Path
	out.Args = c.Args
	out.Stdout = make([]byte, 0)
	out.Stderr = make([]byte, 0)
	// use utc time
	out.StartedAt = time.Now().UTC()

	err := c.Start()
	if err != nil {
		return nil, err
	}

	err = c.Wait()
	if err != nil {
		return nil, err
	}
	out.EndedAt = time.Now().UTC()
	out.Code = c.ProcessState.ExitCode()
	if c.TempFile != nil {
		out.TempFile = c.TempFile
	}

	return &out, nil
}

// Run runs the command with inherited stdio.
func (c *Cmd) Run() (*Result, error) {
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	var out Result
	out.FileName = c.Path
	out.Args = c.Args
	// use utc time
	out.StartedAt = time.Now().UTC()
	out.Stdout = make([]byte, 0)
	out.Stderr = make([]byte, 0)

	err := c.Start()
	if err != nil {
		out.EndedAt = time.Now().UTC()
		out.Code = 1
		return &out, err
	}

	err = c.Wait()
	if err != nil {
		out.EndedAt = time.Now().UTC()
		out.Code = 1
		return &out, err
	}

	out.EndedAt = time.Now().UTC()
	out.Code = c.ProcessState.ExitCode()
	if c.TempFile != nil {
		out.TempFile = c.TempFile
	}

	return &out, nil
}

// Output runs the command and captures stdout and stderr.
func (c *Cmd) Output() (*Result, error) {

	var out Result
	out.Stdout = make([]byte, 0)
	out.Stderr = make([]byte, 0)
	out.StartedAt = time.Now().UTC()
	out.FileName = c.Path
	out.Args = c.Args

	var outb, errb bytes.Buffer
	c.Stdout = &outb
	c.Stderr = &errb

	err := c.Start()
	if err != nil {
		out.EndedAt = time.Now().UTC()
		out.Code = 1
		return &out, err
	}

	err = c.Wait()
	if err != nil {
		out.EndedAt = time.Now().UTC()
		out.Code = 1
		return &out, err
	}

	out.EndedAt = time.Now().UTC()
	out.Code = c.ProcessState.ExitCode()
	out.Stdout = outb.Bytes()
	out.Stderr = errb.Bytes()
	if c.TempFile != nil {
		out.TempFile = c.TempFile
	}

	return &out, nil
}

// Start starts the command, invoking any configured loggers first.
func (c *Cmd) Start() error {
	if c.disableLogger {
		return c.Cmd.Start()
	}

	if c.logger != nil {
		c.logger(c)
	}

	if logger != nil {
		logger(c)
	}

	p := c.Path
	if p != "" && !filepath.IsAbs(p) {
		p2, err := Find(p, nil)
		if err == nil {
			c.Path = p2
		}
	}

	return c.Cmd.Start()
}

// Wait waits for the command to exit.
func (c *Cmd) Wait() error {
	return c.Cmd.Wait()
}
