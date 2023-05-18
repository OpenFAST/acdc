package fio

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseAll(t *testing.T) {

	for _, path := range []string{
		"5MW_Land_BD_DLL_WTurb/5MW_Land_BD_DLL_WTurb.fst",
		"5MW_Land_BD_Init/5MW_Land_BD_Init.fst",
		"EllipticalWing_OLAF/EllipticalWing_OLAF.fst",
		"AOC_WSt/AOC_WSt.fst",
	} {
		main := &Main{}
		err := main.Parse("testdata/reference/" + path)
		if err != nil {
			t.Fatal(err)
		}
		outDir := filepath.Join("testdata", filepath.Dir(path))
		if err := os.MkdirAll(outDir, 0777); err != nil {
			t.Fatal(err)
		}
		err = main.Format(filepath.Join(outDir, "main.fst"))
		if err != nil {
			t.Fatal(err)
		}
		// bs, err := json.MarshalIndent(main, "", "\t")
		// bs, err := json.Marshal(main)
		// if err != nil {
		// 	t.Fatal(err)
		// }
		// err = os.WriteFile("testdata/"+filepath.Dir(path)+"/test.json", bs, 0777)
		// if err != nil {
		// 	t.Fatal(err)
		// }
		// err = json.Unmarshal(bs, main)
		// if err != nil {
		// 	t.Fatal(err)
		// }
	}
}
