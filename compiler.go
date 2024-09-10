package huffc

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var (
	ErrCompilationFailed = errors.New("compilation failed")

	codePrefix       = []byte("runtime: ")
	deployCodePrefix = []byte("bytecode: ")
)

// Compiler represents a compiler for Huff contracts.
type Compiler struct {
	cmd string
}

// New returns a new instance of the compiler.
func New() *Compiler {
	c := &Compiler{
		cmd: "huffc",
	}
	return c
}

// Compile the given huff-file with the given options and return the compiled
// contract.
func (c *Compiler) Compile(filename string, opts *Options) (*Contract, error) {
	// check if file exists
	stat, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		return nil, fmt.Errorf("file %q is a directory", filename)
	}

	// check options
	if opts == nil {
		opts = new(Options)
	}
	opts.setDefaults()

	ex := exec.Command(c.cmd,
		"--bin-runtime",
		"--bytecode",
		"--evm-version", string(opts.EVMVersion),
		filename,
	)
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	ex.Stdout = outBuf
	ex.Stderr = errBuf

	if err := ex.Run(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrCompilationFailed, strings.TrimSpace(errBuf.String()))
	}

	contract := new(Contract)
	s := bufio.NewScanner(outBuf)
	for s.Scan() {
		if code, ok := bytes.CutPrefix(s.Bytes(), codePrefix); ok {
			contract.Runtime = make([]byte, hex.DecodedLen(len(code)))
			if _, err := hex.Decode(contract.Runtime, code); err != nil {
				return nil, err
			}
		} else if code, ok := bytes.CutPrefix(s.Bytes(), deployCodePrefix); ok {
			contract.Constructor = make([]byte, hex.DecodedLen(len(code)))
			if _, err := hex.Decode(contract.Constructor, code); err != nil {
				return nil, err
			}
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	if contract.Runtime == nil || contract.Constructor == nil {
		return nil, fmt.Errorf("%w: unexpected error", ErrCompilationFailed)
	}
	return contract, nil
}

// Contract represents a compiled contract.
type Contract struct {
	Runtime     []byte // The runtime bytecode of the contract.
	Constructor []byte // The constructor bytecode of the contract.
}
