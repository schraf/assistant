package content

import (
	"regexp"
	"strings"
)

// CleanLaTeX removes LaTeX text formatting from the input string,
// preserving the text content while removing formatting commands such as:
// - Bold (\textbf{text})
// - Italics (\textit{text}, \emph{text})
// - Strikethrough (\sout{text})
// - Underline (\underline{text})
// - Code (\texttt{text}, \verb|text|)
// - Links (\href{url}{text}, \url{url})
// - Headers (\section{title}, \subsection{title}, etc.)
// - Lists (\item in itemize/enumerate environments)
// - Blockquotes (\begin{quote}...\end{quote})
// - Math mode ($...$, \(...\), \[...\])
// - Comments (% comment)
func CleanLaTeX(text string) string {
	result := text

	// Remove comments (% comment, but not \% escaped percent)
	// Match % and everything after it on the line, but keep text before %
	commentPattern := regexp.MustCompile("(?m)([^%\\\\]*(?:\\\\.[^%\\\\]*)*)%[^\n]*")
	result = commentPattern.ReplaceAllString(result, "$1")

	// Remove math mode (inline and display)
	// Use a placeholder that will be replaced with double space after space cleanup
	mathPlaceholder := "___MATH_PLACEHOLDER___"
	// Inline math: $...$ or \(...\)
	inlineMathPattern := regexp.MustCompile("\\$\\(?[^$)]+\\)?\\$")
	result = inlineMathPattern.ReplaceAllString(result, mathPlaceholder)
	// Display math: \[...\] or $$...$$
	displayMathPattern := regexp.MustCompile("\\\\\\[.*?\\\\\\]|\\$\\$.*?\\$\\$")
	result = displayMathPattern.ReplaceAllString(result, mathPlaceholder)

	// Extract text from href links (\href{url}{text} -> text)
	hrefPattern := regexp.MustCompile("\\\\href\\{[^}]*\\}\\{([^}]+)\\}")
	result = hrefPattern.ReplaceAllString(result, "$1")
	// Extract text from url (\url{url} -> url)
	urlPattern := regexp.MustCompile("\\\\url\\{([^}]+)\\}")
	result = urlPattern.ReplaceAllString(result, "$1")

	// Remove verbatim text (\verb|text| or \verb*|text|)
	verbPattern := regexp.MustCompile("\\\\verb\\*?[|!@#$%^&*+=\\-_~`'\":;<>?/.,\\[\\]{}()\\\\]+([^|!@#$%^&*+=\\-_~`'\":;<>?/.,\\[\\]{}()\\\\]+)[|!@#$%^&*+=\\-_~`'\":;<>?/.,\\[\\]{}()\\\\]+")
	result = verbPattern.ReplaceAllString(result, "$1")

	// Remove code formatting (\texttt{text} -> text)
	textttPattern := regexp.MustCompile("\\\\texttt\\{([^}]+)\\}")
	result = textttPattern.ReplaceAllString(result, "$1")

	// Remove strikethrough (\sout{text} -> text)
	soutPattern := regexp.MustCompile("\\\\sout\\{([^}]+)\\}")
	result = soutPattern.ReplaceAllString(result, "$1")

	// Remove underline (\underline{text} -> text)
	underlinePattern := regexp.MustCompile("\\\\underline\\{([^}]+)\\}")
	result = underlinePattern.ReplaceAllString(result, "$1")

	// Remove bold (\textbf{text} -> text)
	textbfPattern := regexp.MustCompile("\\\\textbf\\{([^}]+)\\}")
	result = textbfPattern.ReplaceAllString(result, "$1")
	// Remove bold series (\bfseries, but this affects following text, so just remove the command)
	bfseriesPattern := regexp.MustCompile("\\\\bfseries\\b")
	result = bfseriesPattern.ReplaceAllString(result, "")

	// Remove italics (\textit{text} -> text)
	textitPattern := regexp.MustCompile("\\\\textit\\{([^}]+)\\}")
	result = textitPattern.ReplaceAllString(result, "$1")
	// Remove emphasis (\emph{text} -> text)
	emphPattern := regexp.MustCompile("\\\\emph\\{([^}]+)\\}")
	result = emphPattern.ReplaceAllString(result, "$1")
	// Remove italic shape (\itshape)
	itshapePattern := regexp.MustCompile("\\\\itshape\\b")
	result = itshapePattern.ReplaceAllString(result, "")

	// Remove sectioning commands (\section{title} -> title, add newline)
	sectionPattern := regexp.MustCompile("\\\\section\\*?\\{([^}]+)\\}")
	result = sectionPattern.ReplaceAllString(result, "$1\n\n")
	subsectionPattern := regexp.MustCompile("\\\\subsection\\*?\\{([^}]+)\\}")
	result = subsectionPattern.ReplaceAllString(result, "$1\n")
	subsubsectionPattern := regexp.MustCompile("\\\\subsubsection\\*?\\{([^}]+)\\}")
	result = subsubsectionPattern.ReplaceAllString(result, "$1\n")
	paragraphPattern := regexp.MustCompile("\\\\paragraph\\*?\\{([^}]+)\\}")
	result = paragraphPattern.ReplaceAllString(result, "$1\n")
	subparagraphPattern := regexp.MustCompile("\\\\subparagraph\\*?\\{([^}]+)\\}")
	result = subparagraphPattern.ReplaceAllString(result, "$1\n")
	chapterPattern := regexp.MustCompile("\\\\chapter\\*?\\{([^}]+)\\}")
	result = chapterPattern.ReplaceAllString(result, "$1\n")
	partPattern := regexp.MustCompile("\\\\part\\*?\\{([^}]+)\\}")
	result = partPattern.ReplaceAllString(result, "$1\n")

	// Remove quote environments (\begin{quote}...\end{quote} -> content)
	quotePattern := regexp.MustCompile("(?s)\\\\begin\\{quote\\}(.*?)\\\\end\\{quote\\}")
	result = quotePattern.ReplaceAllString(result, "$1")
	quotationPattern := regexp.MustCompile("(?s)\\\\begin\\{quotation\\}(.*?)\\\\end\\{quotation\\}")
	result = quotationPattern.ReplaceAllString(result, "$1")

	// Remove list environments and items
	// Extract content from itemize/enumerate environments
	itemizePattern := regexp.MustCompile("(?s)\\\\begin\\{(?:itemize|enumerate)\\}(.*?)\\\\end\\{(?:itemize|enumerate)\\}")
	result = itemizePattern.ReplaceAllString(result, "$1")
	// Remove \item commands, keeping the content, add newline
	// Match \item with optional bracket argument, then capture content (either in braces or until next \item or end)
	itemPattern := regexp.MustCompile("\\\\item\\s*(?:\\[[^\\]]*\\])?\\s*(?:\\{([^}]+)\\}|([^\\\\{]+))")
	result = itemPattern.ReplaceAllStringFunc(result, func(match string) string {
		// Extract the captured group (either from braces or plain text)
		if matches := itemPattern.FindStringSubmatch(match); len(matches) > 1 {
			content := matches[1]
			if content == "" && len(matches) > 2 {
				content = matches[2]
			}
			return strings.TrimSpace(content) + "\n"
		}
		return "\n"
	})

	// Remove horizontal rules (\hrule, \rule)
	hrPattern := regexp.MustCompile("\\\\hrule\\b|\\\\rule\\{[^}]*\\}\\{[^}]*\\}")
	result = hrPattern.ReplaceAllString(result, "")

	// Remove common environments (keep content)
	centerPattern := regexp.MustCompile("(?s)\\\\begin\\{center\\}(.*?)\\\\end\\{center\\}")
	result = centerPattern.ReplaceAllString(result, "$1")
	flushleftPattern := regexp.MustCompile("(?s)\\\\begin\\{flushleft\\}(.*?)\\\\end\\{flushleft\\}")
	result = flushleftPattern.ReplaceAllString(result, "$1")
	flushrightPattern := regexp.MustCompile("(?s)\\\\begin\\{flushright\\}(.*?)\\\\end\\{flushright\\}")
	result = flushrightPattern.ReplaceAllString(result, "$1")

	// Remove common LaTeX commands that don't affect text content
	commonCommands := []string{
		"\\newpage", "\\clearpage", "\\pagebreak",
		"\\vspace", "\\hspace", "\\vfill", "\\hfill",
		"\\label", "\\ref", "\\cite", "\\footnote",
		"\\maketitle", "\\tableofcontents",
	}
	for _, cmd := range commonCommands {
		cmdPattern := regexp.MustCompile(regexp.QuoteMeta(cmd) + "\\{[^}]*\\}|" + regexp.QuoteMeta(cmd) + "\\b")
		result = cmdPattern.ReplaceAllString(result, "")
	}

	// Remove remaining backslash commands (basic ones, be careful not to break everything)
	// This is a catch-all for simple commands like \LaTeX, \today, etc.
	simpleCommandPattern := regexp.MustCompile("\\\\[a-zA-Z]+\\*?\\b")
	result = simpleCommandPattern.ReplaceAllString(result, "")

	// Replace math placeholder with single space (before space cleanup)
	result = strings.ReplaceAll(result, "___MATH_PLACEHOLDER___", " ")

	// Clean up multiple consecutive newlines
	newlinePattern := regexp.MustCompile("\\n{3,}")
	result = newlinePattern.ReplaceAllString(result, "\n\n")

	// Clean up multiple consecutive spaces
	// Collapse 3+ spaces to single space, leave 1-2 spaces as is
	spacePattern := regexp.MustCompile("[ \t]{3,}")
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
