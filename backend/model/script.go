package model

import "time"

type Section struct {
	ID                string `json:"id"`
	Type              string `json:"type"`
	Title             string `json:"title"`
	Content           string `json:"content"`
	EstimatedDuration int    `json:"estimatedDuration"`
}

type Script struct {
	ID                     string    `json:"id"`
	Topic                  string    `json:"topic"`
	Discipline             string    `json:"discipline"`
	TargetAudience         string    `json:"targetAudience"`
	Style                  string    `json:"style"`
	EstimatedTotalDuration int       `json:"estimatedTotalDuration"`
	Sections               []Section `json:"sections"`
	CreatedAt              time.Time `json:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt"`
}

type GenerateRequest struct {
	Topic            string `json:"topic" binding:"required"`
	Discipline       string `json:"discipline"`
	TargetAudience   string `json:"targetAudience"`
	TargetDuration   int    `json:"targetDuration"`
	Style            string `json:"style"`
	CustomStructure  string `json:"customStructure"`
	ReferenceText    string `json:"referenceText"`
	ReferenceFileURL  string `json:"referenceFileURL"`
}
