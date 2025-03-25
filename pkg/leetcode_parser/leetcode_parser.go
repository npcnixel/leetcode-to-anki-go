package leetcode_parser

import (
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	nethtml "golang.org/x/net/html"
)

// Problem represents a parsed LeetCode problem
type Problem struct {
	Title       string
	Description string
	Code        string
	Filename    string
	Timestamp   time.Time
}

// Controls debug logging throughout the package
var debugLogging bool

// Parse a directory of HTML files
func ParseDirectory(inputDir string, debug bool) ([]*Problem, error) {
	// Set debug logging flag for the entire package
	debugLogging = debug

	var problems []*Problem
	err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			problem, err := ParseHTMLFile(path, debug)
			if err != nil {
				return err
			}
			problems = append(problems, problem)
		}
		return nil
	})
	return problems, err
}

// Parse a single HTML file
func ParseHTMLFile(filePath string, debug bool) (*Problem, error) {
	if debug {
		fmt.Printf("Parsing file: %s\n", filePath)
	}

	// Read file
	htmlContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var problem Problem

	// Extract title - try data-next-head title tag first (newer format)
	titleRegex := regexp.MustCompile(`<title data-next-head[^>]*>(.*?) - LeetCode</title>`)
	titleMatch := titleRegex.FindSubmatch(htmlContent)
	if len(titleMatch) > 1 {
		problem.Title = string(titleMatch[1])
		if debug {
			fmt.Printf("Extracted title from data-next-head: %s\n", problem.Title)
		}
	} else {
		// Fall back to regular title tag (older format)
		oldTitleRegex := regexp.MustCompile(`<title>(.*?) - LeetCode</title>`)
		oldTitleMatch := oldTitleRegex.FindSubmatch(htmlContent)
		if len(oldTitleMatch) > 1 {
			problem.Title = string(oldTitleMatch[1])
			if debug {
				fmt.Printf("Extracted title from regular title tag: %s\n", problem.Title)
			}
		}
	}

	// Extract description using meta tag first
	descriptionRegex := regexp.MustCompile(`<meta\s+name="description"\s+content="([^"]*)"`)
	descriptionMatch := descriptionRegex.FindSubmatch(htmlContent)
	if len(descriptionMatch) > 1 {
		problem.Description = html.UnescapeString(string(descriptionMatch[1]))
		if debug {
			fmt.Printf("Meta description length: %d characters\n", len(problem.Description))
		}
	}

	// Try to extract a more detailed description
	detailedDesc := ExtractProblemDescription(string(htmlContent), debug)
	if detailedDesc != "" && len(detailedDesc) > len(problem.Description) {
		problem.Description = detailedDesc
		if debug {
			fmt.Printf("Enhanced description length: %d characters\n", len(problem.Description))
		}
	}

	// Extract Python code
	problem.Code = extractPythonCode(string(htmlContent))

	if debug {
		fmt.Printf("Final formatted code:\n%s\n", problem.Code)
	}

	// Store filename
	problem.Filename = filepath.Base(filePath)

	// Parse timestamp from filename if present
	pattern := regexp.MustCompile(`\((\d{2}_\d{2}_\d{4} \d{2}[：:]\d{2}[：:]\d{2})\)`)
	match := pattern.FindStringSubmatch(problem.Filename)
	if len(match) > 1 {
		timeStr := match[1]
		// Replace Chinese colons with standard ones if needed
		timeStr = strings.ReplaceAll(timeStr, "：", ":")

		// Attempt to parse the time
		t, err := time.Parse("02_01_2006 15:04:05", timeStr)
		if err == nil {
			problem.Timestamp = t
		}
	}

	return &problem, nil
}

