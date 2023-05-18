package schema

import (
	"bytes"
	"go/format"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

const (
	BDStations   = "BDStations"
	Bool         = "Bool"
	BoolDefault  = "BoolDefault"
	Float        = "Float"
	FloatDefault = "FloatDefault"
	Floats       = "Floats"
	Header       = "Header"
	Int          = "Int"
	IntDefault   = "IntDefault"
	Ints         = "Ints"
	OutList      = "OutList"
	OutList2     = "OutList2"
	Path         = "Path"
	Paths        = "Paths"
	String       = "String"
	Table        = "Table"
	Text         = "Text"
	Title        = "Title"
	File         = "File"
	Files        = "Files"
)

type Schema struct {
	Order  int    `json:"Order"`
	Name   string `json:"Name"`
	Field  string `json:"Field"`
	Type   string `json:"Type"`
	Desc   string `json:"Desc"`
	Unit   string `json:"Unit"`
	NumVar string `json:"NumVar,omitempty"`
	NumInt int    `json:"NumInt,omitempty"`

	SwitchOptions   []Option `json:"SwitchOptions,omitempty"`
	TableHeaderSize int      `json:"TableHeaderSize,omitempty"`
	TableColumns    []Column `json:"TableColumns,omitempty"`
	FileType        string   `json:"FileType,omitempty"`
}

type Column struct {
	Name  string `json:"Name"`
	Field string `json:"Field"`
	Type  string `json:"Type"`
	Unit  string `json:"Unit"`
}

type Option struct {
	Value any    `json:"Value"`
	Label string `json:"Label"`
}

var genMap = map[string][]Schema{}
var Map = map[string]map[string]*Schema{}

func RegisterSchemas(name string, schemas []Schema) map[string]*Schema {
	genMap[name] = schemas
	headerNum := 0
	schemaMap := map[string]*Schema{}
	for i := range schemas {
		s := &schemas[i]
		s.Order = i + 1
		if s.Type == Header {
			headerNum++
			s.Name = "Header" + strconv.Itoa(headerNum)
		}
		s.Name = strings.TrimSpace(s.Name)
		s.Field = VarToField(s.Name)
		if s.Unit == "" {
			s.Unit = "-"
		}
		for j := range s.TableColumns {
			s.TableColumns[j].Field = VarToField(s.TableColumns[j].Name)
		}
		schemaMap[s.Field] = s
	}
	Map[name] = schemaMap
	return schemaMap
}

var reParen = regexp.MustCompile(`\(|\)`)

func VarToField(v string) string {
	field := reParen.ReplaceAllString(v, "")
	return strings.ToUpper(field[:1]) + field[1:]
}

func GenerateStructs(path string) error {

	w := &bytes.Buffer{}

	names := []string{}
	for name := range genMap {
		names = append(names, name)
	}
	sort.Strings(names)

	data := []any{}
	for _, name := range names {
		data = append(data, map[string]any{"Name": name, "Schemas": genMap[name]})
	}

	structsTemplate.Execute(w, data)

	contents, err := format.Source(w.Bytes())
	if err != nil {
		contents = w.Bytes()
	}

	return os.WriteFile(path, contents, 0777)
}

var structsTemplate = template.Must(template.New("files").Parse(`package fio

import (
	"acdc/fio/schema"
	"fmt"
	"bytes"
)

type File struct {
	Path
	{{- range $s := .}}
	{{$s.Name}} *{{$s.Name}} 
	{{- end}}
}

type Files struct {
	Paths
	{{- range $s := .}}
	{{$s.Name}} []*{{$s.Name}}
	{{- end}}
}

{{- range $s := .}}{{$name := $s.Name}}
//------------------------------------------------------------------------------
// {{$s.Name}}
//------------------------------------------------------------------------------

type {{$s.Name}} struct{
{{- range $schema := $s.Schemas}}
	{{- if ne $schema.Type "Table"}}
		{{$schema.Field}} {{$schema.Type}} // {{$schema.Desc}}
	{{- else}}
		{{$schema.Field}} Table{{$schema.Field}} // {{$schema.Desc}}
	{{- end}}
{{- end}}
}

func (s *{{$s.Name}}) Parse(path string) error {
	return parse(s, path, schema.{{$s.Name}})
}

func (s *{{$s.Name}}) Format(path string) error {
	return format(s, path, schema.{{$s.Name}})
}

{{- range $schema := $s.Schemas}}{{if ne $schema.Type "Table"}}{{continue}}{{end}}
type Table{{$schema.Field}} struct {
	Rows []Table{{$schema.Field}}Row
}

type Table{{$schema.Field}}Row struct {
{{- range $col := $schema.TableColumns}}
	{{$col.Field}} {{if eq $col.Type "Float" -}}
		float64
		{{- else if eq $col.Type "Int" -}}
		int
		{{- else if eq $col.Type "String" -}} 
		string
		{{- end}} // {{$col.Unit}}
	{{- end}}
}

func (t *Table{{$schema.Field}}) Parse(s *schema.Schema, lines []string, num int) ([]string, error) {
	if n:=num+s.TableHeaderSize; len(lines) < n {
		return nil, fmt.Errorf("insufficient lines, need %d", n)
	}
	lines = lines[s.TableHeaderSize:]
	t.Rows = make([]Table{{$schema.Field}}Row, num)
	for i, line := range lines[:num] {
		if err := lineToStruct(&t.Rows[i], line); err != nil{
			return nil, err
		}
	}
	return lines[num:], nil
}

func (t Table{{$schema.Field}}) Format(s *schema.Schema, w *bytes.Buffer) {
	if s.TableHeaderSize > 0 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", c.Name)
		}
		w.WriteString("\n")
	}
	if s.TableHeaderSize > 1 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", "("+c.Unit+")")
		}
		w.WriteString("\n")
	}
	for _, row := range t.Rows {
		structToLine(row, w)
	}
}

{{- end}}

{{- end}}
`))
