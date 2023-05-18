package fio

import (
	"acdc/fio/schema"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

type Parser interface {
	Parse(*schema.Schema, []string, int) ([]string, error)
}

type PostParser interface {
	PostParse() error
}

type Formatter interface {
	Format(*schema.Schema, *bytes.Buffer)
}

type FileParser interface {
	ParseFile(string) error
}

//------------------------------------------------------------------------------
// Header
//------------------------------------------------------------------------------

type Header struct{}

func (t *Header) Parse(s *schema.Schema, lines []string, _ int) ([]string, error) {
	return lines[1:], nil
}

func (t Header) Format(s *schema.Schema, w *bytes.Buffer) {
	suffix := ""
	if n := 132 - len(s.Desc) - 8; n > 0 {
		suffix = strings.Repeat("-", n)
	}
	fmt.Fprintf(w, "------ %s %s\n", s.Desc, suffix)
}

//------------------------------------------------------------------------------
// Title
//------------------------------------------------------------------------------

type Title struct {
	Value string
}

func (t *Title) Parse(s *schema.Schema, lines []string, _ int) ([]string, error) {
	t.Value = lines[0]
	return lines[1:], nil
}

func (t Title) Format(s *schema.Schema, w *bytes.Buffer) {
	fmt.Fprintf(w, "%s\n", t.Value)
}

//------------------------------------------------------------------------------
// Bool
//------------------------------------------------------------------------------

type Bool struct {
	Value bool
}

func (t *Bool) Parse(s *schema.Schema, lines []string, num int) ([]string, error) {
	values, _, err := splitLine(lines[0], s.Name, num)
	if err != nil {
		return nil, err
	}
	t.Value, err = strconv.ParseBool(values[0])
	if err != nil {
		return nil, err
	}
	return lines[1:], nil
}

func (t Bool) Format(s *schema.Schema, w *bytes.Buffer) {
	fmt.Fprintf(w, "%12v    %-14s - %s (%s)\n", t.Value, s.Name, s.Desc, s.Unit)
}

//------------------------------------------------------------------------------
// Bool Optional
//------------------------------------------------------------------------------

type BoolOptional struct {
	Value *bool
}

type BoolDefault BoolOptional

func (t *BoolDefault) Parse(s *schema.Schema, lines []string, num int) ([]string, error) {
	values, _, err := splitLine(lines[0], s.Name, num)
	if err != nil {
		return nil, err
	}
	if !strings.Contains(strings.ToLower(values[0]), "default") {
		v, err := strconv.ParseBool(values[0])
		if err != nil {
			return nil, err
		}
		t.Value = &v
	}
	return lines[1:], nil
}

func (t BoolDefault) Format(s *schema.Schema, w *bytes.Buffer) {
	if t.Value == nil {
		fmt.Fprintf(w, "%12s    %-14s - %s (%s)\n", "default", s.Name, s.Desc, s.Unit)
	} else {
		fmt.Fprintf(w, "%12v    %-14s - %s (%s)\n", *t.Value, s.Name, s.Desc, s.Unit)
	}
}

//------------------------------------------------------------------------------
// Float
//------------------------------------------------------------------------------

type Float struct {
	Value float64
}

func (t *Float) Parse(s *schema.Schema, lines []string, num int) ([]string, error) {
	values, _, err := splitLine(lines[0], s.Name, num)
	if err != nil {
		return nil, err
	}
	t.Value, err = strconv.ParseFloat(values[0], 64)
	if err != nil {
		return nil, err
	}
	return lines[1:], nil
}

func (t Float) Format(s *schema.Schema, w *bytes.Buffer) {
	fmt.Fprintf(w, "%12g    %-14s - %s (%s)\n", t.Value, s.Name, s.Desc, s.Unit)
}

//------------------------------------------------------------------------------
// Float Optional
//------------------------------------------------------------------------------

type FloatOptional struct {
	Value *float64
}

type FloatDefault FloatOptional
type FloatAll FloatOptional

func (t *FloatDefault) Parse(s *schema.Schema, lines []string, num int) ([]string, error) {
	values, _, err := splitLine(lines[0], s.Name, num)
	if err != nil {
		return nil, err
	}
	if !strings.Contains(strings.ToLower(values[0]), "default") {
		v, err := strconv.ParseFloat(values[0], 64)
		if err != nil {
			return nil, err
		}
		t.Value = &v
	}
	return lines[1:], nil
}

func (t FloatDefault) Format(s *schema.Schema, w *bytes.Buffer) {
	if t.Value == nil {
		fmt.Fprintf(w, "%12s    %-14s - %s (%s)\n", "default", s.Name, s.Desc, s.Unit)
	} else {
		fmt.Fprintf(w, "%12g    %-14s - %s (%s)\n", *t.Value, s.Name, s.Desc, s.Unit)
	}
}

func (t *FloatAll) Parse(s *schema.Schema, lines []string, num int) ([]string, error) {
	values, _, err := splitLine(lines[0], s.Name, num)
	if err != nil {
		return nil, err
	}
	if !strings.Contains(strings.ToLower(values[0]), "all") {
		v, err := strconv.ParseFloat(values[0], 64)
		if err != nil {
			return nil, err
		}
		t.Value = &v
	}
	return lines[1:], nil
}

func (t FloatAll) Format(s *schema.Schema, w *bytes.Buffer) {
	if t.Value == nil {
		fmt.Fprintf(w, "%12s    %-14s - %s (%s)\n", "ALL", s.Name, s.Desc, s.Unit)
	} else {
		fmt.Fprintf(w, "%12g    %-14s - %s (%s)\n", *t.Value, s.Name, s.Desc, s.Unit)
	}
}

//------------------------------------------------------------------------------
// Floats
//------------------------------------------------------------------------------

type Floats struct {
	Values []float64
}

func (t *Floats) Parse(s *schema.Schema, lines []string, num int) ([]string, error) {
	values, _, err := splitLine(lines[0], s.Name, num)
	if err != nil {
		return nil, err
	}
	t.Values = make([]float64, num)
	for j, v := range values[:num] {
		t.Values[j], err = strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}
	}
	return lines[1:], nil
}

