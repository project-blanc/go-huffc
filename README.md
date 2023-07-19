# `go-huff`: Go Bindings for the Huff Compiler

`go-huff` provides an easy way to compile [Huff](https://github.com/huff-language/huff-rs) contracts from Go.

> **Note**
> `go-huff` requires the `huffc` binary to be installed. See [huff.sh](https://huff.sh) for installation instructions.

```
go get github.com/project-blanc/go-huff
```

## Getting Started

```go
// Compile a contract with default compiler settings
c := huff.New(nil)
contract, err := c.Compile("contract.huff")

// Compile a contract with custom compiler settings
c := huff.New(&huff.Options{
    EVMVersion: huff.EVMVersionIstanbul,
})
contract, err := c.Compile("contract.huff")
```
