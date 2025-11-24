package functions

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const templatesDir = "templates"

func GenerateFunctions() {
	baseDir := "."
	if len(os.Args) > 1 {
		baseDir = os.Args[1]
	}

	for _, fn := range GetAll() {
		// Prepare template data
		data := map[string]interface{}{
			"Path":           fn.Path,
			"HandlerName":    fn.HandlerName,
			"Route":          fn.Route,
			"Method":         fn.Method,
			"MethodUpper":    strings.Title(fn.Method), // Title case
			"AuthLevel":      fn.AuthLevel,
			"NeedsAuth":      fn.NeedsAuth,
			"NeedsQuery":     fn.NeedsQuery,
			"HasPathParams":  fn.HasPathParams,
			"PathParamName":  fn.PathParamName,
			"PathParamRoute": fn.PathParamRoute,
		}

		funcDir := filepath.Join(baseDir, "tmp", fn.Path)

		// Obtener el nombre de la función (último segmento del path)
		// Ejemplo: "status/health_check" -> "health_check"
		functionName := filepath.Base(fn.Path)

		// Create directory structure
		// La carpeta debe tener el nombre de la función para evitar sobreescritura en Azure
		functionFolder := filepath.Join(funcDir, functionName)
		os.MkdirAll(functionFolder, 0755)
		os.MkdirAll(filepath.Join(funcDir, "bin"), 0755)

		// Generate main.go
		generateFileFromTemplate(funcDir, "main.go", "main.go.tmpl", data)

		// Generate host.json
		generateFileFromTemplate(funcDir, "host.json", "host.json.tmpl", data)

		// Generate local.settings.json
		generateFileFromTemplate(funcDir, "local.settings.json", "local.settings.json.tmpl", data)

		// Generate functionName/function.json (carpeta con nombre de función)
		generateFileFromTemplate(functionFolder, "function.json", "function.json.tmpl", data)

		// Generate go.mod
		generateFileFromTemplate(funcDir, "go.mod", "go.mod.tmpl", data)

		fmt.Printf("Generated function: %s\n", fn.Path)
	}

	fmt.Println("\nAll functions generated successfully!")
}

func generateFileFromTemplate(dir, filename, templateName string, data map[string]interface{}) {
	// Read template file
	templatePath := filepath.Join(templatesDir, templateName)
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading template %s: %v\n", templatePath, err)
		return
	}

	// Parse template
	tmpl, err := template.New(templateName).Parse(string(templateContent))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing template %s: %v\n", templateName, err)
		return
	}

	// Create output file
	filePath := filepath.Join(dir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file %s: %v\n", filePath, err)
		return
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, data); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing template %s: %v\n", templateName, err)
		return
	}
}
