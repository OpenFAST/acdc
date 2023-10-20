package mbc3_test

import (
	"acdc/mbc3"
	"path/filepath"
	"testing"
)

func TestAnalyze(t *testing.T) {

	// Find all model FAST files
	modelLinFiles, err := filepath.Glob("testdata/mbc/model.*.lin")
	if err != nil {
		t.Fatal(err)
	}

	linFileSets := [][]string{
		{"testdata/5MW_Land_BD_Linear.1.lin"},
		{"testdata/5MW_OC4Semi_Linear.1.lin"},
		{"testdata/Fake5MW_AeroLin_B1_UA4_DBEMT3.1.lin"},
		{"testdata/Fake5MW_AeroLin_B3_UA6.1.lin"},
		{"testdata/Ideal_Beam_Fixed_Free_Linear.1.lin"},
		{"testdata/Ideal_Beam_Free_Free_Linear.1.lin"},
		{"testdata/StC_test_OC4Semi_Linear_Nac.1.lin"},
		{"testdata/StC_test_OC4Semi_Linear_Tow.1.lin"},
		{"testdata/WP_Stationary_Linear.1.lin"},
		modelLinFiles,
	}

	for _, linFiles := range linFileSets {

		// Read linearization data from files
		linFileData := make([]*mbc3.LinData, len(linFiles))
		for i, f := range linFiles {
			if linFileData[i], err = mbc3.ReadLinFile(f); err != nil {
				t.Fatal(err)
			}
		}

		// Create matrix data from linearization file data
		matData := mbc3.NewMatData(linFileData)

		// Perform multi-blade coordinate transform
		mbc, err := matData.MBC3()
		if err != nil {
			t.Fatal(err)
		}

		// Perform Eigenanalysis
		_, err = mbc.EigenAnalysis()
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestMBC3(t *testing.T) {

	ld, err := mbc3.ReadLinFile("testdata/5MW_Land_BD_Linear.1.lin")
	// ld, err := mbc3.ReadLinFile("testdata/Fake5MW_AeroLin_B3_UA6.1.lin")
	if err != nil {
		t.Fatal(err)
	}

	// Create matrix data from linearization file data
	matData := mbc3.NewMatData([]*mbc3.LinData{ld})

	// Perform multi-blade coordinate transform
	mbc, err := matData.MBC3()
	if err != nil {
		t.Fatal(err)
	}

	mbc3.ToCSV(mbc.AvgA, "testdata/AvgA.csv", "%.7e")
}
