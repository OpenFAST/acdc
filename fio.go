package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

type PostParser interface {
	PostParse() error
}

//------------------------------------------------------------------------------
// Variable Types
//------------------------------------------------------------------------------

type FieldBase struct {
	Name string `json:"Name"`
	Type string `json:"Type"`
	Desc string `json:"Desc"`
	Line int    `json:"Line"`
}

type Path struct {
	FieldBase
	Value    string `json:"Value"`
	FileType string `json:"FileType"`
	Root     bool   `json:"Root"`
}

type Paths struct {
	FieldBase
	Value     []string `json:"Value"`
	FileType  string   `json:"FileType"`
	Condensed bool     `json:"Condensed"`
}

type Bool struct {
	FieldBase
	Value bool `json:"Value"`
}

type Integer struct {
	FieldBase
	Value int  `json:"Value"`
	Size  bool `json:"Size"`
}

type Real struct {
	FieldBase
	Value float64 `json:"Value"`
}

type Reals struct {
	FieldBase
	Value []float64 `json:"Value"`
}

type String struct {
	FieldBase
	Value string `json:"Value"`
}

//------------------------------------------------------------------------------
// Files Type
//------------------------------------------------------------------------------

// Files structure contains slices of all file types
// file types must be slices for the parsing and writing code to work
type Files struct {
	Main        []Main            `json:"Main"`
	ElastoDyn   []ElastoDyn       `json:"ElastoDyn"`
	BeamDyn     []BeamDyn         `json:"BeamDyn"`
	AeroDyn     []AeroDyn         `json:"AeroDyn"`
	AeroDyn14   []AeroDyn14       `json:"AeroDyn14"`
	HydroDyn    []HydroDyn        `json:"HydroDyn"`
	ServoDyn    []ServoDyn        `json:"ServoDyn"`
	InflowWind  []InflowWind      `json:"InflowWind"`
	OLAF        []OLAF            `json:"OLAF"`
	Misc        []Misc            `json:"Misc"`
	StControl   []StControl       `json:"StControl"`
	AirfoilInfo []AirfoilInfo     `json:"AirfoilInfo"`
	PathMap     map[string]string `json:"-"`
}

// NewFiles returns the Model structure with all slices initialized to empty
func NewFiles() *Files {
	return &Files{
		Main:        []Main{},
		ElastoDyn:   []ElastoDyn{},
		BeamDyn:     []BeamDyn{},
		AeroDyn:     []AeroDyn{},
		AeroDyn14:   []AeroDyn14{},
		HydroDyn:    []HydroDyn{},
		InflowWind:  []InflowWind{},
		OLAF:        []OLAF{},
		ServoDyn:    []ServoDyn{},
		StControl:   []StControl{},
		Misc:        []Misc{},
		AirfoilInfo: []AirfoilInfo{},
		PathMap:     map[string]string{},
	}
}

// Add adds a file structure of the given fileType to the Model structure
// and returns a pointer to it. fileType must match the name of a field
// in the Model struct or it will be treated as a Misc type file.
func (m *Files) Add(fileType string) any {

	mVal := reflect.Indirect(reflect.ValueOf(m))

	slice := mVal.FieldByName(fileType)
	if !slice.IsValid() {
		slice = mVal.FieldByName("Misc")
	}

	newVal := reflect.New(slice.Type().Elem()).Elem()
	slice.Set(reflect.Append(slice, newVal))
	return slice.Index(slice.Len() - 1).Addr().Interface()
}

func (fs *Files) Copy() (*Files, error) {
	bs, err := json.Marshal(fs)
	if err != nil {
		return nil, err
	}
	newFiles := &Files{}
	if err := json.Unmarshal(bs, newFiles); err != nil {
		return nil, err
	}
	return newFiles, err
}

//------------------------------------------------------------------------------
// File Types
//------------------------------------------------------------------------------

