package fio

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func FormatAll(files FileMap, rootPath string) error {

	// Get directory
	dir := filepath.Dir(rootPath)

	// Get base name without suffix
	base := filepath.Base(rootPath)
	base = strings.TrimSuffix(base, filepath.Ext(base))

	// Create directory
	if err := os.MkdirAll(dir, 0777); err != nil {
		return fmt.Errorf("error creating directory '%s': %w", dir, err)
	}

	// Get main file
	main, ok := files["Main"]
	if !ok {
		return fmt.Errorf("main file not in list of files")
	}

	// Format all files
	if _, err := formatFiles(files, main, dir, base, ""); err != nil {
		return fmt.Errorf("error writing files: %w", err)
	}

	return nil
}

func formatFiles(files FileMap, file *File, dir, base, suffix string) (string, error) {

	// Create name for this file
	name := base + "_" + file.Name + suffix + ".dat"
	if file.Name == "Main" {
		name = base + ".fst"
	}

	// Create path for writing this file
	path := filepath.Join(dir, name)

	// Loop through fields in file
	for i := range file.Fields {

		// Get pointer to field
		field := &file.Fields[i]

		// If field is the title, update
		if field.Type == Title {
			field.Value = fmt.Sprintf("%s (Generated at %v)",
				path, time.Now().Format(time.RFC3339))
		}

		// If field is a path and path type is not empty
		if field.Type == Path && field.PathFileType != "" {

			// Get subfile from files, continue if it doesn't exist
			subFile, ok := files[file.Name+"."+field.Name]
			if !ok {
				field.Value = "unused"
				continue
			}

			suffix := ""
			if i := strings.Index(field.Name, "("); i != -1 {
				suffix = "_" + field.Name[i+1:strings.LastIndex(field.Name, ")")]
			}

			// Write file and get path
			subPath, err := formatFiles(files, subFile, dir, base, suffix)
			if err != nil {
				return "", fmt.Errorf("error writing file '%s': %w", subPath, err)
			}

			// Update path in field
			field.Value = subPath
		}

		// If field is multiple paths and path type is not empty
		if field.Type == Paths && field.PathFileType != "" {

			// Loop through field values
			for i := range field.Values {

				suffix := fmt.Sprintf("_%d", i+1)

				// Get subfile from files, continue if it doesn't exist
				subFile, ok := files[file.Name+"."+field.Name+suffix]
				if !ok {
					field.Value = "unused"
					continue
				}

				// Write file and get path
				subPath, err := formatFiles(files, subFile, dir, base, suffix)
				if err != nil {
					return "", fmt.Errorf("error writing file '%s': %w", subPath, err)
				}

				// Update path in field
				field.Values[i] = subPath
			}
		}
	}

	// Create buffer to write this file
	w := &bytes.Buffer{}
	if err := file.Format(w); err != nil {
		return "", fmt.Errorf("error formatting file '%s': %w", path, err)
	}

	// Write file to disk
	if err := os.WriteFile(path, w.Bytes(), 0777); err != nil {
		return "", fmt.Errorf("error writing file '%s': %w", path, err)
	}

	return name, nil
}

//------------------------------------------------------------------------------
// Formatting Functions
//------------------------------------------------------------------------------

type FormatFunc func(f Field, w io.Writer) error