func (t Floats) Format(s *schema.Schema, w *bytes.Buffer) {
	fio := fmt.Sprintf("%v", t.Values)
	fmt.Fprintf(w, "%12s    %-14s - %s (%s)\n", fio[1:len(fio)-1], s.Name, s.Desc, s.Unit)
}

//------------------------------------------------------------------------------
// Int
//------------------------------------------------------------------------------

type Int struct {
	Value int
}

func (t *Int) Parse(s *schema.Schema, lines []string, num int) ([]string, error) {
	values, _, err := splitLine(lines[0], s.Name, num)
	if err != nil {
		return nil, err
	}
	t.Value, err = strconv.Atoi(values[0])
	if err != nil {
		return nil, err
	}
	return lines[1:], nil
}

func (t Int) Format(s *schema.Schema, w *bytes.Buffer) {
	fmt.Fprintf(w, "%12d    %-14s - %s (%s)\n", t.Value, s.Name, s.Desc, s.Unit)
}

//------------------------------------------------------------------------------
// Int Optional
//------------------------------------------------------------------------------

type IntOptional struct {
	Value *int
}

type IntDefault IntOptional

func (t *IntDefault) Parse(s *schema.Schema, lines []string, num int) ([]string, error) {
	values, _, err := splitLine(lines[0], s.Name, num)
	if err != nil {
		return nil, err
	}
	if !strings.Contains(strings.ToLower(values[0]), "default") {
		v, err := strconv.Atoi(values[0])
		if err != nil {
			return nil, err
		}
		t.Value = &v
	}
	return lines[1:], nil
}

func (t IntDefault) Format(s *schema.Schema, w *bytes.Buffer) {
	if t.Value == nil {
		fmt.Fprintf(w, "%12s    %-14s - %s (%s)\n", "default", s.Name, s.Desc, s.Unit)
	} else {
		fmt.Fprintf(w, "%12d    %-14s - %s (%s)\n", *t.Value, s.Name, s.Desc, s.Unit)
	}
}

//------------------------------------------------------------------------------
// Ints
//------------------------------------------------------------------------------

type Ints struct {
	Values []int
}

func (t *Ints) Parse(s *schema.Schema, lines []string, num int) ([]string, error) {
	if num < 1 {
		num = 1
	}
	values, _, err := splitLine(lines[0], s.Name, num)
	if err != nil {
		return nil, err
	}
	t.Values = make([]int, num)
	for j, v := range values[:num] {
		t.Values[j], err = strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
	}
	return lines[1:], nil
}

