package auth

import (
	"fmt"
	"os"
	"strings"
)

// GetToken returns the Terraform Cloud API token from the environment.
func GetToken() (string, error) {
	token := os.Getenv("TFC_TOKEN")
	if token == "" {
		return "", fmt.Errorf("TFC_TOKEN not set — export your Terraform Cloud API token")
	}
	return strings.TrimSpace(token), nil
}

// GetAddress returns the Terraform Cloud base URL, defaulting to app.terraform.io.
func GetAddress() string {
	if addr := os.Getenv("TFC_ADDRESS"); addr != "" {
		return strings.TrimRight(addr, "/")
	}
	return "https://app.terraform.io"
}
