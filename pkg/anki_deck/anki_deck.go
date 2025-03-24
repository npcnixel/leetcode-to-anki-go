package anki_deck

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/npcnixel/genanki-go"
	"github.com/npcnixel/leetcode-to-anki-go/pkg/leetcode_parser"
)

func CreateDeck(problems []*leetcode_parser.Problem, outputDir string, debug bool) error {
	outputPath := filepath.Join(outputDir, "leetcode_deck.apkg")

	model := genanki.StandardBasicModel("LeetCode Problem")
	deck := genanki.StandardDeck("LeetCode Problems", "Collection of LeetCode problems and solutions")

	pkg := genanki.NewPackage([]*genanki.Deck{deck}).AddModel(model.Model)
	pkg.SetDebug(debug)

	// Add new notes
	for i, problem := range problems {
		front := fmt.Sprintf(`
<style>
h1 {
    color: #2196F3;
    font-size: 24px;
    margin-bottom: 20px;
    border-bottom: 2px solid #2196F3;
    padding-bottom: 10px;
}
.example h3 {
    font-size: 16px;
    margin: 0 0 10px 0;
    color: #2196F3;
}
.example pre {
    font-size: 12px;
    margin: 5px 0;
    font-family: monospace;
}
.example p {
    font-size: 14px;
}
</style>
<h1>%s</h1>
%s`, problem.Title, problem.Description)

		back := fmt.Sprintf(`
<style>
.solution {
    background-color: #f8f9fa;
    padding: 15px;
    border-radius: 5px;
    border: 1px solid #dee2e6;
}
.solution pre {
    margin: 0;
    font-family: 'Courier New', Courier, monospace;
    font-size: 14px;
    line-height: 1.5;
}
</style>
<div class="solution">
<pre><code>%s</code></pre>
</div>`, problem.Code)

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
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	if err := pkg.WriteToFile(outputPath); err != nil {
		return fmt.Errorf("failed to write package: %v", err)
	}

	return nil
}
