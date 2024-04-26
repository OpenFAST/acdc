package viz_test

import (
	"acdc/viz"
	"testing"
)

func TestBuildModeViz(t *testing.T) {

	data, err := viz.BuildModeViz([]string{
		"testdata/03_NREL_5MW-ED.Mode1.LinTime1.ED_TowerLn2Mesh_motion.001.vtp",
		"testdata/03_NREL_5MW-ED.Mode1.LinTime1.ED_TowerLn2Mesh_motion.002.vtp",
		"testdata/03_NREL_5MW-ED.Mode1.LinTime1.ED_TowerLn2Mesh_motion.003.vtp",
		"testdata/03_NREL_5MW-ED.Mode1.LinTime1.ED_TowerLn2Mesh_motion.004.vtp",
		"testdata/03_NREL_5MW-ED.Mode1.LinTime1.ED_TowerLn2Mesh_motion.005.vtp",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", data)
}
