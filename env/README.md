# env

## Overview

The env package provides utilities for working with environment variables in Go.

It includes functions for expanding environment variables, handling default values,
manipulating the path variable, and parsing bash-style variable interpolations.

## Public API

- `Get`, `Set`, `Unset`, `Has`, and `All` for basic environment access
- `GetPath`, `SetPath`, `SplitPath`, `JoinPath`, `PrependPath`, `AppendPath`, and `HasPath` for path helpers
- `Expand` and `ExpandWithOptions` for shell-style variable expansion
- `ExpandOptions` and `With...` helpers for configuring expansion behavior

## Usage

To use `env`, import the module in your Go project:

```go
import "github.com/neostd/go/env"

func main() {
    // Example of expanding an environment variable
    value, err := env.Expand("$HOME")
    if err != nil {
        panic(err)
    }
    fmt.Println("Expanded value:", value)

    if err := env.Set("FOO", "bar"); err != nil {
        panic(err)
    }

    // Example of using default values
    valueWithDefault, err := env.Expand("${FOO:-default}")
    if err != nil {
        panic(err)
    }
    fmt.Println("Value with default:", valueWithDefault)


    // command substitution
    output, err := env.Expand(
        "Value: $(echo hi)",
        env.WithCommandSubstitution(true),
        env.WithEnableShellExpansion(true),
    )
    if err != nil {
        panic(err)
    }
    fmt.Println("Command substitution output:", output)
}

```

## Expansion Options

`Expand` supports bash-style interpolation patterns such as:

- `$HOME`
- `${NAME}`
- `${NAME:-default}`
- `${NAME:=default}`
- `${NAME:?message}`
- `$(command)` when command substitution is enabled

You can also provide custom getters and setters:

```go
value, err := env.Expand(
    "${APP_ENV:-development}",
    env.WithGet(func(key string) string {
        if key == "APP_ENV" {
            return "production"
        }
        return ""
    }),
)
if err != nil {
    panic(err)
}

println(value)
```
