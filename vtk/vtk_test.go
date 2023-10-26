package vtk_test

import (
	"acdc/vtk"
	"testing"
)

func TestLoadVTK(t *testing.T) {

	testFiles := []string{
		// "testdata/BD_BldMotion1.001.vtp",
		"testdata/ED_Hub.001.vtp",
		"testdata/ED_Nacelle.001.vtp",
		"testdata/ED_TailFin.001.vtp",
		"testdata/ED_TowerLn2Mesh_motion.001.vtp",
		"testdata/GroundSurface.vtp",
	}

	for _, testFile := range testFiles {
		_, err := vtk.LoadVTK(testFile)
		if err != nil {
			t.Fatal(err)
		}
	}

}