type FileBase struct {
	Name  string   `json:"Name"`
	Type  string   `json:"Type"`
	Lines []string `json:"Lines"`
}

type Misc struct {
	FileBase
}

type Main struct {
	FileBase
	TMax        Real    `json:"TMax"`
	DT          Real    `json:"DT"`
	CompElast   Integer `json:"CompElast"`
	CompInflow  Integer `json:"CompInflow"`
	CompAero    Integer `json:"CompAero"`
	CompServo   Integer `json:"CompServo"`
	CompHydro   Integer `json:"CompHydro"`
	CompSub     Integer `json:"CompSub"`
	CompMooring Integer `json:"CompMooring"`
	CompIce     Integer `json:"CompIce"`
	MHK         Integer `json:"MHK"`
	Gravity     Real    `json:"Gravity"`
	EDFile      Path    `json:"EDFile" ftype:"ElastoDyn"`
	BDBldFile1  Path    `json:"BDBldFile1" key:"BDBldFile(1)" ftype:"BeamDyn"`
	BDBldFile2  Path    `json:"BDBldFile2" key:"BDBldFile(2)" ftype:"BeamDyn"`
	BDBldFile3  Path    `json:"BDBldFile3" key:"BDBldFile(3)" ftype:"BeamDyn"`
	InflowFile  Path    `json:"InflowFile" ftype:"InflowWind"`
	AeroFile    Path    `json:"AeroFile" ftype:"AeroDyn"`
	ServoFile   Path    `json:"ServoFile" ftype:"ServoDyn"`
	HydroFile   Path    `json:"HydroFile" ftype:"HydroDyn"`
	SubFile     Path    `json:"SubFile" ftype:"SubDyn"`
	MooringFile Path    `json:"MooringFile" ftype:"Misc"`
	IceFile     Path    `json:"IceFile" ftype:"Misc"`
	SttsTime    Real    `json:"SttsTime"`
	OutFileFmt  Integer `json:"OutFileFmt"`
	OutFmt      String  `json:"OutFmt"`
	Linearize   Bool    `json:"Linearize"`
	CalcSteady  Bool    `json:"CalcSteady"`
	TrimCase    Integer `json:"TrimCase"`
	TrimTol     Real    `json:"TrimTol"`
	TrimGain    Real    `json:"TrimGain"`
	Twr_Kdmp    Real    `json:"Twr_Kdmp"`
	Bld_Kdmp    Real    `json:"Bld_Kdmp"`
	NLinTimes   Integer `json:"NLinTimes"`
	LinTimes    Reals   `json:"LinTimes" num:"NLinTimes"`
	LinInputs   Integer `json:"LinInputs"`
	LinOutputs  Integer `json:"LinOutputs"`
	LinOutJac   Bool    `json:"LinOutJac"`
	LinOutMod   Bool    `json:"LinOutMod"`
	WrVTK       Integer `json:"WrVTK"`
	VTK_type    Integer `json:"VTK_type"`
	VTK_fps     Integer `json:"VTK_fps"`
}

func (m *Main) PostParse() error {

	switch m.CompAero.Value {
	case 0:
		m.AeroFile.FileType = ""
	case 1:
		m.AeroFile.FileType = "AeroDyn14"
	case 2:
		m.AeroFile.FileType = "AeroDyn"
	}

	m.NLinTimes.Size = false

	return nil
}

