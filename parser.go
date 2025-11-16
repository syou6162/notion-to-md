package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jomei/notionapi"
)

// extractBlockID extracts block ID from a Notion URL or returns the input as-is if it's already an ID
func extractBlockID(input string) (notionapi.BlockID, error) {
	// If it's a URL, extract the ID
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		// Match 32-character hex string (without hyphens)
		re := regexp.MustCompile(`([a-f0-9]{32})`)
		matches := re.FindStringSubmatch(input)
		if len(matches) < 2 {
			return "", fmt.Errorf("invalid Notion URL: cannot extract block ID")
		}
		// Convert to UUID format with hyphens
		id := matches[1]
		formatted := fmt.Sprintf("%s-%s-%s-%s-%s",
			id[0:8], id[8:12], id[12:16], id[16:20], id[20:32])
		return notionapi.BlockID(formatted), nil
	}

	// Assume it's already a block ID
	return notionapi.BlockID(input), nil
}
