package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// DeploymentConfig holds deployment configuration
type DeploymentConfig struct {
	ProjectName string
	Environment string
	Region      string
}

// DeployFunctions deploys Lambda functions to AWS
func DeployFunctions(functionsJSONPath string, functionFilter string) error {
	// Read functions.json
	functionsJSON, err := os.ReadFile(functionsJSONPath)
	if err != nil {
		return fmt.Errorf("error reading %s: %w", functionsJSONPath, err)
	}

	var functions []FunctionConfig
	if err := json.Unmarshal(functionsJSON, &functions); err != nil {
		return fmt.Errorf("error parsing %s: %w", functionsJSONPath, err)
	}

	// Get deployment config from environment or defaults
	config := DeploymentConfig{
		ProjectName: getEnvOrDefault("PROJECT_NAME", "go-project-skeleton"),
		Environment: getEnvOrDefault("ENVIRONMENT", "development"),
		Region:      getEnvOrDefault("AWS_REGION", "us-east-1"),
	}

	// Calculate name prefix (same logic as Terraform)
	projectNameClean := strings.ReplaceAll(strings.ReplaceAll(config.ProjectName, "-", ""), "_", "")
	envPrefix := config.Environment
	if len(envPrefix) > 3 {
		envPrefix = envPrefix[:3]
	}
	namePrefix := fmt.Sprintf("%s-%s", projectNameClean, envPrefix)

	// Filter functions if filter provided
	functionsToDeploy := functions
	if functionFilter != "" {
		functionsToDeploy = filterFunctions(functions, functionFilter)
		if len(functionsToDeploy) == 0 {
			return fmt.Errorf("function '%s' not found", functionFilter)
		}
	}

	fmt.Printf("🚀 Deploying %d function(s) to AWS Lambda\n", len(functionsToDeploy))
	fmt.Printf("📦 Project: %s | Environment: %s | Region: %s\n", config.ProjectName, config.Environment, config.Region)
	fmt.Printf("🏷️  Name prefix: %s\n\n", namePrefix)

	// Deploy each function
	successCount := 0
	failCount := 0

	for _, fn := range functionsToDeploy {
		functionName := getLambdaFunctionName(namePrefix, fn.Name)
		if err := deployFunction(fn, functionName, config.Region); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to deploy %s: %v\n\n", fn.Path, err)
			failCount++
			continue
		}
		fmt.Printf("✅ Successfully deployed: %s\n\n", fn.Path)
		successCount++
	}

	// Summary
	fmt.Println("=" + strings.Repeat("=", 50))
	fmt.Printf("📊 Deployment Summary:\n")
	fmt.Printf("   ✅ Success: %d\n", successCount)
	fmt.Printf("   ❌ Failed:  %d\n", failCount)
	fmt.Println("=" + strings.Repeat("=", 50))

	if failCount > 0 {
		return fmt.Errorf("deployment completed with %d failure(s)", failCount)
	}

	return nil
}

func deployFunction(fn FunctionConfig, lambdaFunctionName, region string) error {
	fmt.Printf("📦 Deploying: %s → %s\n", fn.Path, lambdaFunctionName)

	// Get absolute path to function directory
	funcDir, err := filepath.Abs("tmp/" + fn.Path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Check if function directory exists
	if _, err := os.Stat(funcDir); os.IsNotExist(err) {
		return fmt.Errorf("function directory does not exist: %s", funcDir)
	}

	// Step 1: go mod tidy
	fmt.Printf("   🔧 Running go mod tidy...\n")
	if err := runCommand(funcDir, "go", "mod", "tidy"); err != nil {
		return fmt.Errorf("go mod tidy failed: %w", err)
	}

	// Ensure bin directory exists
	binDir := filepath.Join(funcDir, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	// Step 2: Build the function (bootstrap binary)
	fmt.Printf("   🔨 Building function...\n")
	bootstrapPath := filepath.Join(binDir, "bootstrap")
	buildCmd := exec.Command("go", "build",
		"-o", bootstrapPath,
		"-tags", "lambda.norpc",
		"main.go")
	buildCmd.Dir = funcDir
	buildCmd.Env = append(os.Environ(),
		"GOOS=linux",
		"GOARCH=amd64",
		"CGO_ENABLED=0",
	)

	if output, err := buildCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("build failed: %w\nOutput: %s", err, string(output))
	}

	// Step 3: Create zip file with bootstrap only
	// Templates are now stored in S3 and loaded at runtime
	fmt.Printf("   📦 Creating zip file...\n")
	zipName := fmt.Sprintf("%s.zip", fn.Name)
	zipPath := filepath.Join(binDir, zipName)

	// Run zip from bin directory so that entries are:
	//   - bootstrap
	zipCmd := exec.Command("zip", "-r", zipName, "bootstrap")
	zipCmd.Dir = binDir
	if output, err := zipCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("zip creation failed: %w\nOutput: %s", err, string(output))
	}

	// Step 4: Deploy to Lambda
	fmt.Printf("   ☁️  Uploading to AWS Lambda...\n")
	updateCmd := exec.Command("aws", "lambda", "update-function-code",
		"--function-name", lambdaFunctionName,
		"--zip-file", fmt.Sprintf("fileb://%s", zipPath),
		"--region", region,
	)
	if output, err := updateCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("lambda update failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

func runCommand(dir, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func getLambdaFunctionName(namePrefix, functionName string) string {
	// Remove hyphens from function name (same as Terraform)
	cleanName := strings.ReplaceAll(functionName, "-", "")
	return fmt.Sprintf("%s-%s", namePrefix, cleanName)
}

func filterFunctions(functions []FunctionConfig, name string) []FunctionConfig {
	var filtered []FunctionConfig
	for _, fn := range functions {
		if fn.Name == name || strings.Contains(fn.Name, name) || strings.Contains(fn.Path, name) {
			filtered = append(filtered, fn)
		}
	}
	return filtered
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
