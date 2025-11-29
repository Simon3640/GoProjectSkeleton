package utils

// FunctionConfig represents a function configuration from functions.json
type FunctionConfig struct {
	Name          string `json:"name"`
	Path          string `json:"path"`
	Handler       string `json:"handler"`
	Route         string `json:"route"`
	Method        string `json:"method"`
	AuthLevel     string `json:"authLevel"`
	NeedsAuth     bool   `json:"needsAuth"`
	NeedsQuery    bool   `json:"needsQuery"`
	HasPathParams bool   `json:"hasPathParams"`
	PathParamName string `json:"pathParamName"`
}

const FunctionsJSONPath = "../functions.json"
