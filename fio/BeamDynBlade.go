package fio

var BeamDynBlade = NewFile("BeamDynBlade", []Field{
	{Type: Heading, Desc: "BeamDynBlade Input File"},
	{Name: "Title", Type: Title},
	{Type: Heading, Desc: "Blade Parameters"},
	{Name: "station_total", Type: Int, Desc: `Number of blade input stations (-)`},
	{Name: "damp_type", Type: Int, Desc: `Damping type: 0: no damping; 1: damped`},
	{Type: Heading, Desc: "Damping Coefficient"},
	{Name: "mu", Type: Table, TableHeaderSize: 2, TableColumns: []Column{
		{Name: "mu1", Type: Float, Unit: "-"},
		{Name: "mu2", Type: Float, Unit: "-"},
		{Name: "mu3", Type: Float, Unit: "-"},
		{Name: "mu4", Type: Float, Unit: "-"},
		{Name: "mu5", Type: Float, Unit: "-"},
		{Name: "mu6", Type: Float, Unit: "-"},
	}},
	{Type: Heading, Desc: "Distributed Properties"},
	{Name: "Stations", Type: BDStations, Num: "station_total"},
})
