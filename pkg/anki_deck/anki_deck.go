package anki_deck

import (
	"fmt"
	"html"
	"os"
	"path/filepath"

	"github.com/npcnixel/genanki-go"
	"github.com/npcnixel/leetcode-to-anki-go/pkg/leetcode_to_anki"
)

func CreateDeck(problems []*leetcode_to_anki.Problem, outputDir string, debug bool) error {
	outputPath := filepath.Join(outputDir, "leetcode_deck.apkg")

	model := genanki.StandardBasicModel("LeetCode Problem")
	deck := genanki.StandardDeck("LeetCode Problems", "Collection of LeetCode problems and solutions")

	pkg := genanki.NewPackage([]*genanki.Deck{deck}).AddModel(model.Model)
	pkg.SetDebug(debug)

	// Add new notes
	for i, problem := range problems {
		// Format the description to make it more readable
		formattedDescription := leetcode_to_anki.FormatDescription(problem.Description, problem.Images)

		// Front of the card shows the problem title and description
		front := fmt.Sprintf(`
<style>
/* Default styling (light theme) */
body {
    background-color: #ffffff;
    color: #333333;
    font-family: Arial, sans-serif;
    max-width: 800px;
    margin: 0 auto;
    padding: 15px;
    line-height: 1.6;
}
h1 {
    color: #ff6b00;
    font-size: 28px;
    margin-bottom: 18px;
    border-bottom: 3px solid #ff6b00;
    padding-bottom: 10px;
    text-align: center;
    text-shadow: 1px 1px 2px rgba(0,0,0,0.2);
}
.title-container {
    background-color: #f5f5f5;
    padding: 15px;
    margin-bottom: 20px;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}
.description {
    font-size: 16px;
    margin-bottom: 20px;
    color: #333333;
}
.description p {
    margin-bottom: 15px;
}
.example {
    background-color: #f5f5f5 !important;
    padding: 15px;
    border-radius: 8px;
    margin: 15px 0;
    border: 1px solid #e0e0e0;
}
.example-title {
    color: #ff6b00;
    font-weight: bold;
    margin-bottom: 10px;
    font-size: 18px;
}
.example-content {
    margin-left: 15px;
}
.example-content pre, pre {
    background-color: #ffffff !important;
    padding: 10px;
    border-radius: 4px;
    margin: 10px 0;
    font-family: 'Courier New', Courier, monospace;
    overflow-x: auto;
    border: 1px solid #e0e0e0;
    color: #333333 !important;
}
.constraints {
    background-color: #f5f5f5;
    padding: 15px;
    border-radius: 8px;
    margin-top: 20px;
    border: 1px solid #e0e0e0;
}
.constraints-title {
    color: #ff6b00;
    font-weight: bold;
    margin-bottom: 10px;
    font-size: 18px;
}
img {
    max-width: 100%%;
    height: auto;
    display: block;
    margin: 10px auto;
    border-radius: 4px;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}
code {
    background-color: #f8f8f8;
    padding: 2px 4px;
    border-radius: 3px;
    font-family: 'Courier New', Courier, monospace;
    color: #333333 !important;
    border: 1px solid #e0e0e0;
}
pre code {
    display: block;
    padding: 10px;
    overflow-x: auto;
    background-color: #ffffff !important;
    color: #333333 !important;
    border: none;
}

/* Force light styling for all example blocks and code blocks */
div[style*="background-color: #2d2d2d"], 
div[style*="background-color: #333"], 
div[style*="background-color: #3d3d3d"] {
    background-color: #f5f5f5 !important;
    color: #333333 !important;
    border: 1px solid #e0e0e0 !important;
}

/* Target image containers */
div[style*="background-color"][style*="border-radius"] {
    background-color: #ffffff !important;
    border: 1px solid #e0e0e0 !important;
}

/* Target any div that might be an example container */
div > div {
    background-color: #ffffff !important;
    color: #333333 !important;
}

/* Dark theme styling */
.night_mode body {
    background-color: #2d2d2d;
    color: #e0e0e0;
}
.night_mode h1 {
    color: #ff9800;
    text-shadow: 1px 1px 2px rgba(0,0,0,0.7);
    border-bottom: 3px solid #ff9800;
}
.night_mode .title-container {
    background-color: #333333;
    box-shadow: 0 2px 4px rgba(0,0,0,0.3);
}
.night_mode .description {
    color: #e0e0e0;
}
.night_mode .example, .night_mode div[style*="background-color: #2d2d2d"], .night_mode div[style*="background-color: #333"], .night_mode div[style*="background-color: #3d3d3d"] {
    background-color: #333333 !important;
    border: 1px solid #444444 !important;
    color: #e0e0e0 !important;
}
.night_mode .example-title {
    color: #ff9800;
}
.night_mode .example-content pre, .night_mode pre {
    background-color: #3d3d3d !important;
    border: 1px solid #444444;
    color: #e0e0e0 !important;
}
.night_mode .constraints {
    background-color: #333333;
    border: 1px solid #444444;
}
.night_mode .constraints-title {
    color: #ff9800;
}
.night_mode img {
    box-shadow: 0 2px 4px rgba(0,0,0,0.3);
}
.night_mode code {
    background-color: #3d3d3d;
    color: #e0e0e0 !important;
    border: 1px solid #444444;
}
.night_mode pre code {
    background-color: #3d3d3d !important;
    color: #e0e0e0 !important;
    border: none;
}

/* Target image containers in dark mode */
.night_mode div[style*="background-color"][style*="border-radius"] {
    background-color: #2d2d2d !important;
    border: 1px solid #444444 !important;
}

/* Target any div that might be an example container in dark mode */
.night_mode div > div {
    background-color: #2d2d2d !important;
    color: #e0e0e0 !important;
}
</style>
<div class="title-container">
  <h1>%s</h1>
</div>
<div class="description">
%s
</div>
`, problem.Title, formattedDescription)

		// Back of the card shows only the solution code
		back := fmt.Sprintf(`
<style>
/* Default styling (light theme) */
body {
    background-color: #ffffff;
    color: #333333;
    font-family: Arial, sans-serif;
    max-width: 800px;
    margin: 0 auto;
    padding: 15px;
}
.solution {
    background-color: #f5f5f5 !important;
    padding: 15px;
    border-radius: 8px;
    border: 1px solid #e0e0e0;
}
.solution pre {
    margin: 0;
    font-family: 'Courier New', Courier, monospace;
    font-size: 14px;
    line-height: 1.5;
    color: #333333 !important;
    white-space: pre-wrap;
    overflow-x: auto;
    background-color: #ffffff !important;
}
.solution code {
    display: block;
    padding: 0;
    background-color: transparent;
    color: #333333 !important;
    border: none;
}

/* Force light styling for all dark blocks */
div[style*="background-color: #2d2d2d"], 
div[style*="background-color: #333"], 
div[style*="background-color: #3d3d3d"] {
    background-color: #f5f5f5 !important;
    color: #333333 !important;
    border: 1px solid #e0e0e0 !important;
}

/* Target image containers */
div[style*="background-color"][style*="border-radius"] {
    background-color: #ffffff !important;
    border: 1px solid #e0e0e0 !important;
}

/* Target any div that might be an example container */
div > div {
    background-color: #ffffff !important;
    color: #333333 !important;
}

pre[style*="background-color"], code[style*="background-color"] {
    background-color: #ffffff !important;
    color: #333333 !important;
}
img {
    max-width: 100%%;
    height: auto;
    display: block;
    margin: 10px auto;
    border-radius: 4px;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

/* Dark theme styling */
.night_mode body {
    background-color: #2d2d2d;
    color: #e0e0e0;
}
.night_mode .solution {
    background-color: #3d3d3d !important;
    border: 1px solid #555555;
}
.night_mode .solution pre {
    color: #e0e0e0 !important;
    background-color: #3d3d3d !important;
}
.night_mode .solution code {
    color: #e0e0e0 !important;
}
.night_mode div[style*="background-color: #2d2d2d"], 
.night_mode div[style*="background-color: #333"], 
.night_mode div[style*="background-color: #3d3d3d"] {
    background-color: #3d3d3d !important;
    color: #e0e0e0 !important;
    border: 1px solid #555555 !important;
}

/* Target image containers in dark mode */
.night_mode div[style*="background-color"][style*="border-radius"] {
    background-color: #2d2d2d !important;
    border: 1px solid #444444 !important;
}

/* Target any div that might be an example container in dark mode */
.night_mode div > div {
    background-color: #2d2d2d !important;
    color: #e0e0e0 !important;
}

.night_mode pre[style*="background-color"], .night_mode code[style*="background-color"] {
    background-color: #3d3d3d !important;
    color: #e0e0e0 !important;
}
.night_mode img {
    box-shadow: 0 2px 4px rgba(0,0,0,0.3);
}
</style>
<div class="solution">
<pre><code>%s</code></pre>
</div>
`, html.EscapeString(problem.Code))

		if debug {
			fmt.Printf("\n=== Front of card for %s ===\n%s\n", problem.Title, front)
			fmt.Printf("\n=== Back of card for %s ===\n%s\n", problem.Title, back)
		}

		note := genanki.NewNote(
			model.ID,
			[]string{front, back},
			[]string{"leetcode", "programming"},
		)
		note.ID = int64(i + 1)
		deck.AddNote(note)

		// Add images to the package
		if len(problem.Images) > 0 {
			if debug {
				fmt.Printf("Adding %d images for problem: %s\n", len(problem.Images), problem.Title)
			}

			for _, img := range problem.Images {
				if debug {
					fmt.Printf("  Adding image: %s (%d bytes)\n", img.Filename, len(img.Data))
				}
				pkg.AddMedia(img.Filename, img.Data)
			}
		}
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	if err := pkg.WriteToFile(outputPath); err != nil {
		return fmt.Errorf("failed to write package: %v", err)
	}

	fmt.Printf("Successfully created Anki package: %s\n", outputPath)
	fmt.Printf("Added %d new notes\n", len(problems))

	return nil
}
