package viz

import (
	"encoding/xml"
	"fmt"
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

type OrientationVectors struct {
	X [][]float32 `json:"x"` // nx3 matrix for X orientation vectors
	Y [][]float32 `json:"y"` // nx3 matrix for Y orientation vectors
	Z [][]float32 `json:"z"` // nx3 matrix for Z orientation vectors
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
	PointData      PointData
}

type Points struct {
	DataArray DataArray `xml:"DataArray"`
}

type Lines struct {
	DataArray []DataArray `xml:"DataArray"`
}

type PointData struct {
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

	fmt.Println("decode element: ", d)

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

func LoadVTK(path string) (*VTKFile, *OrientationVectors, error) {

	bs, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	vf := &VTKFile{}

	if err = xml.Unmarshal(bs, vf); err != nil {
		return nil, nil, err
	}

	// Debug: Print the structure
	fmt.Printf("NumberOfPoints: %d\n", vf.PolyData.Piece.NumberOfPoints)
	fmt.Printf("Points DataArray Type: %s, Components: %d\n", vf.PolyData.Piece.Points.DataArray.Type, vf.PolyData.Piece.Points.DataArray.NumberOfComponents)
	fmt.Printf("Points DataArray Values: %v\n", vf.PolyData.Piece.Points.DataArray.MatrixF32)

	// Create orientation vectors
	orientationVectors := &OrientationVectors{}

	// Debug: Print all available DataArrays
	fmt.Printf("Number of PointData arrays: %d\n", len(vf.PolyData.Piece.PointData.DataArray))
	for i, da := range vf.PolyData.Piece.PointData.DataArray {
		fmt.Printf("DataArray[%d]: Name='%s', Type='%s', Components=%d\n", i, da.Name, da.Type, da.NumberOfComponents)
	}

	// Find DataArray that has the name of "OrientationX"
	for _, da := range vf.PolyData.Piece.PointData.DataArray {
		if da.Name == "OrientationX" {
			fmt.Println("Found OrientationX:", da.MatrixF32)
			orientationVectors.X = da.MatrixF32
		} else if da.Name == "OrientationY" {
			fmt.Println("Found OrientationY:", da.MatrixF32)
			orientationVectors.Y = da.MatrixF32
		} else if da.Name == "OrientationZ" {
			fmt.Println("Found OrientationZ:", da.MatrixF32)
			orientationVectors.Z = da.MatrixF32
		} else if da.Name == "OrientationZ" {
			fmt.Println("Found OrientationZ:", da.MatrixF32)
			orientationVectors.Z = da.MatrixF32
		}
	}

	return vf, orientationVectors, nil
}
