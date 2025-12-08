package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// TemplateData holds data for template generation
type TemplateData struct {
	HandlerName  string
	NeedsAuth    bool
	ModulePath   string
	FunctionPath string
}

const templatesDir = "templates"

// GenerateFunctions generates all Lambda functions from functions.json
func GenerateFunctions(functionsJSONPath string) error {
	// Read functions.json
	functionsJSON, err := os.ReadFile(functionsJSONPath)
	if err != nil {
		return fmt.Errorf("error reading %s: %w", functionsJSONPath, err)
	}

	var functions []FunctionConfig
	if err := json.Unmarshal(functionsJSON, &functions); err != nil {
		return fmt.Errorf("error parsing %s: %w", functionsJSONPath, err)
	}

	// Generate each function
	for _, fn := range functions {
		if err := GenerateFunction(fn); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating function %s: %v\n", fn.Path, err)
			continue
		}
		fmt.Printf("✅ Generated function: %s (%s)\n", fn.Path, fn.Handler)
	}

	fmt.Println("\n🎉 All functions generated successfully!")
	return nil
}

// GenerateFunction generates a single Lambda function
func GenerateFunction(fn FunctionConfig) error {
	// Calculate module path (e.g., "status/health_check" -> "github.com/simon3640/goprojectskeleton/functions/aws/status/health_check")
	modulePath := fmt.Sprintf("github.com/simon3640/goprojectskeleton/functions/aws/%s", strings.ReplaceAll(fn.Path, "/", "/"))

	// Create function directory
	funcDir := filepath.Join("tmp", fn.Path)
	if err := os.MkdirAll(funcDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create bin directory for compiled artifacts
	binDir := filepath.Join(funcDir, "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	// Prepare template data
	data := TemplateData{
		HandlerName:  fn.Handler,
		NeedsAuth:    fn.NeedsAuth,
		ModulePath:   modulePath,
		FunctionPath: fn.Path,
	}

	// Generate main.go
	if err := generateFromTemplate(funcDir, "main.go", "main.go.tmpl", data); err != nil {
		return fmt.Errorf("failed to generate main.go: %w", err)
	}

	// Generate go.mod
	if err := generateFromTemplate(funcDir, "go.mod", "go.mod.tmpl", data); err != nil {
		return fmt.Errorf("failed to generate go.mod: %w", err)
	}

	return nil
}

func generateFromTemplate(dir, filename, templateName string, data TemplateData) error {
	// Read template file
	templatePath := filepath.Join(templatesDir, templateName)
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("error reading template %s: %w", templatePath, err)
	}

	// Parse template
	tmpl, err := template.New(templateName).Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("error parsing template %s: %w", templateName, err)
	}

	// Create output file
	filePath := filepath.Join(dir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", filePath, err)
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("error executing template %s: %w", templateName, err)
	}

	return nil
}
