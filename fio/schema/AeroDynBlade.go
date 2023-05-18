package schema

var AeroDynBlade = RegisterSchemas("AeroDynBlade", []Schema{
	{Type: Header, Desc: "AeroDynBlade Input File"},
	{Name: "Title", Type: Title},
	{Type: Header, Desc: "Blade Properties"},
	{Name: "NumBlNds", Type: Int, Desc: "Number of blade nodes used in the analysis"},
	{Name: "BlNds", Type: Table, NumVar: "NumBlNds",
		TableHeaderSize: 2,
		TableColumns: []Column{
			{Name: "BlSpn", Type: Float, Unit: "m"},
			{Name: "BlCrvAC", Type: Float, Unit: "m"},
			{Name: "BlSwpAC", Type: Float, Unit: "m"},
			{Name: "BlCrvAng", Type: Float, Unit: "deg"},
			{Name: "BlTwist", Type: Float, Unit: "deg"},
			{Name: "BlChord", Type: Float, Unit: "m"},
			{Name: "BlAFID", Type: Int, Unit: "-"},
			{Name: "BlCb", Type: Float, Unit: "-"},
			{Name: "BlCenBn", Type: Float, Unit: "m"},
			{Name: "BlCenBt", Type: Float, Unit: "m"},
		}},
})