func (t Ints) Format(s *schema.Schema, w *bytes.Buffer) {
	fio := fmt.Sprintf("%v", t.Values)
	fmt.Fprintf(w, "%12s    %-14s - %s (%s)\n", fio[1:len(fio)-1], s.Name, s.Desc, s.Unit)
}

//------------------------------------------------------------------------------
// String
//------------------------------------------------------------------------------

type String struct {
	Value string
}

func (t *String) Parse(s *schema.Schema, lines []string, _ int) ([]string, error) {
	values, _, err := splitLine(lines[0], s.Name, 1)
	if err != nil {
		return nil, err
	}
	t.Value = values[0]
	return lines[1:], nil
}

func (t String) Format(s *schema.Schema, w *bytes.Buffer) {
	fmt.Fprintf(w, "%12s    %-14s - %s (%s)\n", t.Value, s.Name, s.Desc, s.Unit)
}

//------------------------------------------------------------------------------
// Strings
//------------------------------------------------------------------------------

type Strings struct {
	Values []string
}

//------------------------------------------------------------------------------
// Path
//------------------------------------------------------------------------------

type Path struct {
	Value    string
	FileType string
}

func (t *Path) Parse(s *schema.Schema, lines []string, _ int) ([]string, error) {
	t.FileType = s.FileType
	values, _, err := splitLine(lines[0], s.Name, 1)
	if err != nil {
		return nil, err
	}
	t.Value = strings.Trim(values[0], `"`)
	return lines[1:], nil
}

func (t Path) Format(s *schema.Schema, w *bytes.Buffer) {
	fmt.Fprintf(w, "%12s    %-14s - %s (%s)\n", `"`+t.Value+`"`, s.Name, s.Desc, s.Unit)
}

//------------------------------------------------------------------------------
// Paths
//------------------------------------------------------------------------------

type Paths struct {
	Values   []string
	FileType string
}

func (t *Paths) Parse(s *schema.Schema, lines []string, size int) ([]string, error) {
	t.FileType = s.FileType
	values, _, err := splitLine(lines[0], s.Name, 1)
	if err != nil {
		return nil, err
	}
	for _, v := range values {
		t.Values = append(t.Values, strings.Trim(v, `"`))
	}
	for len(t.Values) < size {
		lines = lines[1:]
		values, _, err := splitLine(lines[0], "", 1)
		if err != nil {
			return nil, err
		}
		for _, v := range values {
			t.Values = append(t.Values, strings.Trim(v, `"`))
		}
	}
	return lines[1:], nil
}

func (t Paths) Format(s *schema.Schema, w *bytes.Buffer) {
	fmt.Fprintf(w, "%12s    %-14s - %s\n", `"`+t.Values[0]+`"`, s.Name, s.Desc)
	for _, p := range t.Values[1:] {
		fmt.Fprintf(w, "\"%s\"\n", p)
	}
}

//------------------------------------------------------------------------------
// Text
//------------------------------------------------------------------------------

type Text struct {
	Value string
}

func (t *Text) Parse(s *schema.Schema, lines []string, _ int) ([]string, error) {
	t.Value = strings.Join(lines, "\n")
	return []string{}, nil
}

func (t Text) Format(s *schema.Schema, w *bytes.Buffer) {
	fmt.Fprint(w, t.Value)
}

//------------------------------------------------------------------------------
// OutList
//------------------------------------------------------------------------------

type OutList Strings
type OutList2 Strings

func (t *OutList) Parse(s *schema.Schema, lines []string, _ int) ([]string, error) {
	line := ""
	line, lines = lines[0], lines[1:]
	_, _, err := splitLine(line, s.Name, 0)
	if err != nil {
		return nil, err
	}
	foundEnd := false
	for len(lines) > 0 {
		line, lines = lines[0], lines[1:]
		if strings.HasPrefix(strings.ToLower(line), "end") {
			foundEnd = true
			break
		}
		line, _, _ = strings.Cut(line, "-")
		if i := strings.LastIndex(line, `"`); i != -1 {
			line = line[:i]
		}
		vars := strings.FieldsFunc(line, func(r rune) bool {
			return r == '"' || r == ',' || r == ' ' || r == '\t'
		})
		t.Values = append(t.Values, vars...)
	}
	if !foundEnd {
		return nil, fmt.Errorf("END not found")
	}
	return lines, nil
}

