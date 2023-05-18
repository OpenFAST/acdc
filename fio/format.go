package fio

import (
	"acdc/fio/schema"
	"bytes"
	"fmt"
	"os"
	"reflect"
)

type FileFormatter interface {
	FormatFile(path string) error
}

func format(s any, path string, schemas map[string]*schema.Schema) error {

	w := &bytes.Buffer{}

	sVal := reflect.ValueOf(s).Elem()
	sType := sVal.Type()

	// Loop through fields in struct
	for i := 0; i < sVal.NumField(); i++ {

		f := sVal.Field(i)
		sch := schemas[sType.Field(i).Name]

		// If field is File field
		if f, ok := f.Addr().Interface().(*File); ok {
			if err := f.FormatFile(path, sch.Field); err != nil {
				return err
			}
		}

		// If field is Files field
		if fs, ok := f.Addr().Interface().(*Files); ok {
			if err := fs.FormatFiles(path, sch.Field); err != nil {
				return err
			}
		}

		// Get field as formatter and write to file
		if p, ok := f.Interface().(Formatter); ok {
			p.Format(sch, w)
		} else {
			return fmt.Errorf("field %s does not support Formatter interface", sType.Field(i).Name)
		}
	}

	return os.WriteFile(path, w.Bytes(), 0777)
}
