package tests

import (
	"archive/zip"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"html"

	_ "github.com/mattn/go-sqlite3"
	"github.com/npcnixel/leetcode-to-anki-go/pkg/anki_deck"
	"github.com/npcnixel/leetcode-to-anki-go/pkg/leetcode_to_anki"
)

// Expected code constants for validation
const expectedStockCode = `class Solution:
    def maxProfit(self, prices: List[int]) -> int:
        if not prices or len(prices) == 1:
            return 0
            
        total_profit = 0
        left = 0  # Buy position
        
        for end in range(1, len(prices)):
            if prices[end] < prices[end - 1]:
                if prices[end - 1] > prices[left]:
                    total_profit += prices[end - 1] - prices[left]
                left = end
        
        if prices[-1] > prices[left]:
            total_profit += prices[-1] - prices[left]
            
        return total_profit`

const expectedRainWaterCode = `class Solution:
    def trap(self, height: List[int]) -> int:
        mh = max(height)
        block_vol = sum(height)

        left = 0
        max_left = 0
        while height[left] < mh:
            if height[left] > max_left:
                max_left = height[left]
            
            height[left] = max_left
            left += 1
        
        right = len(height)-1
        max_right = 0
        while height[right] < mh:
            if height[right] > max_right:
                max_right = height[right]
            
            height[right] = max_right
            right -= 1

        filled_vol = sum(height[0:left]+height[right:]) + (right-left)*mh
        return filled_vol - block_vol`

func TestIntegratedAnkiPackage(t *testing.T) {
	// Create a temporary test directory
	testDir, err := os.MkdirTemp("", "leetcode-anki-integrated-test")
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testDir) // Clean up after test

	// Create subdirectories
	inputDir := filepath.Join(testDir, "input")
	outputDir := filepath.Join(testDir, "output")
	extractDir := filepath.Join(testDir, "extract")

	for _, dir := range []string{inputDir, outputDir, extractDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Copy HTML files from tests folder to input directory
	htmlFiles := []string{
		"Best Time to Buy and Sell Stock II - LeetCode (25_03_2025 21：10：32).html",
		"Trapping Rain Water - LeetCode (25_03_2025 21：29：58).html",
	}

	for _, filename := range htmlFiles {
		sourcePath := filepath.Join(".", filename) // Files are in tests folder
		destPath := filepath.Join(inputDir, filename)

		content, err := os.ReadFile(sourcePath)
		if err != nil {
			t.Fatalf("Failed to read file %s: %v", sourcePath, err)
		}

		if err := os.WriteFile(destPath, content, 0644); err != nil {
			t.Fatalf("Failed to write file %s: %v", destPath, err)
		}
	}

	// Parse HTML files
	problems, err := leetcode_to_anki.ParseDirectory(inputDir, true)
	if err != nil {
		t.Fatalf("Failed to parse HTML files: %v", err)
	}

	// Verify parsing results
	if len(problems) != 2 {
		t.Fatalf("Expected 2 problems, got %d", len(problems))
	}

	// Check problem titles and code
	expectedTitles := map[string]bool{
		"Best Time to Buy and Sell Stock II": false,
		"Trapping Rain Water":                false,
	}

	// Define expected code for each problem
	expectedCodeMap := map[string]string{
		"Best Time to Buy and Sell Stock II": expectedStockCode,
		"Trapping Rain Water":                expectedRainWaterCode,
	}

	for _, problem := range problems {
		if _, exists := expectedTitles[problem.Title]; exists {
			expectedTitles[problem.Title] = true
			t.Logf("Found problem: %s", problem.Title)

			// Verify code against expected
			expectedCode := expectedCodeMap[problem.Title]

			// Normalize both codes for comparison (remove extra whitespace)
			normalizedExpected := normalizeCode(expectedCode)
			normalizedActual := normalizeCode(problem.Code)

			if normalizedActual != normalizedExpected {
				t.Errorf("Code mismatch for problem %s\nExpected:\n%s\n\nActual:\n%s",
					problem.Title, expectedCode, problem.Code)
			} else {
				t.Logf("Code for problem %s matches expected", problem.Title)
			}
		} else {
			t.Errorf("Unexpected problem title: %s", problem.Title)
		}

		// Ensure problem has code and description
		if problem.Code == "" {
			t.Errorf("Problem %s has no code", problem.Title)
		}
		if problem.Description == "" {
			t.Errorf("Problem %s has no description", problem.Title)
		}
	}

	// Check all expected titles were found
	for title, found := range expectedTitles {
		if !found {
			t.Errorf("Expected problem %s not found", title)
		}
	}

	// Create Anki deck
	packagePath := filepath.Join(outputDir, "leetcode_deck.apkg")
	anki_deck.CreateDeck(problems, outputDir, true)

	// Verify the Anki package was created
	if _, err := os.Stat(packagePath); os.IsNotExist(err) {
		t.Fatalf("Anki package was not created at %s", packagePath)
	}

	// Extract and verify the contents of the Anki package
	verifyAnkiPackageContents(t, packagePath, extractDir, expectedTitles, expectedCodeMap)
}

// normalizeCode removes extra whitespace and normalizes line endings for code comparison
func normalizeCode(code string) string {
	// Replace all whitespace sequences with a single space
	re := regexp.MustCompile(`\s+`)
	normalized := re.ReplaceAllString(code, " ")

	// Trim leading/trailing whitespace
	normalized = strings.TrimSpace(normalized)

	return normalized
}

