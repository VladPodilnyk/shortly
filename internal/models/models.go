package models

import (
	"fmt"
)

type UserRequest struct {
	Url   string `json:"url"`
	Alias string `json:"alias,omitempty"`
}

type DecodedUrl struct {
	OriginalUrl string `json:"original_url"`
}

type EncodedUrl struct {
	ShortUrl string `json:"short_url"`
}

type SystemInfo struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

func (status *SystemInfo) Show() string {
	return fmt.Sprintf("(status: %s environment: %s version: %s)", status.Status, status.Environment, status.Version)
}
