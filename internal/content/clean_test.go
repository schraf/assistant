package content

import "testing"

func TestCleanText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "markdown formatting",
			input:    "This is **bold** and *italic* text",
			expected: "This is bold and italic text",
		},
		{
			name:     "html formatting",
			input:    "This is <strong>bold</strong> and <em>italic</em> text",
			expected: "This is bold and italic text",
		},
		{
			name:     "latex formatting",
			input:    "This is \\textbf{bold} and \\textit{italic} text",
			expected: "This is bold and italic text",
		},
		{
			name:     "escape characters",
			input:    "Line 1\nLine 2\tTabbed",
			expected: "Line 1Line 2 Tabbed",
		},
		{
			name:     "combined markdown and html",
			input:    "**Bold** text with <em>italic</em>",
			expected: "Bold text with italic",
		},
		{
			name:     "combined all formats",
			input:    "**Bold** <strong>HTML</strong> \\textbf{LaTeX} text\nwith escapes",
			expected: "Bold HTML LaTeX textwith escapes",
		},
		{
			name:     "markdown link",
			input:    "Visit [Google](https://google.com)",
			expected: "Visit Google",
		},
		{
			name:     "html link",
			input:    "Visit <a href=\"https://google.com\">Google</a>",
			expected: "Visit Google",
		},
		{
			name:     "latex link",
			input:    "Visit \\href{https://google.com}{Google}",
			expected: "Visit Google",
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
			name:     "markdown headers",
			input:    "# Header 1\n## Header 2",
			expected: "",
		},
		{
			name:     "html headers",
			input:    "<h1>Header 1</h1><h2>Header 2</h2>",
			expected: "Header 1Header 2",
		},
		{
			name:     "latex section",
			input:    "\\section{Title}",
			expected: "Title",
		},
		{
			name:     "mixed escapes and formatting",
			input:    "**Bold**\n<em>Italic</em>\t\\textbf{LaTeX}",
			expected: "BoldItalic LaTeX",
		},
		{
			name:     "preserve special characters",
			input:    `Text with "quotes" and 'apostrophes' and ? ! $ % & ( ) + = - characters`,
			expected: `Text with "quotes" and 'apostrophes' and ? ! $ % & ( ) + = - characters`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CleanText(tt.input)
			if result != tt.expected {
				t.Errorf("CleanText() = %q, want %q", result, tt.expected)
			}
		})
	}
}