func (t OutList) Format(s *schema.Schema, w *bytes.Buffer) {
	fmt.Fprintf(w, "%12s    %-14s - %s\n", "", s.Name, s.Desc)
	for _, v := range t.Values {
		fmt.Fprintf(w, "%-12s\n", `"`+v+`"`)
	}
	w.Write([]byte("END of input file (\"END\" must appear " +
		"in the first 3 columns of this last OutList line)\n"))
	w.Write(bytes.Repeat([]byte("-"), 80))
	w.Write([]byte("\n"))
}

func (t *OutList2) Parse(s *schema.Schema, lines []string, _ int) ([]string, error) {
	foundEnd := false
	for len(lines) > 0 {
		if strings.HasPrefix(strings.ToLower(lines[0]), "end") {
			foundEnd = true
			break
		}
		line, _, _ := strings.Cut(lines[0], "-")
		if i := strings.LastIndex(line, `"`); i != -1 {
			line = line[:i]
		}
		values := strings.FieldsFunc(line, func(r rune) bool {
			return r == '"' || r == ',' || r == ' ' || r == '\t'
		})
		t.Values = append(t.Values, values...)
		lines = lines[1:]
	}
	if !foundEnd {
		return nil, fmt.Errorf("END not found")
	}
	return lines, nil
}

func (t OutList2) Format(s *schema.Schema, w *bytes.Buffer) {
	for _, v := range t.Values {
		fmt.Fprintf(w, "%-12s\n", `"`+v+`"`)
	}
	w.Write([]byte("END of input file (\"END\" must appear " +
		"in the first 3 columns of this last OutList2 line)\n"))
	w.Write(append(bytes.Repeat([]byte("-"), 80), []byte("\n")...))
}

//------------------------------------------------------------------------------
// BDStations
//------------------------------------------------------------------------------

type BDStations struct {
	Rows []BDStationsRow
}

type BDStationsRow struct {
	Eta float64
	K   [6][6]float64
	M   [6][6]float64
}

func (t *BDStations) Parse(s *schema.Schema, lines []string, num int) ([]string, error) {
	if n := num*15 - 1; len(lines) < n {
		return nil, fmt.Errorf("insufficient lines, need %d", n)
	}
	t.Rows = make([]BDStationsRow, num)
	for i := range t.Rows {
		if i > 0 {
			lines = lines[1:]
		}
		var err error
		t.Rows[i].Eta, err = strconv.ParseFloat(strings.TrimSpace(lines[0]), 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing eta for station %d: %w", i+1, err)
		}
		lines = lines[1:]
		for j := range t.Rows[i].K {
			for k, v := range strings.Fields(lines[0]) {
				t.Rows[i].K[j][k], err = strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("error parsing stiffness for station %d: %w", i+1, err)
				}
			}
			lines = lines[1:]
		}
		lines = lines[1:]
		for j := range t.Rows[i].M {
			for k, v := range strings.Fields(lines[0]) {
				t.Rows[i].M[j][k], err = strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("error parsing mass for station %d: %w", i+1, err)
				}
			}
			lines = lines[1:]
		}
	}
	return lines, nil
}

func (t BDStations) Format(s *schema.Schema, w *bytes.Buffer) {
	for i, row := range t.Rows {
		if i > 0 {
			fmt.Fprint(w, "\n")
		}
		fmt.Fprintf(w, " %14g\n", row.Eta)
		for _, r := range row.K {
			for _, v := range r {
				fmt.Fprintf(w, " %14g", v)
			}
			fmt.Fprint(w, "\n")
		}
		fmt.Fprint(w, "\n")
		for _, r := range row.M {
			for _, v := range r {
				fmt.Fprintf(w, " %14g", v)
			}
			fmt.Fprint(w, "\n")
		}
	}
}

//------------------------------------------------------------------------------
// Files
//------------------------------------------------------------------------------

func checkPath(dir, path string) (string, bool) {
	path = filepath.Clean(filepath.Join(dir, path))
	if _, err := os.Stat(path); err == nil {
		return path, true
	}
	return path, false
}

func (f *File) ParseFile(dir string) error {

	// If no file type specified, return
	if f.FileType == "" {
		return nil
	}

	// Get schema from file type
	schemaMap, ok := schema.Map[f.FileType]
	if !ok {
		return fmt.Errorf("no schema to format file type '%s'", f.FileType)
	}

	// If path doesn't exist, clear path and continue
	path, ok := checkPath(dir, f.Path.Value)
	if !ok {
		f.Path.Value = ""
		return nil
	}

	// Remove directory from path
	f.Path.Value = filepath.Base(f.Path.Value)

	// Get file fieldVal by name, return error if not found
	fieldVal := reflect.ValueOf(f).Elem().FieldByName(f.FileType)

	// Create new type
	fieldVal.Set(reflect.New(fieldVal.Type().Elem()))

	// Parse file
	if err := parse(fieldVal.Interface(), path, schemaMap); err != nil {
		return fmt.Errorf("error parsing '%s': %w", path, err)
	}

	return nil
}

