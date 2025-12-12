package content

import (
	"regexp"
	"strings"
)

// CleanHTML removes HTML text formatting from the input string,
// preserving the text content while removing formatting tags such as:
// - Bold (<b>, <strong>)
// - Italics (<i>, <em>)
// - Strikethrough (<s>, <strike>, <del>)
// - Underline (<u>)
// - Code (<code>, <pre>)
// - Links (<a> -> text content)
// - Images (<img> -> alt text)
// - Headers (<h1> through <h6>)
// - Lists (<ul>, <ol>, <li>)
// - Blockquotes (<blockquote>)
// - Horizontal rules (<hr>)
// - Line breaks (<br>, <br/>)
func CleanHTML(text string) string {
	result := text

	// Extract alt text from images (<img alt="text"> or <img> with alt attribute)
	imgPattern := regexp.MustCompile("(?i)<img[^>]*alt\\s*=\\s*[\"']([^\"']*)[\"'][^>]*>")
	result = imgPattern.ReplaceAllString(result, "$1")
	// Remove any remaining img tags without alt text
	imgRemainingPattern := regexp.MustCompile("(?i)<img[^>]*>")
	result = imgRemainingPattern.ReplaceAllString(result, "")

	// Remove line breaks (<br>, <br/>) and replace with newline - do this early
	brPattern := regexp.MustCompile("(?i)<br\\s*/?>")
	result = brPattern.ReplaceAllString(result, "\n")

	// Extract text content from links (<a>text</a> -> text)
	linkPattern := regexp.MustCompile("(?i)<a[^>]*>(.*?)</a>")
	result = linkPattern.ReplaceAllString(result, "$1")

	// Remove code blocks (<pre> and <code> tags, keep content)
	prePattern := regexp.MustCompile("(?i)<pre[^>]*>(.*?)</pre>")
	result = prePattern.ReplaceAllString(result, "$1")
	codePattern := regexp.MustCompile("(?i)<code[^>]*>(.*?)</code>")
	result = codePattern.ReplaceAllString(result, "$1")

	// Remove strikethrough tags (<s>, <strike>, <del>)
	strikePattern := regexp.MustCompile("(?i)</?(?:s|strike|del)[^>]*>")
	result = strikePattern.ReplaceAllString(result, "")

	// Remove underline tags (<u>)
	underlinePattern := regexp.MustCompile("(?i)</?u[^>]*>")
	result = underlinePattern.ReplaceAllString(result, "")

	// Remove bold tags (<b>, <strong>)
	boldPattern := regexp.MustCompile("(?i)</?(?:b|strong)[^>]*>")
	result = boldPattern.ReplaceAllString(result, "")

	// Remove italic tags (<i>, <em>)
	italicPattern := regexp.MustCompile("(?i)</?(?:i|em)[^>]*>")
	result = italicPattern.ReplaceAllString(result, "")

	// Remove header tags (<h1> through <h6>, keep content, add newline)
	headerPattern := regexp.MustCompile("(?i)</?h[1-6][^>]*>")
	result = headerPattern.ReplaceAllString(result, "\n")

	// Remove horizontal rule tags (<hr>, replace with newline)
	hrPattern := regexp.MustCompile("(?i)<hr[^>]*/?>")
	result = hrPattern.ReplaceAllString(result, "\n")

	// Remove blockquote tags (<blockquote>, keep content, add newline)
	blockquotePattern := regexp.MustCompile("(?i)</?blockquote[^>]*>")
	result = blockquotePattern.ReplaceAllString(result, "\n")

	// Remove list container tags (<ul>, <ol>, add newline)
	listContainerPattern := regexp.MustCompile("(?i)</?(?:ul|ol)[^>]*>")
	result = listContainerPattern.ReplaceAllString(result, "\n")
	// Remove list item tags (<li>, keep content, add newline)
	listItemPattern := regexp.MustCompile("(?i)</?li[^>]*>")
	result = listItemPattern.ReplaceAllString(result, "\n")

	// Remove paragraph tags (<p>, keep content, add spacing)
	pPattern := regexp.MustCompile("(?i)</?p[^>]*>")
	result = pPattern.ReplaceAllString(result, "\n")

	// Remove div tags (<div>, keep content)
	divPattern := regexp.MustCompile("(?i)</?div[^>]*>")
	result = divPattern.ReplaceAllString(result, "")

	// Decode HTML entities (basic ones)
	htmlEntities := map[string]string{
		"&nbsp;":  " ",
		"&amp;":   "&",
		"&lt;":    "<",
		"&gt;":    ">",
		"&quot;":  "\"",
		"&apos;":  "'",
		"&#39;":   "'",
		"&mdash;": "—",
		"&ndash;": "–",
	}
	for entity, replacement := range htmlEntities {
		result = strings.ReplaceAll(result, entity, replacement)
	}

	// Clean up multiple consecutive newlines
	newlinePattern := regexp.MustCompile("\\n{3,}")
	result = newlinePattern.ReplaceAllString(result, "\n\n")

	// Clean up multiple consecutive spaces
	spacePattern := regexp.MustCompile("[ \t]+")
	result = spacePattern.ReplaceAllString(result, " ")

	// Trim leading/trailing whitespace from each line
	lines := strings.Split(result, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	result = strings.Join(lines, "\n")

	// Trim leading/trailing whitespace
	result = strings.TrimSpace(result)

	return result
}
