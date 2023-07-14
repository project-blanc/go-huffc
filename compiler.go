package huff

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
	evmVersion EVMVersion
}

// Options for the [Compiler]. A zero Options consists entirely of default values.
type Options struct {
	EVMVersion EVMVersion
}

// New returns a new instance of the compiler that is configured with the given options.
func New(opts *Options) *Compiler {
	c := &Compiler{
		evmVersion: defaultEVMVersion,
	}
	if opts == nil {
		return c
	}

	if opts.EVMVersion != "" {
		c.evmVersion = opts.EVMVersion
	}
	return c
}

// Compile the given huff-file and return the compiled contract.
func (c *Compiler) Compile(filename string) (*Contract, error) {
	// check if file exists
	stat, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		return nil, fmt.Errorf("file %q is a directory", filename)
	}

	ex := exec.Command("huffc",
		"--bin-runtime",
		"--bytecode",
		"--evm-version", string(c.evmVersion),
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
			contract.Code = make([]byte, hex.DecodedLen(len(code)))
			if _, err := hex.Decode(contract.Code, code); err != nil {
				return nil, err
			}
		} else if code, ok := bytes.CutPrefix(s.Bytes(), deployCodePrefix); ok {
			contract.DeployCode = make([]byte, hex.DecodedLen(len(code)))
			if _, err := hex.Decode(contract.DeployCode, code); err != nil {
				return nil, err
			}
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	if contract.Code == nil || contract.DeployCode == nil {
		return nil, fmt.Errorf("%w: unexpected error", ErrCompilationFailed)
	}
	return contract, nil
}
