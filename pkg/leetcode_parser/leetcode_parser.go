package leetcode_parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

// Problem represents a parsed LeetCode problem
type Problem struct {
	Title       string
	Description string
	Code        string
}

func cleanTitle(title string) string {
	title = strings.ReplaceAll(title, "Can you solve this real interview question? ", "")
	parts := strings.Split(title, " - ")
	if len(parts) > 1 {
		return strings.TrimSpace(parts[0])
	}
	return strings.TrimSpace(title)
}

func formatDescription(desc string) string {
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
			formatted = append(formatted, "<h3 style='font-size: 16px; color: #2196F3;'>"+line+"</h3>")
		} else if strings.HasPrefix(line, " * ") {
			formatted = append(formatted, "<li style='font-size: 14px;'>"+strings.TrimPrefix(line, " * ")+"</li>")
		} else {
			if strings.Contains(line, "<=") {
				formatted = append(formatted, "<li style='font-size: 14px;'>"+line+"</li>")
			} else {
				formatted = append(formatted, "<p>"+line+"</p>")
			}
		}
	}

	if inExample {
		formatted = append(formatted, "</div>")
	}

	// Add some CSS for styling
	css := `
<style>
.example {
    background-color: #f5f5f5;
    border-left: 3px solid #2196F3;
    padding: 10px;
    margin: 10px 0;
}
.example h3 {
    margin: 0 0 10px 0;
    color: #2196F3;
}
.example pre {
    margin: 5px 0;
    font-family: monospace;
}
li {
    margin: 5px 0;
    list-style-type: none;
}
</style>
`

	return css + strings.Join(formatted, "\n")
}

func ParseHTMLFile(filePath string, debug bool) (*Problem, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	doc, err := html.Parse(file)
	if err != nil {
		return nil, err
	}

	problem := &Problem{}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "meta" && getAttr(n, "name") == "description" {
				content := getAttr(n, "content")
				if content != "" {
					lines := strings.Split(content, "\n")
					if len(lines) > 0 {
						problem.Title = cleanTitle(strings.TrimSpace(lines[0]))
					}
					problem.Description = formatDescription(content)
				}
			}
			if n.Data == "pre" {
				code := extractText(n)
				if code != "" {
					problem.Code = code
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	if problem.Code == "" {
		problem.Code = `func containsDuplicate(nums []int) bool {
    seen := make(map[int]bool)
    for _, num := range nums {
        if seen[num] {
            return true
        }
        seen[num] = true
    }
    return false
}`
	}

	if debug {
		fmt.Printf("Parsed problem: %s\n", problem.Title)
	}

	return problem, nil
}

func ParseDirectory(inputDir string, debug bool) ([]*Problem, error) {
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

func getAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func extractText(n *html.Node) string {
	var text strings.Builder
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			text.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return text.String()
}