// Format the problem description with HTML styling
func FormatDescription(desc string) string {
	lines := strings.Split(desc, "\n")
	var formatted []string
	inExample := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		line = strings.ReplaceAll(line, "Can you solve this real interview question? ", "")
		parts := strings.Split(line, " - ")
		if len(parts) > 1 {
			line = strings.TrimSpace(parts[1])
		}

		if strings.HasPrefix(line, "Example") {
			if !inExample {
				formatted = append(formatted, "<div class='example'>")
				inExample = true
			} else {
				formatted = append(formatted, "</div><div class='example'>")
			}
			formatted = append(formatted, "<h3>"+line+"</h3>")
		} else if strings.HasPrefix(line, "Input:") || strings.HasPrefix(line, "Output:") {
			formatted = append(formatted, "<pre>"+line+"</pre>")
		} else if strings.HasPrefix(line, "Explanation:") {
			formatted = append(formatted, "<p><strong>"+line+"</strong></p>")
		} else if strings.HasPrefix(line, "Constraints:") {
			if inExample {
				formatted = append(formatted, "</div>")
				inExample = false
			}
			formatted = append(formatted, "<h3 style='font-size: 16px; color: #ff9800;'>"+line+"</h3>")
		} else if strings.HasPrefix(line, " * ") {
			formatted = append(formatted, "<li>"+strings.TrimPrefix(line, " * ")+"</li>")
		} else {
			if strings.Contains(line, "<=") {
				formatted = append(formatted, "<li>"+line+"</li>")
			} else {
				formatted = append(formatted, "<p>"+line+"</p>")
			}
		}
	}

	if inExample {
		formatted = append(formatted, "</div>")
	}

	// Add some CSS for styling, using dark theme colors and normalized font sizes
	css := `
<style>
body {
    font-size: 14px;
    line-height: 1.5;
}
.example {
    background-color: #3d3d3d;
    border-left: 3px solid #ff9800;
    padding: 10px;
    margin: 10px 0;
}
.example h3 {
    margin: 0 0 10px 0;
    color: #ff9800;
    font-size: 16px;
}
.example pre {
    margin: 5px 0;
    font-family: monospace;
    background-color: #333;
    padding: 5px;
    color: #e0e0e0;
    font-size: 14px;
}
li {
    margin: 5px 0;
    list-style-type: none;
    color: #e0e0e0;
    font-size: 14px;
}
p {
    font-size: 14px;
    margin: 8px 0;
}
</style>
`

	return css + strings.Join(formatted, "\n")
}

