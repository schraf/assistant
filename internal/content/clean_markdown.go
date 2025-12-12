package content

import (
	"regexp"
	"strings"
)

// CleanMarkdown removes markdown text formatting from the input string,
// preserving the text content while removing formatting markers such as:
// - Bold (**text** or __text__)
// - Italics (*text* or _text_)
// - Strikethrough (~~text~~)
// - Inline code (`text`)
// - Links ([text](url) -> text)
// - Images (![alt](url) -> alt)
// - Headers (entire lines starting with # are removed)
// - List markers (-, *, +, numbered)
// - Blockquotes (> text -> text)
// - Horizontal rules (---, ***)
// - Code blocks (```code``` -> code)
func CleanMarkdown(text string) string {
	result := text

	// Remove code blocks (```language\ncode\n``` or ```code```)
	codeBlockPattern := regexp.MustCompile("(?s)```[a-zA-Z]*\\n?(.*?)```")
	result = codeBlockPattern.ReplaceAllString(result, "$1")

	// Remove inline code (`code`)
	inlineCodePattern := regexp.MustCompile("`([^`]+)`")
	result = inlineCodePattern.ReplaceAllString(result, "$1")

	// Remove images (![alt](url) -> alt, ![alt][ref] -> alt)
	imagePattern := regexp.MustCompile("!\\[([^\\]]*)\\]\\([^\\)]*\\)")
	result = imagePattern.ReplaceAllString(result, "$1")
	imageRefPattern := regexp.MustCompile("!\\[([^\\]]*)\\]\\[[^\\]]*\\]")
	result = imageRefPattern.ReplaceAllString(result, "$1")

	// Remove links ([text](url) -> text, [text][ref] -> text)
	linkPattern := regexp.MustCompile("\\[([^\\]]+)\\]\\([^\\)]*\\)")
	result = linkPattern.ReplaceAllString(result, "$1")
	linkRefPattern := regexp.MustCompile("\\[([^\\]]+)\\]\\[[^\\]]*\\]")
	result = linkRefPattern.ReplaceAllString(result, "$1")

	// Remove strikethrough (~~text~~)
	strikethroughPattern := regexp.MustCompile("~~([^~]+)~~")
	result = strikethroughPattern.ReplaceAllString(result, "$1")

	// Remove bold (**text** or __text__)
	boldPattern := regexp.MustCompile("\\*\\*([^*]+)\\*\\*")
	result = boldPattern.ReplaceAllString(result, "$1")
	boldUnderscorePattern := regexp.MustCompile("__([^_]+)__")
	result = boldUnderscorePattern.ReplaceAllString(result, "$1")

	// Remove italics (*text* or _text_)
	// Match *text* but not **text** (bold) - process after bold removal
	// This pattern matches single asterisks that aren't part of double asterisks
	italicPattern := regexp.MustCompile("\\*([^*\\n]+)\\*")
	result = italicPattern.ReplaceAllString(result, "$1")
	// Match _text_ but not __text__ (bold) - process after bold removal
	italicUnderscorePattern := regexp.MustCompile("_([^_\\n]+)_")
	result = italicUnderscorePattern.ReplaceAllString(result, "$1")

	// Remove headers (entire lines starting with #)
	headerPattern := regexp.MustCompile("(?m)^#{1,}.*$")
	result = headerPattern.ReplaceAllString(result, "")

	// Remove horizontal rules (---, ***, ___)
	horizontalRulePattern := regexp.MustCompile("(?m)^[-*_]{3,}\\s*$")
	result = horizontalRulePattern.ReplaceAllString(result, "")

	// Remove blockquotes (> text -> text)
	blockquotePattern := regexp.MustCompile("(?m)^>\\s+(.+)$")
	result = blockquotePattern.ReplaceAllString(result, "$1")

	// Remove list markers (-, *, +, numbered)
	listPattern := regexp.MustCompile("(?m)^([\\s]*)[-*+]\\s+(.+)$")
	result = listPattern.ReplaceAllString(result, "$1$2")
	numberedListPattern := regexp.MustCompile("(?m)^([\\s]*)\\d+\\.\\s+(.+)$")
	result = numberedListPattern.ReplaceAllString(result, "$1$2")

	// Clean up multiple consecutive newlines
	newlinePattern := regexp.MustCompile("\\n{3,}")
	result = newlinePattern.ReplaceAllString(result, "\n\n")

	// Trim leading/trailing whitespace
	result = strings.TrimSpace(result)

	return result
}
