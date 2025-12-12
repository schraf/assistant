package content

import "testing"

func TestCleanHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "bold text",
			input:    "This is <b>bold</b> text",
			expected: "This is bold text",
		},
		{
			name:     "bold with strong",
			input:    "This is <strong>bold</strong> text",
			expected: "This is bold text",
		},
		{
			name:     "italic text",
			input:    "This is <i>italic</i> text",
			expected: "This is italic text",
		},
		{
			name:     "italic with em",
			input:    "This is <em>italic</em> text",
			expected: "This is italic text",
		},
		{
			name:     "strikethrough",
			input:    "This is <s>strikethrough</s> text",
			expected: "This is strikethrough text",
		},
		{
			name:     "strikethrough with del",
			input:    "This is <del>deleted</del> text",
			expected: "This is deleted text",
		},
		{
			name:     "underline",
			input:    "This is <u>underlined</u> text",
			expected: "This is underlined text",
		},
		{
			name:     "inline code",
			input:    "Use <code>code</code> in text",
			expected: "Use code in text",
		},
		{
			name:     "code block",
			input:    "<pre>func main() {}</pre>",
			expected: "func main() {}",
		},
		{
			name:     "link",
			input:    "Visit <a href=\"https://google.com\">Google</a>",
			expected: "Visit Google",
		},
		{
			name:     "image with alt",
			input:    "<img src=\"image.png\" alt=\"Alt text\">",
			expected: "Alt text",
		},
		{
			name:     "image without alt",
			input:    "<img src=\"image.png\">",
			expected: "",
		},
		{
			name:     "header",
			input:    "<h1>Header 1</h1><h2>Header 2</h2>",
			expected: "Header 1\n\nHeader 2",
		},
		{
			name:     "list",
			input:    "<ul><li>Item 1</li><li>Item 2</li></ul>",
			expected: "Item 1\n\nItem 2",
		},
		{
			name:     "numbered list",
			input:    "<ol><li>First</li><li>Second</li></ol>",
			expected: "First\n\nSecond",
		},
		{
			name:     "blockquote",
			input:    "<blockquote>This is a quote</blockquote>",
			expected: "This is a quote",
		},
		{
			name:     "horizontal rule",
			input:    "Text before<hr>Text after",
			expected: "Text before\nText after",
		},
		{
			name:     "line break",
			input:    "Line 1<br>Line 2",
			expected: "Line 1\nLine 2",
		},
		{
			name:     "paragraph",
			input:    "<p>First paragraph</p><p>Second paragraph</p>",
			expected: "First paragraph\n\nSecond paragraph",
		},
		{
			name:     "combined formatting",
			input:    "<h1>Title</h1><p>This is <strong>bold</strong> and <em>italic</em> text with a <a href=\"url\">link</a>.</p>",
			expected: "Title\n\nThis is bold and italic text with a link.",
		},
		{
			name:     "html entities",
			input:    "Text &amp; more &lt;tags&gt;",
			expected: "Text & more <tags>",
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
			name:     "nested tags",
			input:    "<p>Text with <strong><em>bold italic</em></strong> content</p>",
			expected: "Text with bold italic content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CleanHTML(tt.input)
			if result != tt.expected {
				t.Errorf("CleanHTML() = %q, want %q", result, tt.expected)
			}
		})
	}
}
