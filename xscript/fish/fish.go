package fish

import (
	"context"
	"strings"

	"github.com/neostd/go/exec"
)

const NAME = "fish"

var ScriptArgs = []string{}

func New(args ...string) *exec.Cmd {
	exe, _ := exec.Find(NAME, nil)
	if exe == "" {
		exe = NAME
	}

	return exec.New(exe, args...)
}

func NewContext(ctx context.Context, args ...string) *exec.Cmd {
	exe, _ := exec.Find(NAME, nil)
	if exe == "" {
		exe = NAME
	}

	return exec.NewContext(ctx, exe, args...)
}

func File(path string, args ...string) *exec.Cmd {
	allArgs := append(append([]string{}, ScriptArgs...), path)
	allArgs = append(allArgs, args...)
	return New(allArgs...)
}

func FileContext(ctx context.Context, path string, args ...string) *exec.Cmd {
	allArgs := append(append([]string{}, ScriptArgs...), path)
	allArgs = append(allArgs, args...)
	return NewContext(ctx, allArgs...)
}

func Inline(script string, args ...string) *exec.Cmd {
	allArgs := append(append([]string{}, ScriptArgs...), "-c", script)
	allArgs = append(allArgs, args...)
	return New(allArgs...)
}

func InlineContext(ctx context.Context, script string, args ...string) *exec.Cmd {
	allArgs := append(append([]string{}, ScriptArgs...), "-c", script)
	allArgs = append(allArgs, args...)
	return NewContext(ctx, allArgs...)
}

func Script(script string, args ...string) *exec.Cmd {
	if !strings.ContainsAny(script, "\n\r") {
		trimmed := strings.TrimSpace(script)
		if strings.HasSuffix(trimmed, ".fish") {
			return File(trimmed, args...)
		}
	}

	return Inline(script, args...)
}

func ScriptContext(ctx context.Context, script string, args ...string) *exec.Cmd {
	if !strings.ContainsAny(script, "\n\r") {
		trimmed := strings.TrimSpace(script)
		if strings.HasSuffix(trimmed, ".fish") {
			return FileContext(ctx, trimmed, args...)
		}
	}

	return InlineContext(ctx, script, args...)
}
