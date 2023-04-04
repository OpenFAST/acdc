package fio

// Named Types
const (
	Title        Type = "Title"
	Text         Type = "Text"
	Heading      Type = "Heading"
	String       Type = "String"
	Bool         Type = "Bool"
	BoolDefault  Type = "BoolDefault"
	FloatDefault Type = "FloatDefault"
	FloatAll     Type = "FloatAll"
	Float        Type = "Float"
	Floats       Type = "Floats"
	Int          Type = "Int"
	IntDefault   Type = "IntDefault"
	Ints         Type = "Ints"
	Path         Type = "Path"
	Paths        Type = "Paths"
	OutList      Type = "OutList"
	OutList2     Type = "OutList2"
	Table        Type = "Table"
	BDStations   Type = "BDStations"
)

// Field defines a member of a schema.
type Field struct {
	Name   string `json:"Key"`
	Type   Type   `json:"Type"`
	Desc   string `json:"Comment"`
	Unit   string `json:"Unit"`
	Num    string `json:"Num"`
	Value  any    `json:"Value"`
	Values []any  `json:"Values"`
	Table  [][]any
	Set    bool `json:"Set"`

	// Extra type data
	SwitchOptions   []Option `json:"Options"`
	TableHeaderSize int      `json:"HeaderSize"`
	TableColumns    []Column `json:"Columns"`
	PathFileType    string   `json:"PathFileType"`
}

//------------------------------------------------------------------------------
// Switch
//------------------------------------------------------------------------------

type Option struct {
	Value any `json:"Value"`
	Label any `json:"Label"`
}

//------------------------------------------------------------------------------
// Table
//------------------------------------------------------------------------------

type Column struct {
	Name    string `json:"Key"`
	Type    Type   `json:"Type"`
	Comment string `json:"Comment"`
	Unit    string `json:"Unit"`
}
