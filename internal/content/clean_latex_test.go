package content

import "testing"

func TestCleanLaTeX(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "bold text",
			input:    "This is \\textbf{bold} text",
			expected: "This is bold text",
		},
		{
			name:     "italic text",
			input:    "This is \\textit{italic} text",
			expected: "This is italic text",
		},
		{
			name:     "emphasis",
			input:    "This is \\emph{emphasized} text",
			expected: "This is emphasized text",
		},
		{
			name:     "strikethrough",
			input:    "This is \\sout{strikethrough} text",
			expected: "This is strikethrough text",
		},
		{
			name:     "underline",
			input:    "This is \\underline{underlined} text",
			expected: "This is underlined text",
		},
		{
			name:     "code",
			input:    "Use \\texttt{code} in text",
			expected: "Use code in text",
		},
		{
			name:     "verbatim",
			input:    "Use \\verb|code| in text",
			expected: "Use code in text",
		},
		{
			name:     "href link",
			input:    "Visit \\href{https://google.com}{Google}",
			expected: "Visit Google",
		},
		{
			name:     "url",
			input:    "Visit \\url{https://google.com}",
			expected: "Visit https://google.com",
		},
		{
			name:     "section",
			input:    "\\section{Title}",
			expected: "Title",
		},
		{
			name:     "subsection",
			input:    "\\subsection{Subtitle}",
			expected: "Subtitle",
		},
		{
			name:     "itemize list",
			input:    "\\begin{itemize}\\item First\\item Second\\end{itemize}",
			expected: "First\nSecond",
		},
		{
			name:     "enumerate list",
			input:    "\\begin{enumerate}\\item First\\item Second\\end{enumerate}",
			expected: "First\nSecond",
		},
		{
			name:     "quote",
			input:    "\\begin{quote}This is a quote\\end{quote}",
			expected: "This is a quote",
		},
		{
			name:     "inline math",
			input:    "The formula $E = mc^2$ is famous",
			expected: "The formula is famous",
		},
		{
			name:     "display math",
			input:    "The formula \\[E = mc^2\\] is famous",
			expected: "The formula is famous",
		},
		{
			name:     "comment",
			input:    "Text % this is a comment\nMore text",
			expected: "Text\nMore text",
		},
		{
			name:     "combined formatting",
			input:    "\\section{Title}\\textbf{Bold} and \\textit{italic} text with \\href{url}{link}.",
			expected: "Title\n\nBold and italic text with link.",
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
			name:     "nested commands",
			input:    "\\textbf{\\textit{bold italic}} text",
			expected: "bold italic text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CleanLaTeX(tt.input)
			if result != tt.expected {
				t.Errorf("CleanLaTeX() = %q, want %q", result, tt.expected)
			}
		})
	}
}
