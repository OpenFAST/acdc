package fio

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ParseAll(path string) (FileMap, error) {

	main, err := parse("Main", path)
	if err != nil {
		return nil, fmt.Errorf("error parsing file '%s': %w", path, err)
	}

	// Create map
	files := FileMap{"Main": main}

	// Set file type for aero file
	switch main.Field("CompAero").Value.(int) {
	case 0:
		main.Field("AeroFile").PathFileType = ""
	case 1:
		main.Field("AeroFile").PathFileType = "AeroDyn14"
	case 2:
		main.Field("AeroFile").PathFileType = "AeroDyn15"
	default:
		return nil, fmt.Errorf("unsupported value for CompAero")
	}

	// Set file type for servo file
	switch main.Field("CompServo").Value.(int) {
	case 0:
		main.Field("ServoFile").PathFileType = ""
	case 1:
		main.Field("ServoFile").PathFileType = "ServoDyn"
	default:
		return nil, fmt.Errorf("unsupported value for CompServo")
	}

	// Set file type for sub file
	switch main.Field("CompSub").Value.(int) {
	case 0:
		main.Field("SubFile").PathFileType = ""
	case 1:
		main.Field("SubFile").PathFileType = "SubDyn"
	default:
		return nil, fmt.Errorf("unsupported value for CompSub")
	}

	// Parse any paths in file
	if err = parsePaths(main, filepath.Dir(path), files); err != nil {
		return nil, fmt.Errorf("error parsing paths: %w", err)
	}

	return files, nil
}

