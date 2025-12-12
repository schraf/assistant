package content

import "testing"

func TestCleanMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "bold text",
			input:    "This is **bold** text",
			expected: "This is bold text",
		},
		{
			name:     "bold with underscores",
			input:    "This is __bold__ text",
			expected: "This is bold text",
		},
		{
			name:     "italic text",
			input:    "This is *italic* text",
			expected: "This is italic text",
		},
		{
			name:     "italic with underscores",
			input:    "This is _italic_ text",
			expected: "This is italic text",
		},
		{
			name:     "strikethrough",
			input:    "This is ~~strikethrough~~ text",
			expected: "This is strikethrough text",
		},
		{
			name:     "inline code",
			input:    "Use `code` in text",
			expected: "Use code in text",
		},
		{
			name:     "code block",
			input:    "```go\nfunc main() {}\n```",
			expected: "func main() {}",
		},
		{
			name:     "link",
			input:    "Visit [Google](https://google.com)",
			expected: "Visit Google",
		},
		{
			name:     "image",
			input:    "![Alt text](image.png)",
			expected: "Alt text",
		},
		{
			name:     "header",
			input:    "# Header 1\n## Header 2",
			expected: "",
		},
		{
			name:     "list",
			input:    "- Item 1\n- Item 2\n* Item 3",
			expected: "Item 1\nItem 2\nItem 3",
		},
		{
			name:     "numbered list",
			input:    "1. First\n2. Second",
			expected: "First\nSecond",
		},
		{
			name:     "blockquote",
			input:    "> This is a quote",
			expected: "This is a quote",
		},
		{
			name:     "horizontal rule",
			input:    "Text before\n---\nText after",
			expected: "Text before\n\nText after",
		},
		{
			name:     "combined formatting",
			input:    "# Title\n\nThis is **bold** and *italic* text with a [link](url).",
			expected: "This is bold and italic text with a link.",
		},
		{
			name:     "multiple newlines",
			input:    "Text\n\n\n\nMore text",
			expected: "Text\n\nMore text",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "plain text",
			input:    "Just plain text",
			expected: "Just plain text",
		},
		{
			name:     "header with text before and after",
			input:    "Text before\n# Header\nText after",
			expected: "Text before\n\nText after",
		},
		{
			name:     "multiple headers",
			input:    "# Header 1\nSome content\n## Header 2\nMore content\n### Header 3",
			expected: "Some content\n\nMore content",
		},
		{
			name:     "header with six hashes",
			input:    "###### Header 6",
			expected: "",
		},
		{
			name:     "header with content",
			input:    "# Section Title\n\nThis is paragraph text.",
			expected: "This is paragraph text.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CleanMarkdown(tt.input)
			if result != tt.expected {
				t.Errorf("CleanMarkdown() = %q, want %q", result, tt.expected)
			}
		})
	}
}
