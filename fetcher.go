package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jomei/notionapi"
)

// BlockWithIndent holds a block and its indentation level
type BlockWithIndent struct {
	Block  notionapi.Block
	Indent int
}

// BlockFetcher is an interface for fetching blocks from Notion API
type BlockFetcher interface {
	GetChildren(ctx context.Context, blockID notionapi.BlockID, pagination *notionapi.Pagination) (*notionapi.GetChildrenResponse, error)
}

// fetchAllBlocks fetches all blocks recursively starting from the given block ID
func fetchAllBlocks(ctx context.Context, fetcher BlockFetcher, blockID notionapi.BlockID) ([]BlockWithIndent, error) {
	return fetchAllBlocksRecursive(ctx, fetcher, blockID, 0)
}

// fetchAllBlocksRecursive recursively fetches blocks with depth tracking
func fetchAllBlocksRecursive(ctx context.Context, fetcher BlockFetcher, blockID notionapi.BlockID, depth int) ([]BlockWithIndent, error) {
	const maxDepth = 10
	if depth > maxDepth {
		return nil, fmt.Errorf("maximum recursion depth (%d) exceeded", maxDepth)
	}

	fmt.Fprintf(os.Stderr, "[DEBUG] Fetching children for block %s (depth: %d)\n", blockID, depth)

	// Fetch children blocks with pagination
	blocks, err := fetchBlockChildren(ctx, fetcher, blockID)
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(os.Stderr, "[DEBUG] Found %d children at depth %d\n", len(blocks), depth)

	var result []BlockWithIndent
	for i, block := range blocks {
		// Add current block
		result = append(result, BlockWithIndent{
			Block:  block,
			Indent: depth,
		})

		// Recursively fetch children if HasChildren is true
		if block.GetHasChildren() {
			fmt.Fprintf(os.Stderr, "[DEBUG] Block %d/%d has children, recursing...\n", i+1, len(blocks))
			children, err := fetchAllBlocksRecursive(ctx, fetcher, block.GetID(), depth+1)
			if err != nil {
				return nil, err
			}
			result = append(result, children...)
		}
	}

	return result, nil
}

// fetchBlockChildren fetches children of a block with pagination support
func fetchBlockChildren(ctx context.Context, fetcher BlockFetcher, blockID notionapi.BlockID) ([]notionapi.Block, error) {
	var allBlocks []notionapi.Block
	pagination := &notionapi.Pagination{}
	pageNum := 1

	for {
		fmt.Fprintf(os.Stderr, "[DEBUG] API call: GetChildren page %d for block %s\n", pageNum, blockID)
		resp, err := fetcher.GetChildren(ctx, blockID, pagination)
		if err != nil {
			return nil, fmt.Errorf("failed to get children for block %s: %w", blockID, err)
		}

		fmt.Fprintf(os.Stderr, "[DEBUG] Received %d blocks (HasMore: %v)\n", len(resp.Results), resp.HasMore)
		allBlocks = append(allBlocks, resp.Results...)

		if !resp.HasMore {
			break
		}
		pagination.StartCursor = notionapi.Cursor(resp.NextCursor)
		pageNum++
	}

	return allBlocks, nil
}
