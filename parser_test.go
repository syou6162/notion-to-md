package main

import (
	"testing"

	"github.com/jomei/notionapi"
)

func TestExtractBlockIDFromURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected notionapi.BlockID
		wantErr  bool
	}{
		{
			name:     "Valid Notion URL with https",
			input:    "https://www.notion.so/workspace/Page-title-cec1568190834e1fa0ae72d268507aab",
			expected: "cec15681-9083-4e1f-a0ae-72d268507aab",
			wantErr:  false,
		},
		{
			name:     "Valid Notion URL with http",
			input:    "http://www.notion.so/workspace/Page-cec1568190834e1fa0ae72d268507aab",
			expected: "cec15681-9083-4e1f-a0ae-72d268507aab",
			wantErr:  false,
		},
		{
			name:     "Block ID with hyphens",
			input:    "cec15681-9083-4e1f-a0ae-72d268507aab",
			expected: "cec15681-9083-4e1f-a0ae-72d268507aab",
			wantErr:  false,
		},
		{
			name:     "Block ID without hyphens",
			input:    "cec1568190834e1fa0ae72d268507aab",
			expected: "cec1568190834e1fa0ae72d268507aab",
			wantErr:  false,
		},
		{
			name:     "Invalid URL without block ID",
			input:    "https://www.notion.so/workspace/",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Notion URL with different format",
			input:    "https://notion.so/abc123def456789012345678901234ab",
			expected: "abc123de-f456-7890-1234-5678901234ab",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := extractBlockID(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}
