package test

import (
	"github.com/jifanchn/go-libiec61850/iec61850"
	"testing"
)

func TestIEC61850LoadICD(t *testing.T) {
	scl, err := iec61850.GetSCL("test_icd.icd")
	if err != nil {
		t.Error(err)
	}
	scl.Print()
}