func (f *File) FormatFile(parentPath, fieldName string) error {

	if f.FileType == "" {
		return nil
	}

	// Get schema from file type
	schemaMap, ok := schema.Map[f.FileType]
	if !ok {
		return fmt.Errorf("no schema to format file type '%s'", f.FileType)
	}

	// Get file field by name, if nil, return
	field := reflect.ValueOf(f).Elem().FieldByName(f.FileType)
	if field.IsNil() {
		return nil
	}

	// Update path to file
	parentName := strings.TrimSuffix(filepath.Base(parentPath), filepath.Ext(parentPath))
	suffix := strings.ReplaceAll(fieldName, "Filename_", "")
	suffix = strings.ReplaceAll(suffix, "File", "")
	suffix = strings.ReplaceAll(suffix, "Name", "")
	f.Path.Value = parentName + "_" + suffix + ".dat"

	// Create path to file
	path := filepath.Join(filepath.Dir(parentPath), f.Path.Value)

	// Format file
	err := format(field.Interface(), path, schemaMap)
	if err != nil {
		return fmt.Errorf("error formatting '%s': %w", path, err)
	}

	return nil
}

func (f *Files) ParseFiles(dir string) error {

	// If no file type specified, return
	if f.FileType == "" {
		return nil
	}

	// Get schema from file type
	schemaMap, ok := schema.Map[f.FileType]
	if !ok {
		return fmt.Errorf("no schema to format file type '%s'", f.FileType)
	}

	// Get file fieldVal by name, return error if not found
	fieldVal := reflect.ValueOf(f).Elem().FieldByName(f.FileType)

	// Get structure type
	structType := fieldVal.Type().Elem().Elem()

	// Set field to slice of pointers to structure type
	fieldVal.Set(reflect.MakeSlice(reflect.SliceOf(reflect.PtrTo(structType)),
		len(f.Paths.Values), len(f.Paths.Values)))

	// Loop through paths
	for i := range f.Paths.Values {

		// If path doesn't exist, clear path and continue
		path, ok := checkPath(dir, f.Paths.Values[i])
		if !ok {
			f.Paths.Values[i] = ""
			return nil
		}

		// Remove directory from path
		f.Paths.Values[i] = filepath.Base(f.Paths.Values[i])

		// Create new structure
		fileStruct := reflect.New(structType)

		// Add file struct to field
		fieldVal.Index(i).Set(fileStruct)

		// Parse file
		if err := parse(fileStruct.Interface(), path, schemaMap); err != nil {
			return fmt.Errorf("error parsing '%s': %w", path, err)
		}
	}

	return nil
}

func (f *Files) FormatFiles(parentPath, fieldName string) error {

	if f.FileType == "" {
		return nil
	}

	// Get schema from file type
	schemaMap, ok := schema.Map[f.FileType]
	if !ok {
		return fmt.Errorf("no schema to format file type '%s'", f.FileType)
	}

	// Get file field by name, return error if not found
	field := reflect.ValueOf(f).Elem().FieldByName(f.FileType)

	// Create parts of file name
	parentName := strings.TrimSuffix(filepath.Base(parentPath), filepath.Ext(parentPath))
	suffix := strings.ReplaceAll(fieldName, "File", "")
	suffix = strings.ReplaceAll(suffix, "Names", "")

	// Loop through paths
	for i := range f.Paths.Values {

		// Update file name
		f.Paths.Values[i] = parentName + "_" + suffix + strconv.Itoa(i+1) + ".dat"

		// Create path to file
		path := filepath.Join(filepath.Dir(parentPath), f.Paths.Values[i])

		// Get pointer to file structure in slice, return if nil
		fileStruct := field.Index(i)
		if fileStruct.IsNil() {
			return nil
		}

		// Format file
		err := format(fileStruct.Interface(), path, schemaMap)
		if err != nil {
			return fmt.Errorf("error formatting '%s': %w", path, err)
		}
	}

	return nil
}
