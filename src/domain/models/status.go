package models

type Status struct {
	AppName string `json:"appName"`
	Version string `json:"version"`
	Status  string `json:"status"`
	Date    string `json:"date"`
}
