# exec

## Overview

The `exec` package wraps `os/exec` and adds convenience helpers for building,
running, capturing, piping, and locating commands.

## Public API

- `New`, `NewContext`, `Command`, and `CommandContext` for building commands
- `Run` and `Output` for one-shot string command execution
- `Cmd` chainable helpers for args, env, cwd, stdio, logging, and execution
- `Pipeline` helpers for piping commands together
- `Which`, `WhichFirst`, `Find`, and the executable registry for resolution helpers
- `Result` helpers for stdout/stderr access and validation

## Usage

To use `exec`, import the module in your Go project:

```go
import "github.com/neostd/go/exec"

func main() {
    cmd := exec.New("ls", "-l")
    result, err := cmd.Run()
    if err != nil {
        panic(err)
    }
    println("Command output:", string(result.Stdout))

    o, err := exec.Run("echo test")
    if err != nil {
        panic(err)
    }
    if !o.IsOk() {
        err := o.ToError()
        panic(err)
    }

    o2, err := exec.Output("echo hello")
    if o2.IsOk() {
        for _, line := range o2.Lines() {
            println(line)
        }
    }

    o3, err := exec.Command("echo 'Hello World'").PipeCommand("grep Hello").Output()
    if err != nil {
        panic(err)
    }
    println("Piped command output:", string(o3.Stdout))
}

```
