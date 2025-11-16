package main

import (
	"context"
	"errors"
	"testing"

	"github.com/jomei/notionapi"
)

// mockBlockFetcher is a mock implementation of BlockFetcher for testing
type mockBlockFetcher struct {
	responses []*notionapi.GetChildrenResponse
	callCount int
	err       error
}

func (m *mockBlockFetcher) GetChildren(ctx context.Context, blockID notionapi.BlockID, pagination *notionapi.Pagination) (*notionapi.GetChildrenResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.callCount >= len(m.responses) {
		return &notionapi.GetChildrenResponse{
			Results: []notionapi.Block{},
			HasMore: false,
		}, nil
	}
	resp := m.responses[m.callCount]
	m.callCount++
	return resp, nil
}

// Helper function to create a simple paragraph block
func createParagraphBlock(id string, hasChildren bool) notionapi.Block {
	return &notionapi.ParagraphBlock{
		BasicBlock: notionapi.BasicBlock{
			Object:      "block",
			ID:          notionapi.BlockID(id),
			Type:        notionapi.BlockTypeParagraph,
			HasChildren: hasChildren,
		},
		Paragraph: notionapi.Paragraph{
			RichText: []notionapi.RichText{
				{PlainText: "Test block"},
			},
		},
	}
}

func TestFetchBlockChildrenSinglePage(t *testing.T) {
	ctx := context.Background()
	blockID := notionapi.BlockID("test-block-id")

	blocks := []notionapi.Block{
		createParagraphBlock("block-1", false),
		createParagraphBlock("block-2", false),
	}

	mock := &mockBlockFetcher{
		responses: []*notionapi.GetChildrenResponse{
			{
				Results: blocks,
				HasMore: false,
			},
		},
	}

	result, err := fetchBlockChildren(ctx, mock, blockID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 blocks, got %d", len(result))
	}
}

func TestFetchBlockChildrenMultiplePages(t *testing.T) {
	ctx := context.Background()
	blockID := notionapi.BlockID("test-block-id")

	mock := &mockBlockFetcher{
		responses: []*notionapi.GetChildrenResponse{
			{
				Results:    []notionapi.Block{createParagraphBlock("block-1", false)},
				HasMore:    true,
				NextCursor: "cursor-1",
			},
			{
				Results:    []notionapi.Block{createParagraphBlock("block-2", false)},
				HasMore:    true,
				NextCursor: "cursor-2",
			},
			{
				Results: []notionapi.Block{createParagraphBlock("block-3", false)},
				HasMore: false,
			},
		},
	}

	result, err := fetchBlockChildren(ctx, mock, blockID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 blocks, got %d", len(result))
	}

	if mock.callCount != 3 {
		t.Errorf("Expected 3 API calls, got %d", mock.callCount)
	}
}

func TestFetchBlockChildrenEmpty(t *testing.T) {
	ctx := context.Background()
	blockID := notionapi.BlockID("test-block-id")

	mock := &mockBlockFetcher{
		responses: []*notionapi.GetChildrenResponse{
			{
				Results: []notionapi.Block{},
				HasMore: false,
			},
		},
	}

	result, err := fetchBlockChildren(ctx, mock, blockID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected 0 blocks, got %d", len(result))
	}
}

