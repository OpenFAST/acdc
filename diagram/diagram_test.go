package diagram_test

import (
	"acdc/diagram"
	"acdc/lin"
	"path/filepath"
	"regexp"
	"testing"
)

func TestDiagramNew(t *testing.T) {

	// Search for linearization files
	LinFiles, err := filepath.Glob("../../samples/NREL_5MW-ED/case01/*.lin")
	if err != nil {
		t.Fatal(err)
	}
	linRe := regexp.MustCompile(`.+?\.\d+\.lin`)
	tmp := LinFiles
	LinFiles = []string{}
	for _, f := range tmp {
		if linRe.MatchString(f) {
			LinFiles = append(LinFiles, f)
		}
	}

	// Process linearization files into results
	linOPs, err := lin.ProcessFiles(LinFiles)
	if err != nil {
		t.Fatal(err)
	}

	// Generate diagram
	diag, err := diagram.New(linOPs, diagram.Options{
		MinFreq: 0.0,
		MaxFreq: 3.5,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%#v", diag)
}