type AeroDyn struct {
	FileBase
	WakeMod           Integer `json:"WakeMod"`
	AFAeroMod         Integer `json:"AFAeroMod"`
	TwrPotent         Integer `json:"TwrPotent"`
	TwrShadow         Integer `json:"TwrShadow"`
	FrozenWake        Bool    `json:"FrozenWake"`
	SkewMod           Integer `json:"SkewMod"`
	OLAFInputFileName Path    `json:"OLAFInputFileName" ftype:"OLAF"`
	NumAFfiles        Integer `json:"NumAFfiles"`
	AFNames           Paths   `json:"AFNames" num:"NumAFfiles" ftype:"AirfoilInfo"`
	ADBlFile1         Path    `json:"ADBlFile1" key:"ADBlFile(1)" ftype:"Misc"` // AeroDynBlade
	ADBlFile2         Path    `json:"ADBlFile2" key:"ADBlFile(2)" ftype:"Misc"` // AeroDynBlade
	ADBlFile3         Path    `json:"ADBlFile3" key:"ADBlFile(3)" ftype:"Misc"` // AeroDynBlade
	TFinFile          Path    `json:"TFinFile"`
}

type AirfoilInfo struct {
	FileBase
	BL_File Path `json:"BL_File"`
}

type OLAF struct {
	FileBase
	PrescribedCircFile Path `json:"PrescribedCircFile"`
}

type AeroDyn14 struct {
	FileBase
	NumFoil Integer `json:"NumFoil"`
	FoilNm  Paths   `json:"FoilNm" num:"NumFoil"`
}

type BeamDyn struct {
	FileBase
	RotStates Bool `json:"RotStates"`
	BldFile   Path `json:"BldFile" ftype:"Misc"` // BeamDynBlade
}

type ElastoDyn struct {
	FileBase
	FlapDOF1 Bool    `json:"FlapDOF1"`
	FlapDOF2 Bool    `json:"FlapDOF2"`
	EdgeDOF  Bool    `json:"EdgeDOF"`
	TeetDOF  Bool    `json:"TeetDOF"`
	DrTrDOF  Bool    `json:"DrTrDOF"`
	GenDOF   Bool    `json:"GenDOF"`
	YawDOF   Bool    `json:"YawDOF"`
	TwFADOF1 Bool    `json:"TwFADOF1"`
	TwFADOF2 Bool    `json:"TwFADOF2"`
	TwSSDOF1 Bool    `json:"TwSSDOF1"`
	TwSSDOF2 Bool    `json:"TwSSDOF2"`
	BlPitch1 Real    `json:"BlPitch1" key:"BlPitch(1)"`
	BlPitch2 Real    `json:"BlPitch2" key:"BlPitch(2)"`
	BlPitch3 Real    `json:"BlPitch3" key:"BlPitch(3)"`
	RotSpeed Real    `json:"RotSpeed"`
	NumBl    Integer `json:"NumBl"`
	ShftTilt Real    `json:"ShftTilt"`
	BldFile1 Path    `json:"BldFile1" key:"BldFile(1)" ftype:"Misc"` // ElastoDynBlade
	BldFile2 Path    `json:"BldFile2" key:"BldFile(2)" ftype:"Misc"` // ElastoDynBlade
	BldFile3 Path    `json:"BldFile3" key:"BldFile(3)" ftype:"Misc"` // ElastoDynBlade
	TwrFile  Path    `json:"TwrFile" ftype:"Misc"`                   // ElastoDynTower
}

type HydroDyn struct {
	FileBase
	WaveMod  Integer `json:"WaveMod"`
	ExctnMod Integer `json:"ExctnMod"`
	PotFile  Path    `json:"PotFile" ftype:"Misc"`
}

func (h *HydroDyn) PostParse() error {
	h.PotFile.Root = true
	return nil
}

type InflowWind struct {
	FileBase
	WindType       Integer `json:"WindType"`
	PropagationDir Real    `json:"PropagationDir"`
	VFlowAng       Real    `json:"VFlowAng"`
	HWindSpeed     Real    `json:"HWindSpeed"`
	PLExp          Real    `json:"PLExp"`
}

func (IfW *InflowWind) PostParse() error {
	IfW.WindType.Value = 1
	return nil
}