func verifyAnkiPackageContents(t *testing.T, packagePath, extractDir string, expectedTitles map[string]bool, expectedCodeMap map[string]string) {
	// Extract the Anki package (it's a zip file)
	err := extractAnkiPackage(packagePath, extractDir)
	if err != nil {
		t.Fatalf("Failed to extract Anki package: %v", err)
	}

	// Verify collection.anki2 exists (SQLite database)
	dbPath := filepath.Join(extractDir, "collection.anki2")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Fatalf("collection.anki2 database not found in extracted package")
	}

	// Open the SQLite database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("Failed to open SQLite database: %v", err)
	}
	defer db.Close()

	// Query the notes table to count notes
	var noteCount int
	err = db.QueryRow("SELECT COUNT(*) FROM notes").Scan(&noteCount)
	if err != nil {
		t.Fatalf("Failed to query notes table: %v", err)
	}

	// We expect one note per problem
	expectedCount := len(expectedTitles)
	if noteCount != expectedCount {
		t.Errorf("Expected %d notes, found %d", expectedCount, noteCount)
	} else {
		t.Logf("Found %d notes in the Anki package (expected %d)", noteCount, expectedCount)
	}

	// Check each note for expected content
	rows, err := db.Query("SELECT flds FROM notes")
	if err != nil {
		t.Fatalf("Failed to query notes content: %v", err)
	}
	defer rows.Close()

	var foundTitles = make(map[string]bool)
	for title := range expectedTitles {
		foundTitles[title] = false
	}

	for rows.Next() {
		var fields string
		if err := rows.Scan(&fields); err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}

		// Check for problem title in the fields
		for title := range expectedTitles {
			// Use regexp.QuoteMeta to escape all regex special characters in the title
			escapedTitle := regexp.QuoteMeta(title)

			// Look for title in h1 tags
			titlePattern := fmt.Sprintf(`<h1>%s</h1>`, escapedTitle)
			titleFound, _ := regexp.MatchString(titlePattern, fields)

			// Also check div with title-container class
			titleContainerPattern := fmt.Sprintf(`<div class="title-container">\s*<h1>%s</h1>`, escapedTitle)
			titleContainerFound, _ := regexp.MatchString(titleContainerPattern, fields)

			if titleFound || titleContainerFound {
				foundTitles[title] = true

				// Check for code content - get the code from the HTML
				codePattern := `<pre><code>([\s\S]*?)</code></pre>`
				re := regexp.MustCompile(codePattern)
				matches := re.FindStringSubmatch(fields)

				if len(matches) < 2 {
					t.Errorf("Could not extract code for %s", title)
					continue
				}

				// Extract the code and compare with expected
				extractedCode := matches[1]
				expectedCode := expectedCodeMap[title]

				// Normalize for comparison
				normalizedExtracted := normalizeCode(html.UnescapeString(extractedCode))
				normalizedExpected := normalizeCode(expectedCode)

				// Check if the code matches the expected code
				if !strings.Contains(normalizedExtracted, normalizedExpected) {
					t.Errorf("Code in Anki card doesn't match expected for %s\nExpected to contain:\n%s\n\nGot:\n%s",
						title, expectedCode, extractedCode)
				} else {
					t.Logf("Code in Anki card for %s matches expected", title)
				}

				// Check for description elements
				if !strings.Contains(fields, "example") && !strings.Contains(fields, "Example") {
					t.Errorf("Expected examples in %s description but not found", title)
				}
			}
		}
	}

	// Check all expected titles were found in the notes
	for title, found := range foundTitles {
		if !found {
			t.Errorf("Problem %s not found in Anki package", title)
		}
	}

	if len(foundTitles) == len(expectedTitles) {
		t.Log("All expected problems found in Anki package with correct formatting")
	}
}

func extractAnkiPackage(packagePath, extractDir string) error {
	// Open the zip file
	reader, err := zip.OpenReader(packagePath)
	if err != nil {
		return fmt.Errorf("failed to open zip archive: %w", err)
	}
	defer reader.Close()

	// Extract each file
	for _, file := range reader.File {
		// Open the file inside zip
		zipFile, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in zip: %w", err)
		}

		// Create the file path for extraction
		extractPath := filepath.Join(extractDir, file.Name)

		// Create directories if needed
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(extractPath, 0755); err != nil {
				zipFile.Close()
				return fmt.Errorf("failed to create directory: %w", err)
			}
			continue
		}

		// Create the parent directories if they don't exist
		if err := os.MkdirAll(filepath.Dir(extractPath), 0755); err != nil {
			zipFile.Close()
			return fmt.Errorf("failed to create parent directory: %w", err)
		}

		// Create the file
		outFile, err := os.Create(extractPath)
		if err != nil {
			zipFile.Close()
			return fmt.Errorf("failed to create file: %w", err)
		}

		// Copy contents
		if _, err := io.Copy(outFile, zipFile); err != nil {
			outFile.Close()
			zipFile.Close()
			return fmt.Errorf("failed to copy file contents: %w", err)
		}

		// Close files
		outFile.Close()
		zipFile.Close()
	}

	return nil
}
