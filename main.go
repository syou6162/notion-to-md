package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jomei/notionapi"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: notion-to-md <block-id-or-url>")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Examples:")
		fmt.Fprintln(os.Stderr, "  notion-to-md cec15681-9083-4e1f-a0ae-72d268507aab")
		fmt.Fprintln(os.Stderr, "  notion-to-md https://www.notion.so/10xall/By-name-cec1568190834e1fa0ae72d268507aab")
		os.Exit(1)
	}
	input := os.Args[1]

	// Extract block ID from URL or use directly
	blockID, err := extractBlockID(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Get NOTION_TOKEN
	token := os.Getenv("NOTION_TOKEN")
	if token == "" {
		fmt.Fprintln(os.Stderr, "Error: NOTION_TOKEN environment variable not set")
		os.Exit(1)
	}

	// Initialize Notion client
	client := notionapi.NewClient(notionapi.Token(token))

	// Fetch all blocks recursively
	ctx := context.Background()
	blocks, err := fetchAllBlocks(ctx, client.Block, blockID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching blocks: %v\n", err)
		os.Exit(1)
	}

	// Convert to Markdown
	markdown := convert(blocks)
	fmt.Print(markdown)
}
