# `go-huffc`: Go Bindings for the Huff Compiler

[![Go Reference](https://pkg.go.dev/badge/github.com/project-blanc/go-huffc.svg)](https://pkg.go.dev/github.com/project-blanc/go-huffc)
[![Go Report Card](https://goreportcard.com/badge/github.com/project-blanc/go-huffc)](https://goreportcard.com/report/github.com/project-blanc/go-huffc)

`go-huffc` provides an easy way to compile [Huff](https://github.com/huff-language/huff-rs) contracts from Go.

> **Note**
> `go-huffc` requires the `huffc` binary to be installed. See [huff.sh](https://huff.sh) for installation instructions.

```
go get github.com/project-blanc/go-huffc
```

## Getting Started

```go
// Compile a contract with default compiler settings
c := huffc.New(nil)
contract, err := c.Compile("contract.huff")

// Compile a contract with custom compiler settings
c := huffc.New(&huff.Options{
    EVMVersion: huffc.EVMVersionIstanbul,
})
contract, err := c.Compile("contract.huff")
```
