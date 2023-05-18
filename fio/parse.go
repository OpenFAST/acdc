package fio

import (
	"acdc/fio/schema"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
)

func (main *Main) PostParse() error {

	// Set file type for aero file
	switch main.CompAero.Value {
	case 0:
		main.AeroFile.FileType = ""
	case 1:
		main.AeroFile.FileType = "AeroDyn14"
	case 2:
		main.AeroFile.FileType = "AeroDyn15"
	default:
		return fmt.Errorf("unsupported value for CompAero")
	}

	// Set file type for sub file
	switch main.CompSub.Value {
	case 0:
		main.SubFile.FileType = ""
	case 1:
		main.SubFile.FileType = "SubDyn"
	default:
		return fmt.Errorf("unsupported value for CompSub")
	}

	return nil
}

func parse(s any, path string, schemas map[string]*schema.Schema) error {

	lines, err := readLines(path)
	if err != nil {
		return err
	}

	sVal := reflect.ValueOf(s).Elem()
	sType := sVal.Type()

	for i := 0; i < sVal.NumField(); i++ {

		field := sVal.Field(i)
		sch := schemas[sType.Field(i).Name]

		// Find number of values to read
		num := 1
		if sch.NumInt != 0 {
			num = sch.NumInt
		} else if sch.NumVar != "" {
			numField := sVal.FieldByName(schema.VarToField(sch.NumVar))
			if !numField.IsValid() {
				return fmt.Errorf("unknown num var '%s'", sch.NumVar)
			}
			intVal, ok := numField.Interface().(Int)
			if !ok {
				return fmt.Errorf("num var '%s' type must be Int", sch.NumVar)
			}
			num = intVal.Value
		}

		// Get field as parser
		p, ok := field.Addr().Interface().(Parser)
		if !ok {
			return fmt.Errorf("field %s does not support parser interface", sType.Field(i).Name)
		}

		// Parse field
		var err error
		lines, err = p.Parse(sch, lines, num)
		if err != nil {
			return fmt.Errorf("error parsing field '%s': %w", sType.Field(i).Name, err)
		}
	}

	// If struct has a post-parsing function, call it
	if pp, ok := s.(PostParser); ok {
		if err := pp.PostParse(); err != nil {
			return fmt.Errorf("error in PostParse: %w", err)
		}
	}

	// Loop through fields again and call fields that support file parsing
	for i := 0; i < sVal.NumField(); i++ {

		// If field is the File type
		if f, ok := sVal.Field(i).Addr().Interface().(*File); ok {
			if err := f.ParseFile(filepath.Dir(path)); err != nil {
				return err
			}
		}

		// If field is the Files type
		if fs, ok := sVal.Field(i).Addr().Interface().(*Files); ok {
			if err := fs.ParseFiles(filepath.Dir(path)); err != nil {
				return err
			}
		}
	}

	return nil
}

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
