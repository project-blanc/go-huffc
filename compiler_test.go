package huffc_test

import (
	"errors"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/project-blanc/go-huffc"
)

func TestCompilerCompile(t *testing.T) {
	tests := []struct {
		Filename     string
		Options      *huffc.Options
		WantContract *huffc.Contract
		WantErr      error
	}{
		{
			Filename: "stop.huff",
			WantContract: &huffc.Contract{
				Runtime:     []byte{0x00},
				Constructor: []byte{0x60, 0x01, 0x80, 0x60, 0x09, 0x3d, 0x39, 0x3d, 0xf3, 0x00},
			},
		},
		{
			Filename: "push.huff",
			WantContract: &huffc.Contract{
				Runtime:     []byte{0x5f},
				Constructor: []byte{0x60, 0x01, 0x80, 0x60, 0x09, 0x3d, 0x39, 0x3d, 0xf3, 0x5f},
			},
		},
		{
			Filename: "push.huff",
			Options:  &huffc.Options{EVMVersion: huffc.EVMVersionParis},
			WantContract: &huffc.Contract{
				Runtime:     []byte{0x60, 0x00},
				Constructor: []byte{0x60, 0x02, 0x80, 0x60, 0x09, 0x3d, 0x39, 0x3d, 0xf3, 0x60, 0x00},
			},
		},
		{
			Filename: "error.huff",
			WantErr:  huffc.ErrCompilationFailed,
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			c := huffc.New()
			gotContract, gotErr := c.Compile(filepath.Join("testdata", test.Filename), test.Options)
			if !errors.Is(gotErr, test.WantErr) {
				t.Errorf("Err: want %v, got %v", test.WantErr, gotErr)
			}
			if diff := cmp.Diff(test.WantContract, gotContract); diff != "" {
				t.Errorf("Contract (-want +got):\n%s", diff)
			}
		})
	}
}
