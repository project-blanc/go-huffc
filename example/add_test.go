package examples_test

import (
	"errors"
	"math/big"
	"testing"

	"github.com/lmittmann/w3"
	"github.com/lmittmann/w3/w3types"
	"github.com/lmittmann/w3/w3vm"
	"github.com/project-blanc/go-huffc"
)

var (
	compiler = huffc.New()

	// address of the contract
	contractAddr = w3vm.RandA()

	// ABI binding for the add function
	funcAdd = w3.MustNewFunc("add(uint256,uint256)", "uint256")

	maxBig256 = new(big.Int).Sub(
		new(big.Int).Exp(w3.Big2, big.NewInt(256), nil),
		w3.Big1,
	)
)

func TestAdd(t *testing.T) {
	tests := []struct {
		A, B    *big.Int
		WantC   *big.Int // C = A + B
		WantErr error
	}{
		{A: big.NewInt(0), B: big.NewInt(0), WantC: big.NewInt(0)},
		{A: big.NewInt(0), B: big.NewInt(1), WantC: big.NewInt(1)},
		{A: big.NewInt(1), B: big.NewInt(0), WantC: big.NewInt(1)},
		{A: big.NewInt(1), B: big.NewInt(2), WantC: big.NewInt(3)},
		{A: maxBig256, B: big.NewInt(0), WantC: maxBig256},
		{A: big.NewInt(0), B: maxBig256, WantC: maxBig256},
		{A: maxBig256, B: maxBig256, WantErr: w3vm.ErrRevert},
	}

	contract, err := compiler.Compile("add.huff", &huffc.Options{EVMVersion: huffc.EVMVersionShanghai})
	if err != nil {
		t.Fatalf("Failed to compile contract: %v", err)
	}

	for _, test := range tests {
		// construct a VM for testing with the runtime code of the contract
		// deployed to `contractAddr`
		vm, _ := w3vm.New(
			w3vm.WithState(w3types.State{
				contractAddr: {Code: contract.Runtime},
			}),
		)

		// apply the add function
		receipt, err := vm.Apply(&w3types.Message{
			From: w3vm.RandA(),
			To:   &contractAddr,
			Args: []any{test.A, test.B},
			Func: funcAdd,
		})

		// check for error
		if !errors.Is(err, test.WantErr) {
			t.Fatalf("Err: want %q, got %q", test.WantErr, err)
		} else if err != nil {
			continue
		}

		// check the return value
		var gotC *big.Int
		if err := receipt.DecodeReturns(&gotC); err != nil {
			t.Fatalf("Failed to decode return value: %v", err)
		}
		if test.WantC.Cmp(gotC) != 0 {
			t.Fatalf("C: want %v, got %v", test.WantC, gotC)
		}
	}
}

func FuzzAdd(f *testing.F) {
	contract, err := compiler.Compile("add.huff", &huffc.Options{EVMVersion: huffc.EVMVersionShanghai})
	if err != nil {
		f.Fatalf("Failed to compile contract: %v", err)
	}

	preState := w3vm.WithState(w3types.State{
		contractAddr: {Code: contract.Runtime},
	})

	f.Fuzz(func(t *testing.T, aBytes, bBytes []byte) {
		if len(aBytes) > 32 || len(bBytes) > 32 {
			t.Skip()
		}
		var (
			a     = new(big.Int).SetBytes(aBytes[:])
			b     = new(big.Int).SetBytes(bBytes[:])
			wantC = new(big.Int).Add(a, b)
		)

		vm, _ := w3vm.New(preState)

		receipt, err := vm.Apply(&w3types.Message{
			From: w3vm.RandA(),
			To:   &contractAddr,
			Args: []any{a, b},
			Func: funcAdd,
		})

		if wantC.Cmp(maxBig256) <= 0 {
			// gotC does NOT overflow, we expect correct addition
			if err != nil {
				t.Fatalf("Err: want nil, got %q", err)
			}

			var gotC *big.Int
			if err := receipt.DecodeReturns(&gotC); err != nil {
				t.Fatalf("Failed to decode return value: %v", err)
			}
			if wantC.Cmp(gotC) != 0 {
				t.Fatalf("C: want %v+%v=%v, got %v", a, b, wantC, gotC)
			}
		} else {
			// gotC overflows, we expect a revert
			if wantErr := w3vm.ErrRevert; !errors.Is(err, wantErr) {
				t.Fatalf("Err: want %q, got %q", wantErr, err)
			}
		}
	})
}