type ServoDyn struct {
	FileBase
	PCMode    Integer `json:"PCMode"`
	VSContrl  Integer `json:"VSContrl"`
	VS_RtGnSp Real    `json:"VS_RtGnSp"`
	VS_RtTq   Real    `json:"VS_RtTq"`
	VS_Rgn2K  Real    `json:"VS_Rgn2K"`
	VS_SlPc   Real    `json:"VS_SlPc"`
	HSSBrMode Integer `json:"HSSBrMode"`
	YCMode    Integer `json:"YCMode"`
	NumBStC   Integer `json:"NumBStC"`
	BStCfiles Paths   `json:"BStCfiles" num:"NumBStC" ftype:"StControl"`
	NumNStC   Integer `json:"NumNStC"`
	NStCfiles Paths   `json:"NStCfiles" num:"NumNStC" ftype:"StControl"`
	NumTStC   Integer `json:"NumTStC"`
	TStCfiles Paths   `json:"TStCfiles" num:"NumTStC" ftype:"StControl"`
	NumSStC   Integer `json:"NumSStC"`
	SStCfiles Paths   `json:"SStCfiles" num:"NumSStC" ftype:"StControl"`
}

func (sd *ServoDyn) PostParse() error {
	sd.BStCfiles.Condensed = true
	sd.NStCfiles.Condensed = true
	sd.TStCfiles.Condensed = true
	sd.SStCfiles.Condensed = true
	return nil
}

type StControl struct {
	FileBase
	PrescribedForcesFile Path `json:"PrescribedForcesFile"`
}

//------------------------------------------------------------------------------
// Parsing
//------------------------------------------------------------------------------

func ParseFiles(MainPath string) (*Files, error) {

	// Initialize files structure
	files := NewFiles()

	// Parse Main file and all files it references recursively
	files.Main = []Main{{}}
	if err := files.parseFile(MainPath, &files.Main[0]); err != nil {
		return nil, err
	}

	// Add main path to path map
	files.PathMap[MainPath] = filepath.Base(MainPath)

	// Return files
	return files, nil
}

var parensReplacer = strings.NewReplacer("(", "", ")", "")