func parsePaths(f *File, dir string, files FileMap) error {

	// Loop through fields in file
	for _, field := range f.Fields {

		// If field is a path and path type is not empty,
		if field.Type == Path && field.PathFileType != "" {

			// Build path and see if file exists. if not exist, continue
			path := filepath.Clean(filepath.Join(dir, field.Value.(string)))
			if _, err := os.Stat(path); err != nil {
				field.Value = "unused"
				continue
			}

			// Parse file
			file, err := parse(field.PathFileType, path)
			if err != nil {
				return fmt.Errorf("error parsing file '%s' for field '%s' in file '%s'",
					path, field.Name, f.Name)
			}

			// Add file to map
			files[f.Name+"."+field.Name] = file

			// Parse any paths in file
			if err = parsePaths(file, filepath.Dir(path), files); err != nil {
				return err
			}
		}

		// If field is multiple paths and path type is not empty,
		if field.Type == Paths && field.PathFileType != "" {

			for i, v := range field.Values {

				// Build path and see if file exists. if not exist, continue
				path := filepath.Clean(filepath.Join(dir, v.(string)))
				if _, err := os.Stat(path); err != nil {
					field.Values[i] = "unused"
					continue
				}

				// Parse file
				file, err := parse(field.PathFileType, path)
				if err != nil {
					return fmt.Errorf("error parsing file '%s' for field '%s' in file '%s'",
						path, field.Name, f.Name)
				}

				// Add file to map
				files[f.Name+"."+field.Name+"_"+strconv.Itoa(i+1)] = file

				// Parse any paths in file
				if err = parsePaths(file, filepath.Dir(path), files); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func parse(fileType string, path string) (*File, error) {

	// Get file type from map by name
	ft, ok := fileTypes[fileType]
	if !ok {
		return nil, fmt.Errorf("unknown file type '%s'", fileType)
	}

	// Get function for reading lines from file
	readFile, ok := readFuncs[ft.Name]
	if !ok {
		readFile = readLines
	}

	// Read file
	lines, err := readFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file '%s': %w", path, err)
	}

	// Create a copy of the file type to populate with values
	bs, err := json.Marshal(ft)
	if err != nil {
		return nil, err
	}
	f := &File{}
	if err = json.Unmarshal(bs, f); err != nil {
		return nil, err
	}

	// Create map relating field names to fields
	fieldMap := map[string]*Field{}
	for i := range f.Fields {
		fieldMap[f.Fields[i].Name] = &f.Fields[i]
	}

	// Loop through fields and parse file
	for i := range f.Fields {

		field := &f.Fields[i]

		// TODO: Check version

		if len(lines) == 0 {
			return nil, fmt.Errorf("no lines to parse")
		}

		// Get number of values to read
		num := 1
		if field.Num != "" {
			sizeField, ok := fieldMap[field.Num]
			if ok {
				num, ok = sizeField.Value.(int)
				if !ok {
					return nil, fmt.Errorf("field '%s': size field '%s' is not an Int type", field.Name, sizeField.Name)
				}
			} else if num, err = strconv.Atoi(field.Num); err != nil {
				return nil, fmt.Errorf("field '%s': invalid size '%s'", field.Name, field.Num)
			}
		}

		// Get parse func based on type
		pf, ok := parseFuncs[field.Type]
		if !ok {
			return nil, fmt.Errorf("field '%s': no parse func for type '%s'", field.Name, field.Type)
		}

		// Call parse function
		lines, err = pf(field, lines, num)
		if err != nil {
			return nil, fmt.Errorf("field '%s': parse error: %w", field.Name, err)
		}
	}

	return f, nil
}

//------------------------------------------------------------------------------
// Reading Functions
//------------------------------------------------------------------------------

type ReadFunc func(string) ([]string, error)

var readFuncs = map[string]ReadFunc{
	// "AirfoilInfo": func(path string) ([]string, error) {
	// 	lines, err := readLines(path)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	tmp := lines
	// 	lines = lines[:0]
	// 	for _, line := range tmp {
	// 		if !strings.HasPrefix(line, "!") {
	// 			lines = append(lines, line)
	// 		}
	// 	}
	// 	return lines, nil
	// },
}

//------------------------------------------------------------------------------
// Parsing Functions
//------------------------------------------------------------------------------

type ParseFunc func(f *Field, lines []string, num int) ([]string, error)

var parseFuncs = map[Type]ParseFunc{
	Title: func(f *Field, lines []string, _ int) ([]string, error) {
		if len(lines) < 1 {
			return nil, fmt.Errorf("insufficient lines")
		}
		f.Value = lines[0]
		return lines[1:], nil
	},
	Heading: func(f *Field, lines []string, _ int) ([]string, error) {
		if len(lines) < 1 {
			return nil, fmt.Errorf("insufficient lines")
		}
		return lines[1:], nil
	},
	Path: func(f *Field, lines []string, _ int) ([]string, error) {
		values, _, err := splitLine(lines[0], f.Name, 1)
		if err != nil {
			return nil, err
		}
		f.Value = strings.Trim(values[0], `"`)
		return lines[1:], nil
	},
	Paths: func(f *Field, lines []string, size int) ([]string, error) {
		values, _, err := splitLine(lines[0], f.Name, 1)
		if err != nil {
			return nil, err
		}
		for _, v := range values {
			f.Values = append(f.Values, strings.Trim(v, `"`))
		}
		for len(f.Values) < size {
			lines = lines[1:]
			values, _, err := splitLine(lines[0], "", 1)
			if err != nil {
				return nil, err
			}
			for _, v := range values {
				f.Values = append(f.Values, strings.Trim(v, `"`))
			}
		}
		return lines[1:], nil
	},
	String: func(f *Field, lines []string, _ int) ([]string, error) {
		values, _, err := splitLine(lines[0], f.Name, 1)
		if err != nil {
			return nil, err
		}
		f.Value = values[0]
		return lines[1:], nil
	},
	Text: func(f *Field, lines []string, _ int) ([]string, error) {
		f.Value = strings.Join(lines, "\n")
		return []string{}, nil
	},
	OutList: func(f *Field, lines []string, _ int) ([]string, error) {
		line := ""
		line, lines = lines[0], lines[1:]
		_, _, err := splitLine(line, f.Name, 0)
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
			for _, v := range vars {
				f.Values = append(f.Values, v)
			}
		}
		if !foundEnd {
			return nil, fmt.Errorf("END not found")
		}
		return lines, nil
	},
	OutList2: func(f *Field, lines []string, _ int) ([]string, error) {
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
			for _, v := range values {
				f.Values = append(f.Values, v)
			}
			lines = lines[1:]
		}
		if !foundEnd {
			return nil, fmt.Errorf("END not found")
		}
		return lines, nil
	},
	Float: func(f *Field, lines []string, num int) ([]string, error) {
		values, _, err := splitLine(lines[0], f.Name, num)
		if err != nil {
			return nil, err
		}
		f.Value, err = strconv.ParseFloat(values[0], 64)
		if err != nil {
			return nil, err
		}
		return lines[1:], nil
	},
	FloatDefault: func(f *Field, lines []string, num int) ([]string, error) {
		values, _, err := splitLine(lines[0], f.Name, num)
		if err != nil {
			return nil, err
		}
		if !strings.Contains(strings.ToLower(values[0]), "default") {
			v, err := strconv.ParseFloat(values[0], 64)
			if err != nil {
				return nil, err
			}
			f.Value = v
		}
		return lines[1:], nil
	},
	FloatAll: func(f *Field, lines []string, num int) ([]string, error) {
		values, _, err := splitLine(lines[0], f.Name, num)
		if err != nil {
			return nil, err
		}
		if !strings.Contains(strings.ToLower(values[0]), "all") {
			v, err := strconv.ParseFloat(values[0], 64)
			if err != nil {
				return nil, err
			}
			f.Value = v
		}
		return lines[1:], nil
	},

	Floats: func(f *Field, lines []string, num int) ([]string, error) {
		if num < 1 {
			num = 1
		}
		values, _, err := splitLine(lines[0], f.Name, num)
		if err != nil {
			return nil, err
		}
		f.Values = make([]any, num)
		for j, v := range values[:num] {
			f.Values[j], err = strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, err
			}
		}
		return lines[1:], nil
	},
	Int: func(f *Field, lines []string, num int) ([]string, error) {
		values, _, err := splitLine(lines[0], f.Name, num)
		if err != nil {
			return nil, err
		}
		f.Value, err = strconv.Atoi(values[0])
		if err != nil {
			return nil, err
		}
		return lines[1:], nil
	},
	IntDefault: func(f *Field, lines []string, num int) ([]string, error) {
		values, _, err := splitLine(lines[0], f.Name, num)
		if err != nil {
			return nil, err
		}
		if !strings.Contains(strings.ToLower(values[0]), "default") {
			v, err := strconv.Atoi(values[0])
			if err != nil {
				return nil, err
			}
			f.Value = v
		}
		return lines[1:], nil
	},
	Ints: func(f *Field, lines []string, num int) ([]string, error) {
		if num < 1 {
			num = 1
		}
		values, _, err := splitLine(lines[0], f.Name, num)
		if err != nil {
			return nil, err
		}
		f.Values = make([]any, num)
		for j, v := range values[:num] {
			f.Values[j], err = strconv.Atoi(v)
			if err != nil {
				return nil, err
			}
		}
		return lines[1:], nil
	},
	Bool: func(f *Field, lines []string, num int) ([]string, error) {
		values, _, err := splitLine(lines[0], f.Name, num)
		if err != nil {
			return nil, err
		}
		f.Value, err = strconv.ParseBool(values[0])
		if err != nil {
			return nil, err
		}
		return lines[1:], nil
	},
	BoolDefault: func(f *Field, lines []string, num int) ([]string, error) {
		values, _, err := splitLine(lines[0], f.Name, num)
		if err != nil {
			return nil, err
		}
		if !strings.Contains(strings.ToLower(values[0]), "default") {
			v, err := strconv.ParseBool(values[0])
			if err != nil {
				return nil, err
			}
			f.Value = v
		}
		return lines[1:], nil
	},
	Table: func(f *Field, lines []string, num int) ([]string, error) {
		if len(lines) < num+f.TableHeaderSize {
			return nil, fmt.Errorf("insufficient lines, need %d", num+f.TableHeaderSize)
		}
		lines = lines[f.TableHeaderSize:]
		f.Table = make([][]any, num)
		for i := range f.Table {
			f.Table[i] = make([]any, len(f.TableColumns))
			for j, col := range f.TableColumns {
				switch col.Type {
				case Int:
					f.Table[i][j] = 0
				case Float:
					f.Table[i][j] = 0.0
				case String:
					f.Table[i][j] = ""
				}
			}
		}
		var err error
		for i, line := range lines[:num] {
			line, _, _ = strings.Cut(line, "- ")
			for j, s := range strings.Fields(line) {
				switch f.TableColumns[j].Type {
				case Int:
					f.Table[i][j], err = strconv.Atoi(s)
					if err != nil {
						return nil, fmt.Errorf("error parsing '%s' into field '%s'", s, f.TableColumns[j].Name)
					}
				case Float:
					f.Table[i][j], err = strconv.ParseFloat(s, 64)
					if err != nil {
						return nil, fmt.Errorf("error parsing '%s' into field '%s'", s, f.TableColumns[j].Name)
					}
				case String:
					f.Table[i][j] = s
				}
			}
		}
		return lines[num:], nil
	},
	BDStations: func(f *Field, lines []string, num int) ([]string, error) {
		if n := num*15 - 1; len(lines) < n {
			return nil, fmt.Errorf("insufficient lines, need %d", n)
		}
		f.Table = make([][]any, num)
		for i := range f.Table {
			if i > 0 {
				lines = lines[1:]
			}
			eta, err := strconv.ParseFloat(strings.TrimSpace(lines[0]), 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing eta for station %d: %w", i+1, err)
			}
			lines = lines[1:]
			stiff := [6][6]float64{}
			for j := range stiff {
				for k, v := range strings.Fields(lines[0]) {
					stiff[j][k], err = strconv.ParseFloat(v, 64)
					if err != nil {
						return nil, fmt.Errorf("error parsing stiffness for station %d: %w", i+1, err)
					}
				}
				lines = lines[1:]
			}
			lines = lines[1:]
			mass := [6][6]float64{}
			for j := range mass {
				for k, v := range strings.Fields(lines[0]) {
					mass[j][k], err = strconv.ParseFloat(v, 64)
					if err != nil {
						return nil, fmt.Errorf("error parsing mass for station %d: %w", i+1, err)
					}
				}
				lines = lines[1:]
			}
			f.Table[i] = []any{eta, stiff, mass}
		}
		return lines, nil
	},
}

//------------------------------------------------------------------------------
// Utility
//------------------------------------------------------------------------------

func splitLine(line, keyword string, num int) ([]string, string, error) {
	found := false
	comment := ""
	if line, comment, found = strings.Cut(line, "- "); !found {
		line, comment, _ = strings.Cut(line, "! ")
	}

	if keyword != "" {
		i := strings.Index(strings.ToLower(line), strings.ToLower(keyword))
		if i == -1 {
			return nil, "", fmt.Errorf("keyword '%s' not found", keyword)
		}
		line = line[:i] + line[i+len(keyword):]
	}
	values := strings.FieldsFunc(line, func(r rune) bool {
		return r == ' ' || r == ','
	})
	if n := len(values); n < num {
		return nil, "", fmt.Errorf("found %d values, requested %d", n, num)
	}
	return values, comment, nil
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
		line := strings.TrimSpace(scanner.Text())

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
