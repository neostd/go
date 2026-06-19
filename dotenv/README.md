# dotenv

## Overview

The dotenv package enables parsing, reading, editing, and writing dotenv (`.env`) files in Go.

`Parse` returns an `EnvDoc` representing the parsed content of a dotenv file.
`EnvDoc` provides methods to access and manipulate environment variables while
retaining comments and original ordering.

Preserving order matters when writing the document back out and when later doing
your own expansion or evaluation passes.

If you only need key/value pairs, `ToMap` converts an `EnvDoc` to `map[string]string`.

## Public API

- `Parse(string) (*EnvDoc, error)` parses dotenv content
- `Lex(string) ([]*Token, error)` tokenizes dotenv content
- `NewDoc()` and `NewDocument()` create empty documents
- `EnvDoc` provides add/get/set/merge/string helpers for `.env` content

## Usage

To use `dotenv`, import the module in your Go project:

```go
import "github.com/neostd/go/dotenv"

func main() {
    content := `KEY1="value1"

KEY2='value2'
Key3=value3
Key4=a value with spaces
# This is a comment
Key5="a value with \"escaped quotes\""
Key6='a value with \'single quotes\''
Key7="line1
line2
line3
"
Key8="value with \nnewlines"
Key9="value with \t tabs"
`

    doc, err := dotenv.Parse(content)
    if err != nil {
        panic(err)
    }

    println("Parsed keys:")
    for _, key := range doc.Keys() {
        value, _ := doc.Get(key)
        println(key, "=", value)
    }

    output := doc.String()
    println("Output content:")
    println(output)

    doc2 := dotenv.NewDocument()
    doc2.AddNewline()
    doc2.AddComment("This is a new comment")
    doc2.AddVariable("KEY1", "value1")
    doc2.AddQuotedVariable("KEY2", "value2", '"')

    println("New document content:")
    println(doc2.String())
}
```

## Working With Documents

`EnvDoc` preserves comments and ordering, so it works well when you want to
modify an existing `.env` file and write it back out without flattening it into
just a map.

```go
doc := dotenv.NewDoc()
doc.AddComment("Application settings")
doc.AddVariable("APP_ENV", "development")
doc.AddQuotedVariable("APP_NAME", "my app", '"')

doc.Set("APP_ENV", "production")

println(doc.String())
```
