package project

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ProjectFile represents a HMI/Web project file
type ProjectFile struct {
	Version     string            `json:"version"`
	ProjectType string            `json:"projectType"` // "hmi" or "web"
	ProjectName string            `json:"projectName"`
	Canvas      CanvasConfig      `json:"canvas"`
	Components  []ComponentDef    `json:"components"`
	DataBindings []DataBindingDef `json:"dataBindings"`
}

// CanvasConfig defines the canvas properties
type CanvasConfig struct {
	Width            int    `json:"width"`
	Height           int    `json:"height"`
	BackgroundColor string `json:"backgroundColor,omitempty"`
	Responsive       bool   `json:"responsive,omitempty"`
}

// ComponentDef defines a component in the project
type ComponentDef struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	X          int                    `json:"x"`
	Y          int                    `json:"y"`
	Width      int                    `json:"width"`
	Height     int                    `json:"height"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// DataBindingDef defines data binding for a component
type DataBindingDef struct {
	ComponentID string `json:"componentId"`
	Property    string `json:"property"`
	Variable    string `json:"variable"`   // For HMI: variable ID
	API         string `json:"api,omitempty"` // For Web: API endpoint
}

// Loader loads project files from disk
type Loader struct {
	basePath string
}

// NewLoader creates a new project loader
func NewLoader(basePath string) *Loader {
	return &Loader{
		basePath: basePath,
	}
}

// Load loads a project file from disk
func (l *Loader) Load(filename string) (*ProjectFile, error) {
 fullPath := filepath.Join(l.basePath, filename)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read project file: %w", err)
	}

	var project ProjectFile
	if err := json.Unmarshal(data, &project); err != nil {
		return nil, fmt.Errorf("failed to parse project file: %w", err)
	}

	// Validate project file
	if err := l.validate(&project); err != nil {
		return nil, fmt.Errorf("invalid project file: %w", err)
	}

	return &project, nil
}

// validate validates the project file
func (l *Loader) validate(project *ProjectFile) error {
	if project.Version == "" {
		return fmt.Errorf("missing version")
	}

	if project.ProjectType != "hmi" && project.ProjectType != "web" {
		return fmt.Errorf("invalid project type: %s (must be 'hmi' or 'web')", project.ProjectType)
	}

	if project.Canvas.Width <= 0 || project.Canvas.Height <= 0 {
		return fmt.Errorf("invalid canvas dimensions")
	}

	return nil
}

// Save saves a project file to disk
func (l *Loader) Save(project *ProjectFile, filename string) error {
	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal project: %w", err)
	}

	fullPath := filepath.Join(l.basePath, filename)

	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write project file: %w", err)
	}

	return nil
}

// List lists all project files in the base directory
func (l *Loader) List() ([]string, error) {
	var projects []string

	err := filepath.Walk(l.basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Check if file has .hmi or .web extension
		ext := filepath.Ext(path)
		if ext == ".hmi" || ext == ".web" {
			relPath, err := filepath.Rel(l.basePath, path)
			if err != nil {
				return err
			}
			projects = append(projects, relPath)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}

	return projects, nil
}
