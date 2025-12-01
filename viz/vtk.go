package viz

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
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

type OrientationMatrix struct {
	X []float32 `json:"x"` // 3x3 matrix for X orientation vectors
	Y []float32 `json:"y"` // 3x3 matrix for Y orientation vectors
	Z []float32 `json:"z"` // 3x3 matrix for Z orientation vectors
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

	// fmt.Println("decode element: ", d)

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

func LoadVTK(path string) (*VTKFile, *VTKFile, error) {

	bs, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	vf := &VTKFile{}

	if err = xml.Unmarshal(bs, vf); err != nil {
		return nil, nil, err
	}

	// Skip for non-blade files for local coords conversion
	if !strings.Contains(filepath.Base(path), "Blade") {
		return vf, nil, nil
	}

	local_vf, err := Global2Local(vf)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to convert global coordinates to local: %w", err)
	}

	return vf, local_vf, nil
}

func GetOrientations(vf *VTKFile) (*OrientationVectors, *OrientationMatrix, error) {
	// Create orientation vectors
	ov := &OrientationVectors{}

	// Find Orientation vectors
	for _, da := range vf.PolyData.Piece.PointData.DataArray {
		if da.Name == "OrientationX" {
			ov.X = da.MatrixF32
		} else if da.Name == "OrientationY" {
			ov.Y = da.MatrixF32
		} else if da.Name == "OrientationZ" {
			ov.Z = da.MatrixF32
		}
	}

	// Create a new OrientationMatrix
	om := &OrientationMatrix{
		X: make([]float32, 3),
		Y: make([]float32, 3),
		Z: make([]float32, 3),
	}

	// Fill the OrientationMatrix with the first 3 vectors from each orientation
	om.X = ov.X[0]
	om.Y = ov.Y[0]
	om.Z = ov.Z[0]

	return ov, om, nil
}

func Global2Local(vf *VTKFile) (*VTKFile, error) {

	// Copy vf (Deep Copy -- so that it works independently)
	var local_vf *VTKFile
	err := DeepCopy(&vf, &local_vf)
	local_coords := local_vf.PolyData.Piece.Points.DataArray.MatrixF32
	// fmt.Println("\nLocal coordinates:", local_coords)

	// Get Orientation vectors and matrices
	ov, om, err := GetOrientations(local_vf)
	if err != nil {
		return nil, fmt.Errorf("Failed to extract orientation vectors and matrices: %w", err)
	}
	// fmt.Println("\nOrientation Matrix:", om)

	// Translational/Rotational operations for the points
	local_coords = TranslateMatrix(local_coords, local_coords[0]) // Translate by the first point -- so that first point will be moved to the origin
	// fmt.Println("\nLocal coordinates:", local_coords)

	transposed_om := TransposeMatrix(om)
	// fmt.Println("\nTransposed Orientation Matrix:", transposed_om)

	local_coords = DotProduct(local_coords, transposed_om) // Rotate by the first orientation vector
	local_vf.PolyData.Piece.Points.DataArray.MatrixF32 = local_coords
	// fmt.Println("\nLocal coordinates from LoadVTK() :", local_coords)

	// Rotational operations for the Orientation vectors
	ov.X = DotProduct(ov.X, transposed_om)
	ov.Y = DotProduct(ov.Y, transposed_om)
	ov.Z = DotProduct(ov.Z, transposed_om)

	// fmt.Println("\nOrientationX:", ov.X)
	// fmt.Println("\nOrientationY:", ov.Y)
	// fmt.Println("\nOrientationZ:", ov.Z)

	return local_vf, nil
}

func TranslateMatrix(points [][]float32, translation []float32) [][]float32 {
	result := make([][]float32, len(points))
	for i, point := range points {
		resPoint := make([]float32, 3)
		for j := 0; j < 3; j++ {
			resPoint[j] = point[j] + (-translation[j])
		}
		result[i] = resPoint
	}
	return result
}

func DotProduct(vectors [][]float32, matrix [][]float32) [][]float32 {
	result := make([][]float32, len(vectors))
	for i, vec := range vectors {
		resVec := make([]float32, 3)
		for j := 0; j < 3; j++ {
			for k := 0; k < 3; k++ {
				resVec[j] += vec[k] * matrix[k][j]
			}
		}
		result[i] = resVec
	}
	return result
}

func TransposeMatrix(om *OrientationMatrix) [][]float32 {
	// Transpose the orientation matrix
	if len(om.X) == 0 {
		return nil
	}
	transposed := make([][]float32, 3)
	transposed[0] = []float32{om.X[0], om.Y[0], om.Z[0]}
	transposed[1] = []float32{om.X[1], om.Y[1], om.Z[1]}
	transposed[2] = []float32{om.X[2], om.Y[2], om.Z[2]}

	return transposed
}

func DeepCopy(src, dst interface{}) error {
	bytes, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, dst)
}
