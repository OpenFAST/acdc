package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Project struct {
	Info    Info
	Exec    Exec
	Model   *Model
	Analyze *Analyze
}

type Info struct {
	Date string `json:"Date"`
	Path string `json:"Path"`
}

type Exec struct {
	Path    string `json:"Path"`
	Version string `json:"Version"`
	Valid   bool   `json:"Valid"`
}

func NewProject() *Project {
	return &Project{}
}

func LoadProject(path string) (*Project, error) {

	// Read project file
	bs, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading project: %w", err)
	}

	// Parse project file
	p := NewProject()
	if err := json.Unmarshal(bs, p); err != nil {
		return nil, fmt.Errorf("error parsing project: %w", err)
	}

	// Save project path
	p.Info.Path = path

	return p, nil
}

func (p *Project) Save(path string) (*Project, error) {

	if p == nil {
		return nil, fmt.Errorf("project not loaded")
	}

	// Update path and time
	p.Info.Path = path
	p.Info.Date = time.Now().Format(time.RFC850)

	// Convert project to json
	bs, err := json.MarshalIndent(p, "", "\t")
	if err != nil {
		return nil, fmt.Errorf("error encoding project: %w", err)
	}

	// Write project file
	err = os.WriteFile(path, bs, 0777)
	if err != nil {
		return nil, fmt.Errorf("error reading project: %w", err)
	}

	return p, nil
}
