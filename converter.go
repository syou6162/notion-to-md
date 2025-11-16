package main

import (
	"strings"

	"github.com/jomei/notionapi"
)

// convert converts blocks with indentation to Markdown
func convert(blocks []BlockWithIndent) string {
	var result strings.Builder

	for _, bwi := range blocks {
		block := bwi.Block
		indent := strings.Repeat("  ", bwi.Indent) // 2 spaces per indent level

		switch block.GetType() {
		case notionapi.BlockTypeHeading1:
			if h1, ok := block.(*notionapi.Heading1Block); ok {
				text := formatRichText(h1.Heading1.RichText)
				result.WriteString("# " + text + "\n\n")
			}

		case notionapi.BlockTypeHeading2:
			if h2, ok := block.(*notionapi.Heading2Block); ok {
				text := formatRichText(h2.Heading2.RichText)
				result.WriteString("## " + text + "\n\n")
			}

		case notionapi.BlockTypeHeading3:
			if h3, ok := block.(*notionapi.Heading3Block); ok {
				text := formatRichText(h3.Heading3.RichText)
				result.WriteString("### " + text + "\n\n")
			}

		case notionapi.BlockTypeParagraph:
			if p, ok := block.(*notionapi.ParagraphBlock); ok {
				text := formatRichText(p.Paragraph.RichText)
				if text != "" {
					result.WriteString(text + "\n\n")
				}
			}

		case notionapi.BlockTypeBulletedListItem:
			if bl, ok := block.(*notionapi.BulletedListItemBlock); ok {
				text := formatRichText(bl.BulletedListItem.RichText)
				result.WriteString(indent + "- " + text + "\n")
			}

		case notionapi.BlockTypeNumberedListItem:
			if nl, ok := block.(*notionapi.NumberedListItemBlock); ok {
				text := formatRichText(nl.NumberedListItem.RichText)
				result.WriteString(indent + "1. " + text + "\n")
			}

		case notionapi.BlockTypeCode:
			if c, ok := block.(*notionapi.CodeBlock); ok {
				text := formatRichText(c.Code.RichText)
				lang := string(c.Code.Language)
				result.WriteString("```" + lang + "\n")
				result.WriteString(text + "\n")
				result.WriteString("```\n\n")
			}

		case notionapi.BlockTypeToggle:
			if t, ok := block.(*notionapi.ToggleBlock); ok {
				text := formatRichText(t.Toggle.RichText)
				result.WriteString(indent + "- " + text + "\n")
			}

		case notionapi.BlockTypeQuote:
			if q, ok := block.(*notionapi.QuoteBlock); ok {
				text := formatRichText(q.Quote.RichText)
				result.WriteString("> " + text + "\n\n")
			}

		case notionapi.BlockTypeDivider:
			result.WriteString("---\n\n")

		case notionapi.BlockTypeCallout:
			if c, ok := block.(*notionapi.CalloutBlock); ok {
				text := formatRichText(c.Callout.RichText)
				result.WriteString("> " + text + "\n\n")
			}
		}
	}

	return result.String()
}

// formatRichText converts Notion RichText to Markdown with annotations
func formatRichText(richTexts []notionapi.RichText) string {
	var result strings.Builder

	for _, rt := range richTexts {
		text := rt.PlainText

		// Apply annotations in order: code, bold, italic, strikethrough
		if rt.Annotations != nil {
			if rt.Annotations.Code {
				text = "`" + text + "`"
			}
			if rt.Annotations.Bold {
				text = "**" + text + "**"
			}
			if rt.Annotations.Italic {
				text = "*" + text + "*"
			}
			if rt.Annotations.Strikethrough {
				text = "~~" + text + "~~"
			}
		}

		// Apply link
		if rt.Href != "" {
			text = "[" + text + "](" + rt.Href + ")"
		}

		result.WriteString(text)
	}

	return result.String()
}
