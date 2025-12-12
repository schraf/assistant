package content

// CleanText applies all cleaning functions to the input text in sequence:
// 1. CleanMarkdown - removes markdown formatting
// 2. CleanHTML - removes HTML formatting
// 3. CleanLaTeX - removes LaTeX formatting (must be before CleanEscapes to preserve backslashes)
// 4. CleanEscapes - removes escape characters (applied last to avoid breaking LaTeX commands)
//
// This provides comprehensive text cleaning by removing formatting from
// multiple markup languages and escape sequences.
func CleanText(text string) string {
	result := text

	// Remove markdown formatting
	result = CleanMarkdown(result)

	// Remove HTML formatting
	result = CleanHTML(result)

	// Remove LaTeX formatting (before CleanEscapes to preserve backslashes)
	result = CleanLaTeX(result)

	// Remove escape characters last (after LaTeX which uses backslashes)
	result = CleanEscapes(result)

	return result
}