func (fs *Files) parseFile(path string, s any) error {

	lines, err := readLines(path)
	if err != nil {
		return err
	}
	numLines := len(lines)

	sVal := reflect.ValueOf(s).Elem()
	sTyp := sVal.Type()

	fb := sVal.FieldByName("FileBase").Addr().Interface().(*FileBase)
	fb.Name = filepath.Base(path)
	fb.Type = sTyp.Name()
	fb.Lines = lines

	// Loop through fields in struct
	for i := 1; i < sVal.NumField(); i++ {

		// Get field name
		fieldType := sTyp.Field(i)
		fieldName := fieldType.Name
		if key, ok := fieldType.Tag.Lookup("key"); ok {
			fieldName = key
		}
		fieldNameLower := strings.ToLower(fieldName)
		fieldNameLowerNoParens := parensReplacer.Replace(fieldNameLower)
		if fieldNameLowerNoParens == fieldNameLower {
			fieldNameLowerNoParens = ""
		}

		// Create backup of lines to search
		linesSave := lines

		// Initialize field parsed to false
		fieldParsed := false

		// Loop through lines
		for len(lines) > 0 {
			line := lines[0]
			lines = lines[1:]

			// Remove comment from line and trim whitespace
			found, desc := false, ""
			if line, desc, found = strings.Cut(line, "- "); !found {
				line, desc, _ = strings.Cut(line, "! ")
			}
			line = strings.TrimSpace(line)

			// Find index of field name in line
			lineLower := strings.ToLower(line)
			j := strings.LastIndex(lineLower, fieldNameLower)
			if j == -1 && fieldNameLowerNoParens != "" {
				j = strings.LastIndex(lineLower, fieldNameLowerNoParens)
			}

			// Field name not found in line, continue
			if j == -1 {
				continue
			}

			// Split line into field while respecting quotes
			quoted := false
			values := strings.FieldsFunc(line[:j], func(r rune) bool {
				if r == '"' {
					quoted = !quoted
				}
				return !quoted && (unicode.IsSpace(r) || r == ',')
			})

			// Get field value
			fieldVal := sVal.Field(i)

			// Set field base info
			if base, ok := fieldVal.FieldByName("FieldBase").Addr().Interface().(*FieldBase); ok {
				base.Name = fieldName
				base.Type = fieldType.Type.Name()
				base.Desc = desc
				base.Line = numLines - len(lines)
			}

			// Get size based on 'num' key if defined for field
			var numField *Integer
			if numFieldName, ok := fieldType.Tag.Lookup("num"); ok {
				numFieldVal := sVal.FieldByName(numFieldName)
				if numFieldVal.IsZero() {
					return fmt.Errorf("unknown field for num of paths in '%s'", fieldType.Name)
				}
				numField, ok = numFieldVal.Addr().Interface().(*Integer)
				if !ok {
					return fmt.Errorf("field for num of paths in '%s' is not an Int", fieldType.Name)
				}
				numField.Size = true
			}

			// Switch based on field type
			switch v := fieldVal.Addr().Interface().(type) {
			case *Path:
				v.Value = strings.Trim(values[0], `"`)
			case *Paths:
				if numField == nil {
					return fmt.Errorf("number of paths in '%s' not specified", v.Name)
				}
				for _, value := range values[:min(len(values), numField.Value)] {
					v.Value = append(v.Value, strings.Trim(value, `"`))
				}
				for i := len(v.Value); i < numField.Value; i++ {
					line, lines = lines[0], lines[1:]
					v.Value = append(v.Value, strings.Trim(strings.TrimSpace(line), `"`))
				}
			case *Bool:
				v.Value, err = strconv.ParseBool(values[0])
			case *String:
				v.Value = values[0]
			case *Integer:
				v.Value, err = strconv.Atoi(values[0])
			case *Real:
				v.Value, err = strconv.ParseFloat(values[0], 64)
			case *Reals:
				if numField == nil {
					return fmt.Errorf("number of paths in '%s' not specified", v.Name)
				}
				for i := range values {
					f, e := strconv.ParseFloat(values[i], 64)
					if e != nil {
						err = e
						break
					}
					v.Value = append(v.Value, f)
				}
			}
			if err != nil {
				return fmt.Errorf("error parsing field '%s' in file '%s': %w",
					fieldName, path, err)
			}

			// Variable parsed, from line, continue to next variable
			fieldParsed = true
			break
		}

		// Return error if field not found
		// if !fieldParsed {
		// 	return fmt.Errorf("error parsing file '%s', field '%s' not found", path, fieldName)
		// }

		// If field was not parsed, restore lines and go to next field
		if !fieldParsed {
			lines = linesSave
		}
	}

	// If file type has a PostParse method, call it
	if pp, ok := s.(PostParser); ok {
		if err := pp.PostParse(); err != nil {
			return fmt.Errorf("error parsing '%s': %w", path, err)
		}
	}

	// Get directory of current file
	dir := filepath.Dir(path)

	// Loop through fields in struct and parse paths
	for i := 1; i < sVal.NumField(); i++ {

		// Get field value
		fieldVal := sVal.Field(i)
		fieldType := sTyp.Field(i)

		// If field is a path
		if p, ok := fieldVal.Addr().Interface().(*Path); ok {

			// If file type not specified, get it from tag
			if p.FileType == "" {
				p.FileType = fieldType.Tag.Get("ftype")
			}

			// Get path to file for reading, if not an absolute path, prepend
			// directory of current file
			subpath := filepath.Clean(p.Value)
			if !filepath.IsAbs(subpath) {
				subpath = filepath.Join(dir, subpath)
			}

			// If path has already been read, change path to name and skip reading
			if name, ok := fs.PathMap[subpath]; ok {
				p.Value = name
				continue
			}

			// If this is a root path
			if p.Root {

				// List of files matching pattern
				matches, err := filepath.Glob(subpath + ".*")
				if err != nil {
					return fmt.Errorf("invalid root path: %w", err)
				}

				// Loop through matches
				for _, match := range matches {

					// Add file type to files, get structure to parse into
					ss := fs.Add(p.FileType)

					// Parse file into struct
					if err := fs.parseFile(match, ss); err != nil {
						return fmt.Errorf("error parsing '%s': %w", subpath, err)
					}

					// Update path map to indicate that file has been read
					fs.PathMap[match] = filepath.Base(match)
				}

			} else {

				// If path doesn't exist, skip reading
				if stat, err := os.Stat(subpath); err != nil || stat.IsDir() {
					p.Value = "FileNotFound"
					continue
				}

				// Add file type to files, get structure to parse into
				ss := fs.Add(p.FileType)

				// Parse file
				if err := fs.parseFile(subpath, ss); err != nil {
					return fmt.Errorf("error parsing '%s': %w", subpath, err)
				}

				// Update path map to indicate that file has been read
				fs.PathMap[subpath] = filepath.Base(subpath)
			}

			// Change path to name of file
			p.Value = filepath.Base(subpath)
		}

		// If field is a path
		if p, ok := fieldVal.Addr().Interface().(*Paths); ok {

			// If file type not specified, get from tag, skip if still unknown
			if p.FileType == "" {
				p.FileType = fieldType.Tag.Get("ftype")
			}

			// Loop through path values
			for i := range p.Value {

				// Get path to file for reading, if not an absolute path, prepend
				// directory of current file
				subpath := filepath.Clean(p.Value[i])
				if !filepath.IsAbs(subpath) {
					subpath = filepath.Join(dir, subpath)
				}

				// If path doesn't exist, skip reading
				if stat, err := os.Stat(subpath); err != nil || stat.IsDir() {
					p.Value[i] = "FileNotFound"
					continue
				}

				// If path has already been read, change path to name and skip reading
				if name, ok := fs.PathMap[subpath]; ok {
					p.Value[i] = name
					continue
				}

				// Add file type to model, get structure to parse into
				ss := fs.Add(p.FileType)

				// Parse file
				if err := fs.parseFile(subpath, ss); err != nil {
					return fmt.Errorf("error parsing '%s': %w", subpath, err)
				}

				// Change path to name of file
				p.Value[i] = filepath.Base(subpath)

				// Update path map to indicate that file has been read
				fs.PathMap[subpath] = p.Value[i]
			}
		}
	}

	return nil
}

