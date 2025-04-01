package anki_deck

import (
	"fmt"
	"html"
	"os"
	"path/filepath"

	"github.com/npcnixel/genanki-go"
	"github.com/npcnixel/leetcode-to-anki-go/pkg/leetcode_to_anki"
)

// Controls debug logging
var debugLogging bool

func CreateDeck(problems []*leetcode_to_anki.Problem, outputDir string, debug bool) error {
	// Set debug flag for the package
	debugLogging = debug

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
body {
    background-color: #2d2d2d;
    color: #e0e0e0;
    font-family: Arial, sans-serif;
    max-width: 800px;
    margin: 0 auto;
    padding: 15px;
    line-height: 1.6;
}
h1 {
    color: #ff9800;
    font-size: 28px;
    margin-bottom: 18px;
    border-bottom: 3px solid #ff9800;
    padding-bottom: 10px;
    text-align: center;
    text-shadow: 1px 1px 2px rgba(0,0,0,0.7);
}
.title-container {
    background-color: #333;
    padding: 15px;
    margin-bottom: 20px;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0,0,0,0.3);
}
.description {
    font-size: 16px;
    margin-bottom: 20px;
    color: #e0e0e0;
}
.description p {
    margin-bottom: 15px;
}
.example {
    background-color: #333;
    padding: 15px;
    border-radius: 8px;
    margin: 15px 0;
    border: 1px solid #444;
}
.example-title {
    color: #ff9800;
    font-weight: bold;
    margin-bottom: 10px;
    font-size: 18px;
}
.example-content {
    margin-left: 15px;
}
.example-content pre {
    background-color: #2d2d2d;
    padding: 10px;
    border-radius: 4px;
    margin: 10px 0;
    font-family: 'Courier New', Courier, monospace;
    overflow-x: auto;
    border: 1px solid #444;
}
.constraints {
    background-color: #333;
    padding: 15px;
    border-radius: 8px;
    margin-top: 20px;
    border: 1px solid #444;
}
.constraints-title {
    color: #ff9800;
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
    box-shadow: 0 2px 4px rgba(0,0,0,0.3);
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
body {
    background-color: #2d2d2d;
    color: #e0e0e0;
    font-family: Arial, sans-serif;
    max-width: 800px;
    margin: 0 auto;
    padding: 15px;
}
.solution {
    background-color: #3d3d3d;
    padding: 15px;
    border-radius: 8px;
    border: 1px solid #555;
}
.solution pre {
    margin: 0;
    font-family: 'Courier New', Courier, monospace;
    font-size: 14px;
    line-height: 1.5;
    color: #e0e0e0;
    white-space: pre-wrap;
    overflow-x: auto;
}
img {
    max-width: 100%%;
    height: auto;
    display: block;
    margin: 10px auto;
    border-radius: 4px;
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
