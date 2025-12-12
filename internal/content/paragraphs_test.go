package content

import (
	"reflect"
	"testing"
)

func TestSplitParagraphs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "single empty paragraph",
			input:    "\n\n",
			expected: []string{},
		},
		{
			name:     "whitespace only",
			input:    "   \n\n   ",
			expected: []string{},
		},
		{
			name:     "single sentence paragraph - should be filtered",
			input:    "This is a single sentence.",
			expected: []string{},
		},
		{
			name:     "single sentence without punctuation - should be filtered",
			input:    "This is a single sentence without punctuation",
			expected: []string{},
		},
		{
			name:     "multiple sentences paragraph - should be kept",
			input:    "This is the first sentence. This is the second sentence. This is the third sentence.",
			expected: []string{"This is the first sentence. This is the second sentence. This is the third sentence."},
		},
		{
			name:     "paragraph with exclamation - multiple sentences",
			input:    "This is exciting! This is also exciting! This is the end.",
			expected: []string{"This is exciting! This is also exciting! This is the end."},
		},
		{
			name:     "paragraph with question - multiple sentences",
			input:    "What is this? This is a test. This is the answer.",
			expected: []string{"What is this? This is a test. This is the answer."},
		},
		{
			name: "multiple paragraphs - filter single sentences",
			input: `This is a single sentence paragraph.

This is the first sentence of a multi-sentence paragraph. This is the second sentence. This is the third sentence.

Another single sentence.

First sentence here. Second sentence here. Third sentence here.`,
			expected: []string{
				"This is the first sentence of a multi-sentence paragraph. This is the second sentence. This is the third sentence.",
				"First sentence here. Second sentence here. Third sentence here.",
			},
		},
		{
			name: "paragraphs without ending punctuation - should be filtered",
			input: `This paragraph has no ending punctuation

This paragraph has ending punctuation. This is the second sentence.`,
			expected: []string{
				"This paragraph has ending punctuation. This is the second sentence.",
			},
		},
		{
			name: "mixed valid and invalid paragraphs",
			input: `Single sentence.

No ending punctuation here

Valid paragraph. Second sentence. Third sentence.

Another valid one! Second sentence! Third sentence!`,
			expected: []string{
				"Valid paragraph. Second sentence. Third sentence.",
				"Another valid one! Second sentence! Third sentence!",
			},
		},
		{
			name: "paragraphs with extra whitespace",
			input: `   This is a single sentence.   

   First sentence.   Second sentence.   Third sentence.   `,
			expected: []string{
				"First sentence.   Second sentence.   Third sentence.",
			},
		},
		{
			name: "paragraph with only question mark - single sentence",
			input: "What is this?",
			expected: []string{},
		},
		{
			name: "paragraph with only exclamation - single sentence",
			input: "Wow!",
			expected: []string{},
		},
		{
			name: "complex paragraph with abbreviations",
			input: "Dr. Smith went to the U.S.A. This is a second sentence. This is a third sentence.",
			expected: []string{"Dr. Smith went to the U.S.A. This is a second sentence. This is a third sentence."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SplitParagraphs(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SplitParagraphs() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHasEndingPunctuation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "ends with period",
			input:    "This ends with a period.",
			expected: true,
		},
		{
			name:     "ends with exclamation",
			input:    "This ends with an exclamation!",
			expected: true,
		},
		{
			name:     "ends with question mark",
			input:    "This ends with a question?",
			expected: true,
		},
		{
			name:     "no ending punctuation",
			input:    "This has no ending punctuation",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "ends with comma",
			input:    "This ends with a comma,",
			expected: false,
		},
		{
			name:     "ends with semicolon",
			input:    "This ends with a semicolon;",
			expected: false,
		},
		{
			name:     "period in middle",
			input:    "This has a period. But no ending punctuation",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasEndingPunctuation(tt.input)
			if result != tt.expected {
				t.Errorf("hasEndingPunctuation(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsSingleSentence(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "single sentence with period",
			input:    "This is a single sentence.",
			expected: true,
		},
		{
			name:     "single sentence with exclamation",
			input:    "This is exciting!",
			expected: true,
		},
		{
			name:     "single sentence with question",
			input:    "What is this?",
			expected: true,
		},
		{
			name:     "no punctuation",
			input:    "This has no punctuation",
			expected: true,
		},
		{
			name:     "two sentences",
			input:    "First sentence. Second sentence.",
			expected: false,
		},
		{
			name:     "three sentences",
			input:    "First sentence. Second sentence. Third sentence.",
			expected: false,
		},
		{
			name:     "multiple exclamations",
			input:    "First! Second! Third!",
			expected: false,
		},
		{
			name:     "mixed punctuation",
			input:    "First sentence? Second sentence! Third sentence.",
			expected: false,
		},
		{
			name:     "abbreviations count as punctuation",
			input:    "Dr. Smith went to the U.S.A.",
			expected: false, // Multiple periods are detected, so it's not considered a single sentence
		},
		{
			name:     "abbreviations with multiple sentences",
			input:    "Dr. Smith went to the U.S.A. This is a second sentence.",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSingleSentence(tt.input)
			if result != tt.expected {
				t.Errorf("isSingleSentence(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
