# xscript

## Overview

`xscript` provides small helpers for launching scripts and inline code with a
specific runtime. Each subpackage wraps `github.com/neostd/go/exec` and exposes
a common shape for file-based and inline execution.

## Available Providers

- `bash`
- `bun`
- `deno`
- `dotnet`
- `elixir`
- `fish`
- `golang`
- `groovy`
- `lua`
- `node`
- `nushell`
- `php`
- `powershell`
- `pwsh`
- `python`
- `ruby`
- `sh`
- `zsh`

## Common API

Most providers expose the same helpers:

- `New` and `NewContext` to invoke the runtime directly
- `File` and `FileContext` to execute a script file
- `Inline` and `InlineContext` to execute inline code
- `Script` and `ScriptContext` to choose between file and inline modes based on the input

## Executable Resolution

Each provider resolves its runtime with `exec.Find(NAME, nil)` before creating the
command.

That lookup is driven by the executable registry entries each provider registers at
package init time. A provider typically supplies:

- the executable name, such as `node`, `php`, `bash`, or `nu`
- a provider-specific environment variable override
- a list of platform-specific candidate paths

In practice, that means `xscript` can find a runtime from:

1. a provider-specific environment variable
2. a matching executable already available on `PATH`
3. one of the registered fallback paths for the current platform

If no registered path is found, the providers still fall back to the bare runtime
name and let normal process execution resolve it.

## Overriding The Executable

You can control where a provider looks in four ways.

### 1. Set The Runtime On `PATH`

If the executable is already on `PATH`, the provider can resolve it there.

Examples:

- `node`
- `php`
- `python`
- `nu`
- `bash`

### 2. Set The Provider-Specific Environment Variable

Each provider registers a very specific override variable using `xscript.GetVar(...)`.
By default, that means variables such as:

- `BASH_EXE`
- `BUN_EXE`
- `DENO_EXE`
- `DOTNET_EXE`
- `ELIXIR_EXE`
- `FISH_EXE`
- `GO_EXE`
- `GROOVY_EXE`
- `LUA_EXE`
- `NODE_EXE`
- `NU_EXE`
- `PHP_EXE`
- `POWERSHELL_EXE`
- `PWSH_EXE`
- `PYTHON_EXE`
- `RUBY_EXE`
- `SH_EXE`
- `ZSH_EXE`

Example:

```powershell
$env:PHP_EXE = 'D:\tools\php\php.exe'
```

```bash
export NU_EXE="$HOME/.local/bin/nu"
```

### 3. Update The Registered Lookup Paths Programmatically

You can also change the registered executable metadata directly through
`exec.Registry`.

Important: `exec.Registry.Get(...)` returns a copy of the registered
`Executable`. If you modify it, write it back with `exec.Registry.Set(...)`.

```go
package main

import (
    neosexec "github.com/neostd/go/exec"
    _ "github.com/neostd/go/xscript/bash"
)

func main() {
    exe, ok := neosexec.Registry.Get("bash")
    if !ok {
        panic("bash is not registered")
    }

    exe.Linux = append([]string{"/opt/bash/bin/bash"}, exe.Linux...)
    exe.Darwin = append([]string{"/opt/homebrew/bin/bash"}, exe.Darwin...)
    exe.Windows = append([]string{"${SystemRoot}\\..\\custom\\bash.exe"}, exe.Windows...)

    neosexec.Registry.Set("bash", exe)
}
```

You can use the same pattern to:

- add or reorder `Windows`, `Linux`, or `Darwin` candidate paths
- set a fixed `Path`
- change the `Variable` name used for environment-based overrides

Example of changing the override variable name for a provider:

```go
package main

import (
    neosexec "github.com/neostd/go/exec"
    _ "github.com/neostd/go/xscript/php"
)

func main() {
    exe, ok := neosexec.Registry.Get("php")
    if !ok {
        panic("php is not registered")
    }

    exe.Variable = "MYAPP_PHP_BIN"
    neosexec.Registry.Set("php", exe)
}
```

### 4. Change The Variable Naming Convention

`xscript` exposes helpers for changing how those environment variable names are built:

- `xscript.SetPrefix("MYAPP")`
- `xscript.SetWindows("WIN")`

On non-Windows platforms, `GetVar("PHP")` normally produces `PHP_EXE`.
With `xscript.SetPrefix("MYAPP")`, it becomes `MYAPP_PHP_EXE`.

On Windows, if you also call `xscript.SetWindows("WIN")`, `GetVar("PHP")` becomes
`MYAPP_PHP_WIN_EXE`.

This is useful when you want application-specific overrides without colliding with
other environment variables.

## Registered Fallback Paths

Each provider also registers a small list of OS-specific fallback paths. These are
intended to catch common install locations when the runtime is not already on `PATH`
and no override variable is set.

Examples include locations such as:

- `/usr/bin/node`
- `/usr/local/bin/php`
- `${ProgramFiles}\PowerShell\7\pwsh.exe`
- `${SystemRoot}\..\msys64\usr\bin\zsh.exe`

If you need different fallback behavior, update the provider's `posix.go` and `win.go`
registration files.

If you need to change those paths at runtime instead of in source, use
`exec.Registry.Get(...)` and `exec.Registry.Set(...)` as shown above.

## Usage

```go
package main

import (
    "github.com/neostd/go/xscript/node"
)

func main() {
    cmd := node.Script("console.log('hello from node')")
    _, err := cmd.Run()
    if err != nil {
        panic(err)
    }
}
```

