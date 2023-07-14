package huff_test

import (
	"errors"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/project-blanc/go-huff"
)

func TestCompilerCompile(t *testing.T) {
	tests := []struct {
		Filename     string
		Options      *huff.Options
		WantContract *huff.Contract
		WantErr      error
	}{
		{
			Filename: "stop.huff",
			WantContract: &huff.Contract{
				Code:       []byte{0x00},
				DeployCode: []byte{0x60, 0x01, 0x80, 0x60, 0x09, 0x3d, 0x39, 0x3d, 0xf3, 0x00},
			},
		},
		{
			Filename: "push.huff",
			WantContract: &huff.Contract{
				Code:       []byte{0x5f},
				DeployCode: []byte{0x60, 0x01, 0x80, 0x60, 0x09, 0x3d, 0x39, 0x3d, 0xf3, 0x5f},
			},
		},
		{
			Filename: "push.huff",
			Options:  &huff.Options{EVMVersion: huff.EVMVersionParis},
			WantContract: &huff.Contract{
				Code:       []byte{0x60, 0x00},
				DeployCode: []byte{0x60, 0x02, 0x80, 0x60, 0x09, 0x3d, 0x39, 0x3d, 0xf3, 0x60, 0x00},
			},
		},
		{
			Filename: "error.huff",
			WantErr:  huff.ErrCompilationFailed,
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			c := huff.New(test.Options)
			gotContract, gotErr := c.Compile(filepath.Join("testdata", test.Filename))
			if !errors.Is(gotErr, test.WantErr) {
				t.Errorf("Err: want %v, got %v", test.WantErr, gotErr)
			}
			if diff := cmp.Diff(test.WantContract, gotContract); diff != "" {
				t.Errorf("Contract (-want +got):\n%s", diff)
			}
		})
	}
}
