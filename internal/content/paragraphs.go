package content

import (
	"regexp"
	"strings"
)

var splitParagraphPattern *regexp.Regexp
var sentenceEndPattern *regexp.Regexp

func init() {
	splitParagraphPattern = regexp.MustCompile(`\n\s*\n`)
	sentenceEndPattern = regexp.MustCompile(`[.!?]`)
}

// hasEndingPunctuation checks if a paragraph ends with sentence-ending punctuation
func hasEndingPunctuation(paragraph string) bool {
	if len(paragraph) == 0 {
		return false
	}
	lastChar := paragraph[len(paragraph)-1]
	return lastChar == '.' || lastChar == '!' || lastChar == '?'
}

// isSingleSentence checks if a paragraph contains only a single sentence
func isSingleSentence(paragraph string) bool {
	matches := sentenceEndPattern.FindAllString(paragraph, -1)
	return len(matches) <= 1
}

func SplitParagraphs(content string) []string {
	paragraphs := splitParagraphPattern.Split(content, -1)
	filtered := make([]string, 0)

	for _, paragraph := range paragraphs {
		paragraph = strings.TrimSpace(paragraph)
		
		// Skip empty paragraphs
		if paragraph == "" {
			continue
		}
		
		// Skip paragraphs that are single sentences or lack ending punctuation
		if isSingleSentence(paragraph) || !hasEndingPunctuation(paragraph) {
			continue
		}
		
		filtered = append(filtered, paragraph)
	}

	return filtered
}
