package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/npcnixel/leetcode-to-anki-go/pkg/anki_deck"
	"github.com/npcnixel/leetcode-to-anki-go/pkg/leetcode_parser"
)

func main() {
	debug := flag.Bool("debug", false, "Show raw HTML/CSS when creating decks")
	flag.Parse()

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}

	inputDir := filepath.Join(currentDir, "input")
	outputDir := filepath.Join(currentDir, "output")

	problems, err := leetcode_parser.ParseDirectory(inputDir, *debug)
	if err != nil {
		log.Fatalf("Failed to parse HTML files: %v", err)
	}

	if len(problems) == 0 {
		log.Println("No HTML files found in the input directory")
		return
	}

	if err := anki_deck.CreateDeck(problems, outputDir, *debug); err != nil {
		log.Fatalf("Failed to create Anki deck: %v", err)
	}

	fmt.Printf("Successfully processed %d problems\n", len(problems))
	fmt.Printf("Anki deck created in: %s\n", filepath.Join(outputDir, "leetcode_deck.apkg"))
}