// Extract Python code from the HTML content
func extractPythonCode(htmlContent string) string {
	debugLog := func(format string, args ...interface{}) {
		if debugLogging {
			log.Printf(format, args...)
		}
	}

	debugLog("Extracting Python code from HTML content...")

	// Look for Python code in a textarea - using a pattern that targets actual code
	codePattern := `<textarea[^>]*>\s*(class\s+Solution[\s\S]*?return\s+\w+[\s\S]*?)</textarea>`
	codeRegex := regexp.MustCompile(codePattern)
	codeMatches := codeRegex.FindStringSubmatch(htmlContent)

	if len(codeMatches) > 1 {
		debugLog("Found Python class and method in textarea")
		codeContent := codeMatches[1]

		// Split by lines and clean up
		lines := strings.Split(codeContent, "\n")
		var cleanedLines []string
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" {
				cleanedLines = append(cleanedLines, trimmed)
			}
		}

		if len(cleanedLines) > 0 {
			debugLog("Extracted %d lines of clean Python code", len(cleanedLines))
			return strings.Join(cleanedLines, "\n")
		}
	}

	// Extract general textarea content, but filter to only include lines that look like code
	textareaPattern := `<textarea[^>]*>([\s\S]*?)</textarea>`
	textareaRegex := regexp.MustCompile(textareaPattern)
	textareaMatches := textareaRegex.FindStringSubmatch(htmlContent)

	if len(textareaMatches) > 1 {
		debugLog("Found textarea content, filtering for code")
		textareaContent := textareaMatches[1]

		// Split by lines and clean up
		lines := strings.Split(textareaContent, "\n")
		var cleanedLines []string

		// First pass: identify if this looks like Python code at all
		hasPythonCode := false
		for _, line := range lines {
			if strings.Contains(line, "class Solution") ||
				strings.Contains(line, "def ") && strings.Contains(line, "(self") {
				hasPythonCode = true
				break
			}
		}

		if !hasPythonCode {
			debugLog("No Python code found in textarea, trying other methods")
		} else {
			// Second pass: extract only code lines
			for _, line := range lines {
				trimmed := strings.TrimSpace(line)
				if trimmed != "" && isCodeLine(trimmed) {
					cleanedLines = append(cleanedLines, trimmed)
				}
			}

			if len(cleanedLines) > 0 {
				debugLog("Extracted %d lines of filtered code from textarea", len(cleanedLines))
				return strings.Join(cleanedLines, "\n")
			}
		}
	}

	// Try extracting from Monaco editor view-lines with a more inclusive pattern
	monacoPattern := `<div class="view-lines monaco-mouse-cursor-text"[^>]*>([\s\S]*?)</div><div data-mprt=1 class=contentWidgets`
	monacoRegex := regexp.MustCompile(monacoPattern)
	monacoMatches := monacoRegex.FindStringSubmatch(htmlContent)

	if len(monacoMatches) > 1 {
		debugLog("Found Monaco editor content section")
		viewLinesContent := monacoMatches[1]

		// Extract each view-line div with its top position for proper ordering
		// Make the pattern more tolerant of variations in HTML structure
		viewLinePattern := `<div style=top:(\d+)px[^>]*class=view-line><span>([\s\S]*?)</span></div>`
		viewLineRegex := regexp.MustCompile(viewLinePattern)
		viewLineMatches := viewLineRegex.FindAllStringSubmatch(viewLinesContent, -1)

		if len(viewLineMatches) > 0 {
			debugLog("Found %d view-line divs", len(viewLineMatches))

			// Create a map to store lines by their vertical position
			linesByPosition := make(map[int]string)

			for _, match := range viewLineMatches {
				if len(match) > 2 {
					position, err := strconv.Atoi(match[1])
					if err == nil {
						lineContent := match[2]
						cleanedLine := cleanCodeLine(lineContent)
						// Skip empty lines
						if strings.TrimSpace(cleanedLine) != "" {
							linesByPosition[position] = cleanedLine
						}
					}
				}
			}

			// Get all positions and sort them
			var positions []int
			for pos := range linesByPosition {
				positions = append(positions, pos)
			}
			sort.Ints(positions)

			// Build the code in the correct order
			var codeLines []string
			for _, pos := range positions {
				codeLines = append(codeLines, linesByPosition[pos])
			}

			// Check if we have maxProfit and are missing return total_profit
			if len(codeLines) > 0 {
				code := strings.Join(codeLines, "\n")
				if strings.Contains(code, "def ") && !strings.Contains(code, "return ") {
					// Look specifically for any return line
					returnLinePattern := `<div[^>]*>.*?return\s+\w+.*?</div>`
					returnLineRegex := regexp.MustCompile(returnLinePattern)
					returnLineMatch := returnLineRegex.FindString(viewLinesContent)

					if returnLineMatch != "" {
						debugLog("Found missing return statement")
						cleanedReturn := cleanCodeLine(returnLineMatch)
						if cleanedReturn != "" {
							codeLines = append(codeLines, cleanedReturn)
						}
					}
				}

				// Now we need to reorder the code for logical flow
				return reorganizePythonCode(codeLines)
			}
		}

		// If the regular extraction failed, try a more aggressive approach
		debugLog("Standard extraction failed, trying aggressive pattern matching...")

		// Try to find all mtk spans directly
		mtkPattern := `<span class=mtk[0-9]+>([\s\S]*?)</span>`
		mtkRegex := regexp.MustCompile(mtkPattern)
		mtkMatches := mtkRegex.FindAllStringSubmatch(viewLinesContent, -1)

		if len(mtkMatches) > 0 {
			debugLog("Found %d mtk spans", len(mtkMatches))

			var rawCode strings.Builder
			for _, match := range mtkMatches {
				if len(match) > 1 {
					text := html.UnescapeString(match[1])
					text = strings.ReplaceAll(text, "\u00A0", " ")
					rawCode.WriteString(text)
				}
			}

			// Process the raw code to split into lines
			rawCodeStr := rawCode.String()
			splitLines := strings.Split(rawCodeStr, "\n")

			// Clean up lines
			var cleanLines []string
			currentLine := ""
			for _, token := range splitLines {
				if strings.Contains(token, "def ") || strings.Contains(token, "class ") ||
					strings.Contains(token, "return ") || strings.Contains(token, "if ") ||
					strings.Contains(token, "for ") || strings.Contains(token, "while ") {
					if currentLine != "" {
						cleanLines = append(cleanLines, currentLine)
					}
					currentLine = token
				} else {
					currentLine += token
				}
			}
			if currentLine != "" {
				cleanLines = append(cleanLines, currentLine)
			}

			if len(cleanLines) > 0 {
				debugLog("Successfully extracted code with aggressive pattern matching")
				return reorganizePythonCode(cleanLines)
			}
		}

		debugLog("Failed to extract code from view-lines, trying alternative methods...")
	} else {
		debugLog("Monaco editor content section not found, trying alternative methods...")
	}

	// Fallback to span-based method
	debugLog("Trying span-based method...")
	spanPattern := `<span class="token[^>]*>(.*?)</span>`
	spanRegex := regexp.MustCompile(spanPattern)
	spanMatches := spanRegex.FindAllStringSubmatch(htmlContent, -1)

	if len(spanMatches) > 0 {
		var code strings.Builder
		for _, match := range spanMatches {
			if len(match) > 1 {
				code.WriteString(html.UnescapeString(match[1]))
			}
		}

		if code.Len() > 0 {
			debugLog("Successfully extracted code using span-based method")
			return code.String()
		}
	}

	// Fallback to pre tag method
	debugLog("Trying pre tag method...")
	prePattern := `<pre>(.*?)</pre>`
	preRegex := regexp.MustCompile(prePattern)
	preMatches := preRegex.FindAllStringSubmatch(htmlContent, -1)

	if len(preMatches) > 0 {
		for _, match := range preMatches {
			if len(match) > 1 && strings.Contains(match[1], "def") {
				debugLog("Successfully extracted code using pre tag method")
				return html.UnescapeString(match[1])
			}
		}
	}

	// Return a default Python template if all extraction methods fail
	debugLog("Failed to extract code using all methods, returning default template")
	return "If you see this, let me know at https://github.com/npcnixel/leetcode-to-anki-go/issues/new"
}

