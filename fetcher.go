package main

import (
	"context"
	"fmt"
	"strings"

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

	// Fetch children blocks with pagination
	blocks, err := fetchBlockChildren(ctx, fetcher, blockID)
	if err != nil {
		return nil, err
	}

	var result []BlockWithIndent
	for _, block := range blocks {
		// Add current block
		result = append(result, BlockWithIndent{
			Block:  block,
			Indent: depth,
		})

		// Recursively fetch children if HasChildren is true
		if block.GetHasChildren() {
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

	for {
		resp, err := fetcher.GetChildren(ctx, blockID, pagination)
		if err != nil {
			return nil, fmt.Errorf("failed to get children for block %s: %w", blockID, err)
		}

		allBlocks = append(allBlocks, resp.Results...)

		if !resp.HasMore {
			break
		}
		pagination.StartCursor = notionapi.Cursor(resp.NextCursor)
	}

	return allBlocks, nil
}

// fetchPageInfo fetches page metadata from Notion API
func fetchPageInfo(ctx context.Context, client *notionapi.Client, pageID notionapi.PageID) (PageInfo, error) {
	page, err := client.Page.Get(ctx, pageID)
	if err != nil {
		return PageInfo{}, fmt.Errorf("failed to get page info: %w", err)
	}

	// Extract title from properties
	var title string
	for _, prop := range page.Properties {
		if titleProp, ok := prop.(*notionapi.TitleProperty); ok {
			var titleBuilder strings.Builder
			for _, rt := range titleProp.Title {
				titleBuilder.WriteString(rt.PlainText)
			}
			title = titleBuilder.String()
			break
		}
	}

	return PageInfo{
		Title:          title,
		URL:            page.URL,
		CreatedTime:    page.CreatedTime,
		LastEditedTime: page.LastEditedTime,
	}, nil
}