var formatFuncs = map[Type]FormatFunc{
	Title: func(f Field, w io.Writer) error {
		_, err := fmt.Fprintf(w, "%s\n", f.Value)
		return err
	},
	Heading: func(f Field, w io.Writer) error {
		suffix := ""
		if n := 132 - len(f.Desc) - 8; n > 0 {
			suffix = strings.Repeat("-", n)
		}
		_, err := fmt.Fprintf(w, "------ %s %s\n", f.Desc, suffix)
		return err
	},
	Path: func(f Field, w io.Writer) error {
		_, err := fmt.Fprintf(w, "%12s    %-14s - %s (%s)\n", `"`+f.Value.(string)+`"`, f.Name, f.Desc, f.Unit)
		return err
	},
	Paths: func(f Field, w io.Writer) error {
		fmt.Fprintf(w, "%12s    %-14s - %s\n", `"`+f.Values[0].(string)+`"`, f.Name, f.Desc)
		for _, p := range f.Values[1:] {
			fmt.Fprintf(w, "\"%s\"\n", p)
		}
		return nil
	},
	String: func(f Field, w io.Writer) error {
		_, err := fmt.Fprintf(w, "%12s    %-14s - %s (%s)\n", f.Value, f.Name, f.Desc, f.Unit)
		return err
	},
	OutList: func(f Field, w io.Writer) error {
		fmt.Fprintf(w, "%12s    %-14s - %s\n", "", f.Name, f.Desc)
		for _, v := range f.Values {
			fmt.Fprintf(w, "%-12s\n", `"`+v.(string)+`"`)
		}
		w.Write([]byte("END of input file (\"END\" must appear " +
			"in the first 3 columns of this last OutList line)\n"))
		w.Write(bytes.Repeat([]byte("-"), 80))
		w.Write([]byte("\n"))
		return nil
	},
	OutList2: func(f Field, w io.Writer) error {
		for _, v := range f.Values {
			fmt.Fprintf(w, "%-12s\n", `"`+v.(string)+`"`)
		}
		w.Write([]byte("END of input file (\"END\" must appear " +
			"in the first 3 columns of this last OutList2 line)\n"))
		w.Write(append(bytes.Repeat([]byte("-"), 80), []byte("\n")...))
		return nil
	},
	Float: func(f Field, w io.Writer) (err error) {
		_, err = fmt.Fprintf(w, "%12g    %-14s - %s (%s)\n", f.Value, f.Name, f.Desc, f.Unit)
		return err
	},
	FloatDefault: func(f Field, w io.Writer) (err error) {
		if f.Value == nil {
			_, err = fmt.Fprintf(w, "%12s    %-14s - %s (%s)\n", "default", f.Name, f.Desc, f.Unit)
		} else {
			_, err = fmt.Fprintf(w, "%12g    %-14s - %s (%s)\n", f.Value, f.Name, f.Desc, f.Unit)
		}
		return err
	},
	FloatAll: func(f Field, w io.Writer) (err error) {
		if f.Value == nil {
			_, err = fmt.Fprintf(w, "%12s    %-14s - %s (%s)\n", "ALL", f.Name, f.Desc, f.Unit)
		} else {
			_, err = fmt.Fprintf(w, "%12g    %-14s - %s (%s)\n", f.Value, f.Name, f.Desc, f.Unit)
		}
		return err
	},
	Floats: func(f Field, w io.Writer) error {
		fio := fmt.Sprintf("%v", f.Values)
		_, err := fmt.Fprintf(w, "%12s    %-14s - %s (%s)\n", fio[1:len(fio)-1], f.Name, f.Desc, f.Unit)
		return err
	},
	Int: func(f Field, w io.Writer) (err error) {
		_, err = fmt.Fprintf(w, "%12d    %-14s - %s (%s)\n", f.Value, f.Name, f.Desc, f.Unit)
		return err
	},
	IntDefault: func(f Field, w io.Writer) (err error) {
		if f.Value == nil {
			_, err = fmt.Fprintf(w, "%12s    %-14s - %s (%s)\n", "default", f.Name, f.Desc, f.Unit)
		} else {
			_, err = fmt.Fprintf(w, "%12d    %-14s - %s (%s)\n", f.Value, f.Name, f.Desc, f.Unit)
		}
		return err
	},
	Ints: func(f Field, w io.Writer) error {
		fio := fmt.Sprintf("%v", f.Values)
		_, err := fmt.Fprintf(w, "%12s    %-14s - %s (%s)\n", fio[1:len(fio)-1], f.Name, f.Desc, f.Unit)
		return err
	},
	Bool: func(f Field, w io.Writer) (err error) {
		_, err = fmt.Fprintf(w, "%12v    %-14s - %s (%s)\n", f.Value, f.Name, f.Desc, f.Unit)
		return err
	},
	BoolDefault: func(f Field, w io.Writer) (err error) {
		if f.Value == nil {
			_, err = fmt.Fprintf(w, "%12s    %-14s - %s (%s)\n", "default", f.Name, f.Desc, f.Unit)
		} else {
			_, err = fmt.Fprintf(w, "%12v    %-14s - %s (%s)\n", f.Value, f.Name, f.Desc, f.Unit)
		}
		return err
	},
	Table: func(f Field, w io.Writer) (err error) {
		if f.TableHeaderSize > 0 {
			for _, col := range f.TableColumns {
				_, err = fmt.Fprintf(w, " %14s", col.Name)
			}
			fmt.Fprint(w, "\n")
		}
		if f.TableHeaderSize > 1 {
			for _, col := range f.TableColumns {
				_, err = fmt.Fprintf(w, " %14s", "("+col.Unit+")")
			}
			fmt.Fprint(w, "\n")
		}
		for _, row := range f.Table {
			for _, v := range row {
				_, err = fmt.Fprintf(w, " %14v", v)
			}
			fmt.Fprint(w, "\n")
		}
		return err
	},
	BDStations: func(f Field, w io.Writer) (err error) {
		for i, row := range f.Table {
			if i > 0 {
				fmt.Fprint(w, "\n")
			}
			fmt.Fprintf(w, " %14g\n", row[0])
			for _, r := range row[1].([6][6]float64) {
				for _, v := range r {
					fmt.Fprintf(w, " %14g", v)
				}
				fmt.Fprint(w, "\n")
			}
			fmt.Fprint(w, "\n")
			for _, r := range row[2].([6][6]float64) {
				for _, v := range r {
					fmt.Fprintf(w, " %14g", v)
				}
				fmt.Fprint(w, "\n")
			}
		}
		return nil
	},
	Text: func(f Field, w io.Writer) error {
		_, err := fmt.Fprint(w, f.Value)
		return err
	},
}
