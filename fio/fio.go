package fio

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

//go:generate go run schema/gen.go

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

func lineToStruct(s any, line string) error {
	sVal := reflect.ValueOf(s).Elem()
	line, _, _ = strings.Cut(line, "- ")
	for i, f := range strings.Fields(line) {
		_, err := fmt.Sscan(f, sVal.Field(i).Addr().Interface())
		if err != nil {
			return fmt.Errorf("error parsing '%s' into field '%s': %w", f, sVal.Type().Field(i).Name, err)
		}
	}
	return nil
}

func structToLine(s any, w *bytes.Buffer) {
	sVal := reflect.ValueOf(s)
	for i := 0; i < sVal.NumField(); i++ {
		fmt.Fprintf(w, " %14v", sVal.Field(i).Interface())
	}
	w.WriteByte('\n')
}
