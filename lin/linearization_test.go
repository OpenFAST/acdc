package lin_test

import (
	"acdc/lin"
	"path/filepath"
	"testing"
)

func TestLoadResults(t *testing.T) {

	dir := "testdata"

	LinFiles, err := filepath.Glob(filepath.Join(dir, "*.lin"))
	if err != nil {
		t.Fatal(err)
	}
	if len(LinFiles) == 0 {
		t.Fatal("no lin files found")
	}

	_, err = lin.ProcessFiles(LinFiles)
	if err != nil {
		t.Fatal(err)
	}
}
