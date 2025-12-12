package content

import "testing"

func TestCleanEscapes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "newline",
			input:    "Line 1\nLine 2",
			expected: "Line 1Line 2",
		},
		{
			name:     "carriage return",
			input:    "Line 1\rLine 2",
			expected: "Line 1Line 2",
		},
		{
			name:     "tab",
			input:    "Column1\tColumn2",
			expected: "Column1Column2",
		},
		{
			name:     "backspace",
			input:    "text\bmore",
			expected: "textmore",
		},
		{
			name:     "form feed",
			input:    "page1\fpage2",
			expected: "page1page2",
		},
		{
			name:     "vertical tab",
			input:    "line1\vline2",
			expected: "line1line2",
		},
		{
			name:     "bell",
			input:    "alert\asound",
			expected: "alertsound",
		},
		{
			name:     "backslash",
			input:    "path\\to\\file",
			expected: "pathtofile",
		},
		{
			name:     "multiple escapes",
			input:    "text\n\t\r\b",
			expected: "text",
		},
		{
			name:     "octal escape",
			input:    "text\\101more",
			expected: "textmore",
		},
		{
			name:     "hex escape",
			input:    "text\\x41more",
			expected: "textmore",
		},
		{
			name:     "unicode escape",
			input:    "text\\u0041more",
			expected: "textmore",
		},
		{
			name:     "unicode long escape",
			input:    "text\\U00000041more",
			expected: "textmore",
		},
		{
			name:     "mixed escapes",
			input:    "Hello\n\tWorld\\x20Test",
			expected: "HelloWorldTest",
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
			name:     "only escapes",
			input:    "\n\t\r\b\f\v\a\\",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CleanEscapes(tt.input)
			if result != tt.expected {
				t.Errorf("CleanEscapes() = %q, want %q", result, tt.expected)
			}
		})
	}
}
