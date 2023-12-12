# `go-huffc`: Go Bindings for the Huff Compiler

[![Go Reference](https://pkg.go.dev/badge/github.com/project-blanc/go-huffc.svg)](https://pkg.go.dev/github.com/project-blanc/go-huffc)
[![Go Report Card](https://goreportcard.com/badge/github.com/project-blanc/go-huffc)](https://goreportcard.com/report/github.com/project-blanc/go-huffc)

`go-huffc` provides an easy way to compile [Huff](https://github.com/huff-language/huff-rs) contracts from Go.

> [!NOTE]
> `go-huffc` requires the `huffc` binary to be installed. See [huff.sh](https://huff.sh) for installation instructions.

```
go get github.com/project-blanc/go-huffc
```

## Getting Started

```go
// Compile a contract with default compiler settings
c := huffc.New()
contract, err := c.Compile("contract.huff", nil)

// Compile a contract with custom compiler settings
c := huffc.New()
contract, err := c.Compile("contract.huff", &huffc.Options{
    EVMVersion: huffc.EVMVersionIstanbul,
})
```

## Example Project

See the [example project](https://github.com/project-blanc/go-huffc/tree/main/example) for a basic reference on how to test and fuzz a Huff contract in Go.


> [!WARNING]
> This package is pre-1.0. There might be breaking changes between minor versions.