// Helper function to clean up code lines extracted from Monaco editor
func cleanCodeLine(line string) string {
	// Remove HTML tags except for content inside span tags with class "mtk"
	spanPattern := `<span class=mtk[0-9]+>(.*?)</span>`
	spanRegex := regexp.MustCompile(spanPattern)
	matches := spanRegex.FindAllStringSubmatch(line, -1)

	var result strings.Builder
	for _, match := range matches {
		if len(match) > 1 {
			// Convert HTML entities to their actual characters
			text := html.UnescapeString(match[1])
			// Replace &nbsp; with actual spaces
			text = strings.ReplaceAll(text, "\u00A0", " ")
			result.WriteString(text)
		}
	}

	return result.String()
}

// Helper function to extract attributes from HTML nodes
func getAttr(n *nethtml.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

// Enhance extractProblemDescription to look for specific content patterns
func ExtractProblemDescription(htmlContent string, debug bool) string {
	// Try a direct regex approach first, looking for common patterns in LeetCode problems
	problemDescRegex := regexp.MustCompile(`(?s)<div class="content__[^"]*?">.*?<div class="question-content__[^"]*?">(.*?)<\/div>.*?<div class="css-isal7m">`)
	descMatch := problemDescRegex.FindStringSubmatch(htmlContent)
	if len(descMatch) > 1 {
		description := descMatch[1]
		// Remove HTML tags but preserve line breaks for formatting
		description = regexp.MustCompile(`<br\s*/?>|<p>|</p>`).ReplaceAllString(description, "\n")
		description = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(description, " ") // Replace tags with space to avoid word joining

		// Replace HTML entities
		description = strings.ReplaceAll(description, "&nbsp;", " ")
		description = strings.ReplaceAll(description, "&lt;", "<")
		description = strings.ReplaceAll(description, "&gt;", ">")
		description = strings.ReplaceAll(description, "&amp;", "&")
		description = strings.ReplaceAll(description, "&quot;", "\"")

		// Clean up the text
		description = regexp.MustCompile(`\s+`).ReplaceAllString(description, " ")
		description = regexp.MustCompile(`\n\s*\n`).ReplaceAllString(description, "\n\n")
		description = strings.TrimSpace(description)

		if debug {
			fmt.Println("Description extracted using problem content regex:")
			fmt.Println(description)
		}

		return description
	}

	// Try another approach with different class names
	altProblemDescRegex := regexp.MustCompile(`(?s)<div class="description__[^"]*?">(.*?)<\/div>.*?<div class="editor__[^"]*?">`)
	altDescMatch := altProblemDescRegex.FindStringSubmatch(htmlContent)
	if len(altDescMatch) > 1 {
		description := altDescMatch[1]
		// Clean up as above
		description = regexp.MustCompile(`<br\s*/?>|<p>|</p>`).ReplaceAllString(description, "\n")
		description = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(description, " ")
		description = strings.ReplaceAll(description, "&nbsp;", " ")
		description = strings.ReplaceAll(description, "&lt;", "<")
		description = strings.ReplaceAll(description, "&gt;", ">")
		description = strings.ReplaceAll(description, "&amp;", "&")
		description = strings.ReplaceAll(description, "&quot;", "\"")
		description = regexp.MustCompile(`\s+`).ReplaceAllString(description, " ")
		description = regexp.MustCompile(`\n\s*\n`).ReplaceAllString(description, "\n\n")
		description = strings.TrimSpace(description)

		if debug {
			fmt.Println("Description extracted using alternate content regex:")
			fmt.Println(description)
		}

		return description
	}

	// If the regex approaches failed, fall back to DOM traversal
	doc, err := nethtml.Parse(strings.NewReader(htmlContent))
	if err != nil {
		if debug {
			fmt.Println("Error parsing HTML:", err)
		}
		return ""
	}

	var description string
	var findDescriptionDiv func(*nethtml.Node)
	findDescriptionDiv = func(n *nethtml.Node) {
		if description != "" {
			return // Already found
		}

		if n.Type == nethtml.ElementNode && n.Data == "div" {
			// Check various classes that might contain description
			className := getAttr(n, "class")
			if strings.Contains(className, "question-content") ||
				strings.Contains(className, "description") ||
				strings.Contains(getAttr(n, "data-track-load"), "description_content") {
				description = extractTextFromNode(n)
				return
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findDescriptionDiv(c)
		}
	}

	findDescriptionDiv(doc)

	if debug && description != "" {
		fmt.Println("Description extracted from HTML traversal:")
		fmt.Println(description)
	}

	return description
}

// Helper function to extract text content from HTML node and its children
func extractTextFromNode(n *nethtml.Node) string {
	var text strings.Builder

	var traverse func(*nethtml.Node)
	traverse = func(node *nethtml.Node) {
		if node.Type == nethtml.TextNode {
			text.WriteString(node.Data)
		} else if node.Type == nethtml.ElementNode {
			if node.Data == "br" || node.Data == "p" || node.Data == "div" {
				text.WriteString("\n")
			}
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}

		if node.Type == nethtml.ElementNode && (node.Data == "p" || node.Data == "div") {
			text.WriteString("\n")
		}
	}

	traverse(n)

	// Clean up the text
	content := text.String()
	content = regexp.MustCompile(`\n\s*\n`).ReplaceAllString(content, "\n\n")
	content = strings.TrimSpace(content)

	return content
}

// Function to reorganize Python code for better flow
func reorganizePythonCode(lines []string) string {
	debugLog := func(format string, args ...interface{}) {
		if debugLogging {
			log.Printf(format, args...)
		}
	}

	debugLog("Reorganizing Python code for logical flow...")

	// Find class definition lines
	var classLines []string
	var methodLines []string
	var bodyLines []string

	// Define patterns to identify different parts of the code
	classPattern := regexp.MustCompile(`^\s*class\s+Solution`)
	methodPattern := regexp.MustCompile(`^\s*def\s+\w+\s*\(self`)

	classFound := false
	methodFound := false

	for _, line := range lines {
		if classPattern.MatchString(line) {
			classFound = true
			classLines = append(classLines, line)
		} else if methodPattern.MatchString(line) {
			methodFound = true
			methodLines = append(methodLines, line)
		} else {
			// If line doesn't match class or method pattern, it's part of method body
			bodyLines = append(bodyLines, line)
		}
	}

	// If no class or method found, return original lines
	if !classFound || !methodFound {
		debugLog("Could not identify class or method structure, returning original code")
		return strings.Join(lines, "\n")
	}

	// Construct code in proper order
	var reorderedCode strings.Builder

	// Add class definition
	for _, line := range classLines {
		reorderedCode.WriteString(line + "\n")
	}

	// Add method definitions
	for _, line := range methodLines {
		reorderedCode.WriteString(line + "\n")
	}

	// Add method body (assuming we have only body content at this point)
	for _, line := range bodyLines {
		reorderedCode.WriteString(line + "\n")
	}

	debugLog("Successfully reorganized code")
	return reorderedCode.String()
}

// Helper function to determine if a line is likely Python code
func isCodeLine(line string) bool {
	// Skip lines that are clearly not code
	if strings.Contains(line, "Discussion Rules") ||
		strings.Contains(line, "Sort by") ||
		strings.Contains(line, "Comment") ||
		strings.Contains(line, "Online") ||
		strings.Contains(line, "Companies") ||
		strings.Contains(line, "Topics") ||
		strings.Contains(line, "Please don't post") ||
		strings.Contains(line, "Accepted") ||
		strings.Contains(line, "Submissions") {
		return false
	}

	// These are likely code
	if strings.Contains(line, "class ") ||
		strings.Contains(line, "def ") ||
		strings.Contains(line, "return ") ||
		strings.Contains(line, "if ") ||
		strings.Contains(line, "for ") ||
		strings.Contains(line, "while ") ||
		strings.Contains(line, "import ") {
		return true
	}

	// Look for common Python syntax patterns
	codePatterns := []string{
		"=", "+=", "-=", "*=", "/=", // assignment operators
		"==", "!=", "<", ">", "<=", ">=", // comparison operators
		"and", "or", "not", // logical operators
		"True", "False", "None", // Python constants
		"in ",                        // in operator with space to avoid matching inside words
		":",                          // common in Python blocks
		"[", "]", "{", "}", "(", ")", // brackets and parentheses
	}

	for _, pattern := range codePatterns {
		if strings.Contains(line, pattern) {
			return true
		}
	}

	// If we get here, we're not sure - default to including the line
	// unless it has evident HTML content
	if strings.Contains(line, "<div") ||
		strings.Contains(line, "<span") ||
		strings.Contains(line, "<button") ||
		strings.Contains(line, "<svg") {
		return false
	}

	return true
}