```go
package main

import (
    "github.com/neostd/go/xscript/bash"
)

func main() {
    cmd := bash.File("./script.sh", "arg1")
    _, err := cmd.Run()
    if err != nil {
        panic(err)
    }
}
```

## Running Commands And Capturing Output

Each provider returns `*exec.Cmd` values from `github.com/neostd/go/exec`.
That means you can use helpers such as `cmd.Run()` and `cmd.Output()` directly.

### `cmd.Run()`

`cmd.Run()` runs the command with inherited stdin, stdout, and stderr. It is
useful when you want the script output to stream directly to the terminal.

It still returns an `*exec.Result`, but `Stdout` and `Stderr` are not captured.
The returned result is mainly useful for metadata such as the executable path,
arguments, timestamps, temp file, and exit code.

```go
package main

import (
    "github.com/neostd/go/xscript/bash"
)

func main() {
    cmd := bash.Inline("echo hello from bash")

    result, err := cmd.Run()
    if err != nil && result == nil {
        panic(err)
    }

    if result != nil && result.Code != 0 {
        panic(result.ToError())
    }
}
```

### `cmd.Output()`

`cmd.Output()` runs the command and captures stdout and stderr into the returned
`*exec.Result`.

```go
package main

import (
    "fmt"

    "github.com/neostd/go/xscript/php"
)

func main() {
    cmd := php.Inline("fwrite(STDERR, 'warn'); echo 'hello';")

    result, err := cmd.Output()
    if err != nil && result == nil {
        panic(err)
    }

    if err != nil {
        fmt.Println("command error:", err)
    }
    fmt.Println("exit code:", result.Code)
    fmt.Println("stdout:", result.Text())
    fmt.Println("stderr:", result.ErrorText())
}
```

## Working With `exec.Result`

The value returned from `cmd.Run()` and `cmd.Output()` is `*exec.Result`.

Useful fields:

- `Code` for the process exit code
- `Stdout` for raw standard output bytes
- `Stderr` for raw standard error bytes
- `FileName` for the resolved executable path
- `Args` for the executed arguments
- `StartedAt` and `EndedAt` for timing information
- `TempFile` when a provider created a temporary source file for inline execution

Useful helper methods:

- `Text()` returns stdout as a string
- `ErrorText()` returns stderr as a string
- `Lines()` splits stdout into lines
- `Json()` unmarshals stdout as JSON
- `ErrorJson()` unmarshals stderr as JSON
- `IsOk()` reports whether the exit code was `0`
- `ToError()` converts a failed result into an error
- `Validate()` checks the default success condition
- `ValidateWith(...)` lets you define custom validation rules

Note: a non-zero exit status may also be returned as the `error` value from
`cmd.Run()` or `cmd.Output()`, so callers that want detailed diagnostics should
inspect both `result` and `err`.

### Exit Code And Failure Handling

```go
package main

import (
    "fmt"

    "github.com/neostd/go/xscript/python"
)

func main() {
    cmd := python.Inline("import sys; print('nope'); sys.exit(3)")

    result, err := cmd.Output()
    if err != nil && result == nil {
        panic(err)
    }

    if err != nil {
        fmt.Println("command error:", err)
    }

    if result != nil && !result.IsOk() {
        fmt.Println("exit code:", result.Code)
        fmt.Println("stdout:", result.Text())
        fmt.Println("stderr:", result.ErrorText())
        fmt.Println(result.ToError())
    }
}
```

### Working With Stdout And Stderr Separately

```go
package main

import (
    "fmt"

    "github.com/neostd/go/xscript/node"
)

func main() {
    cmd := node.Inline("console.log('out'); console.error('err')")

    result, err := cmd.Output()
    if err != nil && result == nil {
        panic(err)
    }

    fmt.Printf("stdout bytes: %q\n", result.Stdout)
    fmt.Printf("stderr bytes: %q\n", result.Stderr)
    fmt.Printf("stdout text: %s\n", result.Text())
    fmt.Printf("stderr text: %s\n", result.ErrorText())
}
```

### Parsing JSON Output

```go
package main

import (
    "fmt"

    "github.com/neostd/go/xscript/php"
)

func main() {
    cmd := php.Inline("echo json_encode(['ok' => true, 'value' => 42]);")

    result, err := cmd.Output()
    if err != nil && result == nil {
        panic(err)
    }

    data, err := result.Json()
    if err != nil {
        panic(err)
    }

    fmt.Printf("json: %#v\n", data)
}
```

### Custom Validation

```go
package main

import (
    "fmt"
    "strings"

    neosexec "github.com/neostd/go/exec"
    "github.com/neostd/go/xscript/bash"
)

func main() {
    cmd := bash.Inline("echo deployment-complete")

    result, err := cmd.Output()
    if err != nil && result == nil {
        panic(err)
    }

    ok, err := result.ValidateWith(func(r *neosexec.Result) (bool, error) {
        if r.Code != 0 {
            return false, r.ToError()
        }
        if !strings.Contains(r.Text(), "deployment-complete") {
            return false, fmt.Errorf("missing success marker")
        }
        return true, nil
    })
    if err != nil || !ok {
        panic(err)
    }
}
```

## Testing

Tests skip automatically when the target executable is not on `PATH`.

If `mise` is installed and `USE_MISE_FOR_TESTS=true`, test packages will try to:

1. install the runtime with `mise`
2. resolve its executable path with `mise which`
3. prepend the resolved bin directory to `PATH`

This keeps local and CI runs lightweight while still allowing broader coverage
when `mise` is available.

## License

This module is licensed under the MIT License. See `LICENSE.md`.
