// Package main provides CLI for generating and deploying AWS Lambda functions
// Usage:
//
//	go run main.go generate
//	go run main.go deploy [function-name]
package main

import (
	"fmt"
	"os"

	"goprojectskeleton/aws/functions/utils"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "generate", "gen", "g":
		if err := utils.GenerateFunctions(utils.FunctionsJSONPath); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
			os.Exit(1)
		}

	case "deploy", "d":
		functionFilter := ""
		if len(os.Args) > 2 {
			functionFilter = os.Args[2]
		}
		if err := utils.DeployFunctions(utils.FunctionsJSONPath, functionFilter); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
			os.Exit(1)
		}

	case "help", "h", "-h", "--help":
		printUsage()
		os.Exit(0)

	default:
		fmt.Fprintf(os.Stderr, "❌ Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("AWS Lambda Functions CLI")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  go run main.go <command> [arguments]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  generate, gen, g    Generate all Lambda functions from functions.json")
	fmt.Println("  deploy, d [name]    Deploy function(s) to AWS Lambda")
	fmt.Println("                      If [name] is provided, deploys only that function")
	fmt.Println("                      If omitted, deploys all functions")
	fmt.Println("  help, h             Show this help message")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  go run main.go generate")
	fmt.Println("  go run main.go deploy")
	fmt.Println("  go run main.go deploy health-check")
	fmt.Println("")
	fmt.Println("Environment Variables:")
	fmt.Println("  PROJECT_NAME    Project name (default: go-project-skeleton)")
	fmt.Println("  ENVIRONMENT     Environment name (default: development)")
	fmt.Println("  AWS_REGION      AWS region (default: us-east-1)")
}
