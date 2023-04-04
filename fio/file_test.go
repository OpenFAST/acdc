package fio

import (
	"bytes"
	"os"
	"testing"
)

func TestParseFormat(t *testing.T) {

	for _, name := range []string{
		"AeroDyn14",
		"AeroDyn15",
		"AeroDynBlade",
		"AirfoilInfo",
		"BeamDyn",
		"BeamDynBlade",
		"ElastoDyn",
		"ElastoDynBlade",
		"ElastoDynTower",
		// "HydroDyn",
		"InflowWind",
		"Main",
		"ServoDyn",
		"SubDyn",
	} {

		s, err := parse(name, "testdata/"+name+".dat")
		if err != nil {
			t.Fatalf("%s: %s", name, err)
		}

		w := &bytes.Buffer{}
		if err := s.Format(w); err != nil {
			t.Fatalf("%s: %s", name, err)
		}

		err = os.WriteFile("testdata/"+name+"_test.dat", w.Bytes(), 0777)
		if err != nil {
			t.Fatalf("%s: %s", name, err)
		}

		_, err = parse(name, "testdata/"+name+"_test.dat")
		if err != nil {
			t.Fatalf("%s: %s", name, err)
		}
	}
}

func TestParseAll(t *testing.T) {

	// mainPath := "testdata/openfast/5MW_Land_BD_DLL_WTurb/5MW_Land_BD_DLL_WTurb.fst"
	// mainPath := "testdata/openfast/5MW_Land_BD_Init/5MW_Land_BD_Init.fst"
	// mainPath := "testdata/openfast/EllipticalWing_OLAF/EllipticalWing_OLAF.fst"
	mainPath := "testdata/openfast/AOC_WSt/AOC_WSt.fst"
	files, err := ParseAll(mainPath)
	if err != nil {
		t.Fatalf("%s: %s", mainPath, err)
	}

	err = FormatAll(files, "testdata/formatAll/base")
	if err != nil {
		t.Fatalf("%s: %s", mainPath, err)
	}
}