func TestFetchBlockChildrenError(t *testing.T) {
	ctx := context.Background()
	blockID := notionapi.BlockID("test-block-id")

	mock := &mockBlockFetcher{
		err: errors.New("API error"),
	}

	_, err := fetchBlockChildren(ctx, mock, blockID)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestFetchAllBlocksRecursiveNoChildren(t *testing.T) {
	ctx := context.Background()
	blockID := notionapi.BlockID("test-block-id")

	blocks := []notionapi.Block{
		createParagraphBlock("block-1", false),
		createParagraphBlock("block-2", false),
	}

	mock := &mockBlockFetcher{
		responses: []*notionapi.GetChildrenResponse{
			{
				Results: blocks,
				HasMore: false,
			},
		},
	}

	result, err := fetchAllBlocksRecursive(ctx, mock, blockID, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 blocks, got %d", len(result))
	}

	for i, bwi := range result {
		if bwi.Indent != 0 {
			t.Errorf("Block %d: expected indent 0, got %d", i, bwi.Indent)
		}
	}
}

func TestFetchAllBlocksRecursiveWithChildren(t *testing.T) {
	ctx := context.Background()
	blockID := notionapi.BlockID("test-block-id")

	// First call: parent blocks
	// Second call: children of block-1
	mock := &mockBlockFetcher{
		responses: []*notionapi.GetChildrenResponse{
			{
				Results: []notionapi.Block{
					createParagraphBlock("block-1", true),
					createParagraphBlock("block-2", false),
				},
				HasMore: false,
			},
			{
				Results: []notionapi.Block{
					createParagraphBlock("child-1", false),
				},
				HasMore: false,
			},
		},
	}

	result, err := fetchAllBlocksRecursive(ctx, mock, blockID, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// block-1 (indent 0), child-1 (indent 1), block-2 (indent 0)
	if len(result) != 3 {
		t.Errorf("Expected 3 blocks, got %d", len(result))
	}

	expectedIndents := []int{0, 1, 0}
	for i, bwi := range result {
		if bwi.Indent != expectedIndents[i] {
			t.Errorf("Block %d: expected indent %d, got %d", i, expectedIndents[i], bwi.Indent)
		}
	}
}

func TestFetchAllBlocksRecursiveDeepNesting(t *testing.T) {
	ctx := context.Background()
	blockID := notionapi.BlockID("test-block-id")

	// Create a chain: parent -> child -> grandchild
	mock := &mockBlockFetcher{
		responses: []*notionapi.GetChildrenResponse{
			{
				Results: []notionapi.Block{createParagraphBlock("parent", true)},
				HasMore: false,
			},
			{
				Results: []notionapi.Block{createParagraphBlock("child", true)},
				HasMore: false,
			},
			{
				Results: []notionapi.Block{createParagraphBlock("grandchild", false)},
				HasMore: false,
			},
		},
	}

	result, err := fetchAllBlocksRecursive(ctx, mock, blockID, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 blocks, got %d", len(result))
	}

	expectedIndents := []int{0, 1, 2}
	for i, bwi := range result {
		if bwi.Indent != expectedIndents[i] {
			t.Errorf("Block %d: expected indent %d, got %d", i, expectedIndents[i], bwi.Indent)
		}
	}
}

func TestFetchAllBlocksRecursiveMaxDepthExceeded(t *testing.T) {
	ctx := context.Background()
	blockID := notionapi.BlockID("test-block-id")

	mock := &mockBlockFetcher{
		responses: []*notionapi.GetChildrenResponse{
			{
				Results: []notionapi.Block{createParagraphBlock("block", false)},
				HasMore: false,
			},
		},
	}

	// Depth 11 exceeds max depth of 10
	_, err := fetchAllBlocksRecursive(ctx, mock, blockID, 11)
	if err == nil {
		t.Fatal("Expected error for max depth exceeded, got nil")
	}

	if err.Error() != "maximum recursion depth (10) exceeded" {
		t.Errorf("Expected max depth error, got: %v", err)
	}
}

func TestFetchAllBlocksRecursiveError(t *testing.T) {
	ctx := context.Background()
	blockID := notionapi.BlockID("test-block-id")

	mock := &mockBlockFetcher{
		err: errors.New("API error"),
	}

	_, err := fetchAllBlocksRecursive(ctx, mock, blockID, 0)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestFetchAllBlocks(t *testing.T) {
	ctx := context.Background()
	blockID := notionapi.BlockID("test-block-id")

	blocks := []notionapi.Block{
		createParagraphBlock("block-1", false),
	}

	mock := &mockBlockFetcher{
		responses: []*notionapi.GetChildrenResponse{
			{
				Results: blocks,
				HasMore: false,
			},
		},
	}

	result, err := fetchAllBlocks(ctx, mock, blockID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 block, got %d", len(result))
	}

	if result[0].Indent != 0 {
		t.Errorf("Expected indent 0, got %d", result[0].Indent)
	}
}
