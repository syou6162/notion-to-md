package main

import (
	"testing"

	"github.com/jomei/notionapi"
)

func TestConvertHeading1(t *testing.T) {
	blocks := []BlockWithIndent{
		{
			Block: &notionapi.Heading1Block{
				BasicBlock: notionapi.BasicBlock{
					Object: "block",
					Type:   notionapi.BlockTypeHeading1,
				},
				Heading1: notionapi.Heading{
					RichText: []notionapi.RichText{
						{
							PlainText: "Test Heading",
						},
					},
				},
			},
			Indent: 0,
		},
	}

	result := convert(blocks)
	expected := "# Test Heading\n\n"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestConvertHeading2(t *testing.T) {
	blocks := []BlockWithIndent{
		{
			Block: &notionapi.Heading2Block{
				BasicBlock: notionapi.BasicBlock{
					Object: "block",
					Type:   notionapi.BlockTypeHeading2,
				},
				Heading2: notionapi.Heading{
					RichText: []notionapi.RichText{
						{
							PlainText: "Test Heading",
						},
					},
				},
			},
			Indent: 0,
		},
	}

	result := convert(blocks)
	expected := "## Test Heading\n\n"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestConvertHeading3(t *testing.T) {
	blocks := []BlockWithIndent{
		{
			Block: &notionapi.Heading3Block{
				BasicBlock: notionapi.BasicBlock{
					Object: "block",
					Type:   notionapi.BlockTypeHeading3,
				},
				Heading3: notionapi.Heading{
					RichText: []notionapi.RichText{
						{
							PlainText: "Test Heading",
						},
					},
				},
			},
			Indent: 0,
		},
	}

	result := convert(blocks)
	expected := "### Test Heading\n\n"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestConvertParagraph(t *testing.T) {
	blocks := []BlockWithIndent{
		{
			Block: &notionapi.ParagraphBlock{
				BasicBlock: notionapi.BasicBlock{
					Object: "block",
					Type:   notionapi.BlockTypeParagraph,
				},
				Paragraph: notionapi.Paragraph{
					RichText: []notionapi.RichText{
						{
							PlainText: "Test paragraph",
						},
					},
				},
			},
			Indent: 0,
		},
	}

	result := convert(blocks)
	expected := "Test paragraph\n\n"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestConvertBulletedListItem(t *testing.T) {
	blocks := []BlockWithIndent{
		{
			Block: &notionapi.BulletedListItemBlock{
				BasicBlock: notionapi.BasicBlock{
					Object: "block",
					Type:   notionapi.BlockTypeBulletedListItem,
				},
				BulletedListItem: notionapi.ListItem{
					RichText: []notionapi.RichText{
						{
							PlainText: "List item",
						},
					},
				},
			},
			Indent: 0,
		},
	}

	result := convert(blocks)
	expected := "- List item\n"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestConvertNumberedListItem(t *testing.T) {
	blocks := []BlockWithIndent{
		{
			Block: &notionapi.NumberedListItemBlock{
				BasicBlock: notionapi.BasicBlock{
					Object: "block",
					Type:   notionapi.BlockTypeNumberedListItem,
				},
				NumberedListItem: notionapi.ListItem{
					RichText: []notionapi.RichText{
						{
							PlainText: "Numbered item",
						},
					},
				},
			},
			Indent: 0,
		},
	}

	result := convert(blocks)
	expected := "1. Numbered item\n"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestConvertCode(t *testing.T) {
	blocks := []BlockWithIndent{
		{
			Block: &notionapi.CodeBlock{
				BasicBlock: notionapi.BasicBlock{
					Object: "block",
					Type:   notionapi.BlockTypeCode,
				},
				Code: notionapi.Code{
					RichText: []notionapi.RichText{
						{
							PlainText: "func main() {}",
						},
					},
					Language: "go",
				},
			},
			Indent: 0,
		},
	}

	result := convert(blocks)
	expected := "```go\nfunc main() {}\n```\n\n"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestConvertToggle(t *testing.T) {
	blocks := []BlockWithIndent{
		{
			Block: &notionapi.ToggleBlock{
				BasicBlock: notionapi.BasicBlock{
					Object: "block",
					Type:   notionapi.BlockTypeToggle,
				},
				Toggle: notionapi.Toggle{
					RichText: []notionapi.RichText{
						{
							PlainText: "Toggle item",
						},
					},
				},
			},
			Indent: 0,
		},
	}

	result := convert(blocks)
	expected := "- Toggle item\n"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestConvertQuote(t *testing.T) {
	blocks := []BlockWithIndent{
		{
			Block: &notionapi.QuoteBlock{
				BasicBlock: notionapi.BasicBlock{
					Object: "block",
					Type:   notionapi.BlockTypeQuote,
				},
				Quote: notionapi.Quote{
					RichText: []notionapi.RichText{
						{
							PlainText: "Quote text",
						},
					},
				},
			},
			Indent: 0,
		},
	}

	result := convert(blocks)
	expected := "> Quote text\n\n"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestConvertCallout(t *testing.T) {
	blocks := []BlockWithIndent{
		{
			Block: &notionapi.CalloutBlock{
				BasicBlock: notionapi.BasicBlock{
					Object: "block",
					Type:   notionapi.BlockTypeCallout,
				},
				Callout: notionapi.Callout{
					RichText: []notionapi.RichText{
						{
							PlainText: "Callout text",
						},
					},
				},
			},
			Indent: 0,
		},
	}

	result := convert(blocks)
	expected := "> Callout text\n\n"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestConvertDivider(t *testing.T) {
	blocks := []BlockWithIndent{
		{
			Block: &notionapi.DividerBlock{
				BasicBlock: notionapi.BasicBlock{
					Object: "block",
					Type:   notionapi.BlockTypeDivider,
				},
			},
			Indent: 0,
		},
	}

	result := convert(blocks)
	expected := "---\n\n"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestConvertBoldAnnotation(t *testing.T) {
	blocks := []BlockWithIndent{
		{
			Block: &notionapi.ParagraphBlock{
				BasicBlock: notionapi.BasicBlock{
					Object: "block",
					Type:   notionapi.BlockTypeParagraph,
				},
				Paragraph: notionapi.Paragraph{
					RichText: []notionapi.RichText{
						{
							PlainText: "bold text",
							Annotations: &notionapi.Annotations{
								Bold: true,
							},
						},
					},
				},
			},
			Indent: 0,
		},
	}

	result := convert(blocks)
	expected := "**bold text**\n\n"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestConvertItalicAnnotation(t *testing.T) {
	blocks := []BlockWithIndent{
		{
			Block: &notionapi.ParagraphBlock{
				BasicBlock: notionapi.BasicBlock{
					Object: "block",
					Type:   notionapi.BlockTypeParagraph,
				},
				Paragraph: notionapi.Paragraph{
					RichText: []notionapi.RichText{
						{
							PlainText: "italic text",
							Annotations: &notionapi.Annotations{
								Italic: true,
							},
						},
					},
				},
			},
			Indent: 0,
		},
	}

	result := convert(blocks)
	expected := "*italic text*\n\n"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestConvertCodeAnnotation(t *testing.T) {
	blocks := []BlockWithIndent{
		{
			Block: &notionapi.ParagraphBlock{
				BasicBlock: notionapi.BasicBlock{
					Object: "block",
					Type:   notionapi.BlockTypeParagraph,
				},
				Paragraph: notionapi.Paragraph{
					RichText: []notionapi.RichText{
						{
							PlainText: "code text",
							Annotations: &notionapi.Annotations{
								Code: true,
							},
						},
					},
				},
			},
			Indent: 0,
		},
	}

	result := convert(blocks)
	expected := "`code text`\n\n"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestConvertStrikethroughAnnotation(t *testing.T) {
	blocks := []BlockWithIndent{
		{
			Block: &notionapi.ParagraphBlock{
				BasicBlock: notionapi.BasicBlock{
					Object: "block",
					Type:   notionapi.BlockTypeParagraph,
				},
				Paragraph: notionapi.Paragraph{
					RichText: []notionapi.RichText{
						{
							PlainText: "strikethrough text",
							Annotations: &notionapi.Annotations{
								Strikethrough: true,
							},
						},
					},
				},
			},
			Indent: 0,
		},
	}

	result := convert(blocks)
	expected := "~~strikethrough text~~\n\n"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestConvertMultipleAnnotations(t *testing.T) {
	blocks := []BlockWithIndent{
		{
			Block: &notionapi.ParagraphBlock{
				BasicBlock: notionapi.BasicBlock{
					Object: "block",
					Type:   notionapi.BlockTypeParagraph,
				},
				Paragraph: notionapi.Paragraph{
					RichText: []notionapi.RichText{
						{
							PlainText: "bold italic",
							Annotations: &notionapi.Annotations{
								Bold:   true,
								Italic: true,
							},
						},
					},
				},
			},
			Indent: 0,
		},
	}

	result := convert(blocks)
	expected := "***bold italic***\n\n"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestConvertLink(t *testing.T) {
	blocks := []BlockWithIndent{
		{
			Block: &notionapi.ParagraphBlock{
				BasicBlock: notionapi.BasicBlock{
					Object: "block",
					Type:   notionapi.BlockTypeParagraph,
				},
				Paragraph: notionapi.Paragraph{
					RichText: []notionapi.RichText{
						{
							PlainText: "link text",
							Href:      "https://example.com",
						},
					},
				},
			},
			Indent: 0,
		},
	}

	result := convert(blocks)
	expected := "[link text](https://example.com)\n\n"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestConvertNestedList(t *testing.T) {
	blocks := []BlockWithIndent{
		{
			Block: &notionapi.BulletedListItemBlock{
				BasicBlock: notionapi.BasicBlock{
					Object: "block",
					Type:   notionapi.BlockTypeBulletedListItem,
				},
				BulletedListItem: notionapi.ListItem{
					RichText: []notionapi.RichText{
						{
							PlainText: "Parent item",
						},
					},
				},
			},
			Indent: 0,
		},
		{
			Block: &notionapi.BulletedListItemBlock{
				BasicBlock: notionapi.BasicBlock{
					Object: "block",
					Type:   notionapi.BlockTypeBulletedListItem,
				},
				BulletedListItem: notionapi.ListItem{
					RichText: []notionapi.RichText{
						{
							PlainText: "Child item",
						},
					},
				},
			},
			Indent: 1,
		},
	}

	result := convert(blocks)
	expected := "- Parent item\n  - Child item\n"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestConvertMultipleBlocks(t *testing.T) {
	blocks := []BlockWithIndent{
		{
			Block: &notionapi.Heading1Block{
				BasicBlock: notionapi.BasicBlock{
					Object: "block",
					Type:   notionapi.BlockTypeHeading1,
				},
				Heading1: notionapi.Heading{
					RichText: []notionapi.RichText{
						{
							PlainText: "Title",
						},
					},
				},
			},
			Indent: 0,
		},
		{
			Block: &notionapi.ParagraphBlock{
				BasicBlock: notionapi.BasicBlock{
					Object: "block",
					Type:   notionapi.BlockTypeParagraph,
				},
				Paragraph: notionapi.Paragraph{
					RichText: []notionapi.RichText{
						{
							PlainText: "Content",
						},
					},
				},
			},
			Indent: 0,
		},
	}

	result := convert(blocks)
	expected := "# Title\n\nContent\n\n"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestConvertEmptyParagraph(t *testing.T) {
	blocks := []BlockWithIndent{
		{
			Block: &notionapi.ParagraphBlock{
				BasicBlock: notionapi.BasicBlock{
					Object: "block",
					Type:   notionapi.BlockTypeParagraph,
				},
				Paragraph: notionapi.Paragraph{
					RichText: []notionapi.RichText{},
				},
			},
			Indent: 0,
		},
	}

	result := convert(blocks)
	expected := ""

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestConvertMultipleRichTexts(t *testing.T) {
	blocks := []BlockWithIndent{
		{
			Block: &notionapi.ParagraphBlock{
				BasicBlock: notionapi.BasicBlock{
					Object: "block",
					Type:   notionapi.BlockTypeParagraph,
				},
				Paragraph: notionapi.Paragraph{
					RichText: []notionapi.RichText{
						{
							PlainText: "Normal ",
						},
						{
							PlainText: "bold",
							Annotations: &notionapi.Annotations{
								Bold: true,
							},
						},
						{
							PlainText: " text",
						},
					},
				},
			},
			Indent: 0,
		},
	}

	result := convert(blocks)
	expected := "Normal **bold** text\n\n"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}