func readLines(path string) ([]string, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := make([]string, 0, 256)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// If additional file is to be read and inserted
		if strings.HasPrefix(line, "@") {
			subPath := strings.Trim(strings.Fields(line[1:])[0], `"`)
			subLines, err := readLines(filepath.Join(filepath.Dir(path), subPath))
			if err != nil {
				return nil, err
			}
			lines = append(lines, subLines...)
			continue
		}

		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning text: %w", err)
	}
	return lines, nil
}

//------------------------------------------------------------------------------
// Write
//------------------------------------------------------------------------------

func (m *Files) Write(dir, prefix string) error {

	val := reflect.ValueOf(m).Elem()

	// Loop through fields in model
	for i := 0; i < val.NumField(); i++ {

		// Get field
		field := val.Field(i)

		// Skip fields that aren't a slice
		if field.Kind() != reflect.Slice {
			continue
		}

		// Loop through file structures in slice and write to file
		for j := 0; j < field.Len(); j++ {
			if err := writeFile(field.Index(j).Addr().Interface(), dir, prefix); err != nil {
				return fmt.Errorf("error writing %s file: %w", val.Type().Field(i).Name, err)
			}
		}
	}

	return nil
}

func writeFile(s any, dir, prefix string) error {

	sVal := reflect.Indirect(reflect.ValueOf(s))

	// Get file base data
	fb := sVal.FieldByName("FileBase").Interface().(FileBase)

	// Create path
	path := filepath.Join(dir, prefix+fb.Name)
	if fb.Type == "Misc" {
		path = filepath.Join(dir, fb.Name)
	}

	// Get lines in file
	lines := fb.Lines

	// Loop through fields
	for i := 1; i < sVal.NumField(); i++ {

		// Get field value
		fieldVal := sVal.Field(i)

		// Switch based on field type, field that has multiple values
		size := 0
		name := ""
		switch v := fieldVal.Interface().(type) {
		case Paths:
			size = len(v.Value)
			name = v.Name
		case Reals:
			size = len(v.Value)
			name = v.Name
		default:
			continue
		}

		// Get name of field that specifies number of items
		fieldTyp := sVal.Type().Field(i)
		numFieldName, ok := fieldTyp.Tag.Lookup("num")
		if !ok {
			return fmt.Errorf("number of paths in '%s' not specified", name)
		}

		// Get number field value
		numFieldVal := sVal.FieldByName(numFieldName)
		a := numFieldVal.Interface().(Integer)
		if numFieldVal.IsZero() {
			return fmt.Errorf("unknown field for num of items in '%s' %v", name, a)
		}
		numField, ok := numFieldVal.Addr().Interface().(*Integer)
		if !ok {
			return fmt.Errorf("field for num of items in '%s' is not an Int", name)
		}

		// Update value with size of array
		numField.Value = size
	}

	// Loop through fields
	for i := 1; i < sVal.NumField(); i++ {

		// Get field value
		fieldVal := sVal.Field(i)

		// Get the field base and skip if line number is zero
		fb := fieldVal.FieldByName("FieldBase").Addr().Interface().(*FieldBase)
		if fb.Line == 0 {
			continue
		}

		// Switch based on field type
		switch v := fieldVal.Interface().(type) {
		case Path:
			pathPrefix := prefix
			if v.FileType == "Misc" {
				pathPrefix = ""
			}
			lines[v.Line-1] = fmt.Sprintf(`%11v   %-15s - %s`, `"`+pathPrefix+v.Value+`"`, v.Name, v.Desc)
		case Paths:
			pathPrefix := prefix
			if v.FileType == "Misc" {
				pathPrefix = ""
			}
			if v.Condensed {
				values := ""
				for _, value := range v.Value {
					values += `"` + pathPrefix + value + `" `
				}
				if len(v.Value) == 0 {
					values = `"unused"`
				}
				lines[v.Line-1] = fmt.Sprintf(`%11v   %-15s - %s`, values, v.Name, v.Desc)
			} else {
				lines[v.Line-1] = fmt.Sprintf(`%11v   %-15s - %s`, `"`+pathPrefix+v.Value[0]+`"`, v.Name, v.Desc)
				for j, value := range v.Value[1:] {
					lines[v.Line+j] = `"` + pathPrefix + value + `"`
				}
			}
		case Bool:
			lines[v.Line-1] = fmt.Sprintf(`%11v   %-15s - %s`, v.Value, v.Name, v.Desc)
		case String:
			lines[v.Line-1] = fmt.Sprintf(`%11v   %-15s - %s`, v.Value, v.Name, v.Desc)
		case Integer:
			lines[v.Line-1] = fmt.Sprintf(`%11v   %-15s - %s`, v.Value, v.Name, v.Desc)
		case Real:
			lines[v.Line-1] = fmt.Sprintf(`%11v   %-15s - %s`, v.Value, v.Name, v.Desc)
		case Reals:
			values := fmt.Sprint(v.Value)
			lines[v.Line-1] = fmt.Sprintf(`%11v   %-15s - %s`, values[1:len(values)-1], v.Name, v.Desc)
		}
	}

	return os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0777)
}
