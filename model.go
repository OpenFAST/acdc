package main

type Model struct {
	HasAero       bool     `json:"HasAero"`
	ImportedPaths []string `json:"ImportedPaths"`
	Files         *Files   `json:"Files"`
	Notes         []string `json:"Notes"`
}

func NewModel() *Model {
	return &Model{}
}
