package fio

import (
	"fmt"
	"io"
)

type Type string

type FileMap map[string]*File

var fileTypes = FileMap{}

type File struct {
	Name   string  `json:"Name"`
	Fields []Field `json:"Fields"`
}

func NewFile(name string, fields []Field) *File {

	ft := &File{
		Name:   name,
		Fields: fields,
	}

	for i := range fields {

		f := &fields[i]

		if f.Unit == "" {
			f.Unit = "-"
		}
	}

	// Register file type in map
	fileTypes[name] = ft

	return ft
}

func (f File) Format(w io.Writer) error {

	for _, field := range f.Fields {

		// Get format func based on type
		ff, ok := formatFuncs[field.Type]
		if !ok {
			return fmt.Errorf("field '%s': no format func for type '%s'", field.Name, field.Type)
		}

		// Call format function
		if err := ff(field, w); err != nil {
			return fmt.Errorf("field '%s': format error: %w", field.Name, err)
		}
	}

	return nil
}

func (f *File) Field(name string) *Field {
	for i := range f.Fields {
		if f.Fields[i].Name == name {
			return &f.Fields[i]
		}
	}
	return nil
}

var AirfoilInfo = NewFile("AirfoilInfo", []Field{
	{Name: "Text", Type: Text},
})

var Foil = NewFile("Foil", []Field{
	{Name: "Text", Type: Text},
})

var TailFin = NewFile("TailFin", []Field{
	{Name: "Text", Type: Text},
})

var UniformWind = NewFile("UniformWind", []Field{
	{Name: "Text", Type: Text},
})
