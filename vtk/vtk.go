package vtk

import (
	"encoding/xml"
	"os"
	"strconv"
	"strings"
)

type VTKFile struct {
	Type      string   `xml:"type,attr"`
	Version   string   `xml:"version,attr"`
	ByteOrder string   `xml:"byte_order,attr"`
	PolyData  PolyData `xml:"PolyData"`
}

type PolyData struct {
	Piece Piece
}

type Piece struct {
	NumberOfPoints int `xml:",attr"`
	NumberOfVerts  int `xml:",attr"`
	NumberOfLines  int `xml:",attr"`
	NumberOfStrips int `xml:",attr"`
	NumberOfPolys  int `xml:",attr"`
	Points         Points
	Lines          Lines
	Polys          Polys
}

type Points struct {
	DataArray DataArray `xml:"DataArray"`
}

type Lines struct {
	DataArray []DataArray `xml:"DataArray"`
}

type Polys struct {
	DataArray []DataArray `xml:"DataArray"`
}

type DataArray struct {
	Name               string      `xml:",attr"`
	Type               string      `xml:"type,attr"`
	NumberOfComponents int         `xml:",attr"`
	Format             string      `xml:"format,attr"`
	RawValues          string      `xml:",chardata"`
	ValuesF32          []float32   `xml:"-"`
	MatrixF32          [][]float32 `xml:"-"`
	ValuesInt32        []int32     `xml:"-"`
	Connectivity       []int32     `xml:"-"`
	Offsets            []int32     `xml:"-"`
}

func (da *DataArray) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type Tmp DataArray
	if err := d.DecodeElement((*Tmp)(da), &start); err != nil {
		return err
	}

	// Split the raw values into a slice of strings (space-separated values)
	valueStrings := strings.Fields(da.RawValues)
	numValues := len(valueStrings)

	// Switch based on data type
	switch da.Type {
	case "Float32":
		switch da.NumberOfComponents {
		case 3:
			values := make([]float32, len(valueStrings))
			for i, v := range valueStrings {
				value, err := strconv.ParseFloat(v, 32)
				if err != nil {
					return err
				}
				values[i] = float32(value)
			}
			da.MatrixF32 = make([][]float32, int(numValues/3))
			for i := range da.MatrixF32 {
				da.MatrixF32[i], values = values[:3:3], values[3:]
			}
		}
	case "Int32":
		values := make([]int32, len(valueStrings))
		for i, v := range valueStrings {
			value, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				return err
			}
			values[i] = int32(value)
		}
		switch strings.ToLower(da.Name) {
		case "connectivity":
			da.Connectivity = values
		case "offsets":
			da.Offsets = values
		default:
			da.ValuesInt32 = values
		}
	}

	return nil
}

func LoadVTK(path string) (*VTKFile, error) {

	bs, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	vf := &VTKFile{}

	if err = xml.Unmarshal(bs, vf); err != nil {
		return nil, err
	}

	return vf, nil
}
