## Overview
leetcode-to-anki-go transforms solved LeetCode problems into personalized Anki flashcards, deliberately designed to reinforce **your unique problem-solving approaches**.

By converting HTML pages of completed problems into structured flashcards with the problem on the front and *your specific solution* on the back, you create a personalized learning system that strengthens your distinctive problem-solving style.

## Features

- Captures and preserves *your personal implementation strategy* for each problem
- Extracts complete problem details including title, description, constraints, and examples
- Maintains your code comments that explain your individual thought process
- Preserves the exact syntax and structure of *your* solution with proper highlighting
- Processes both individual problems and batches of saved pages
- Creates clean, consistently formatted cards optimized for spaced repetition
- Embeds difficulty level and problem categories as searchable tags
- Automatically handles complex formatting elements like tables, math notation, and code blocks

## Installation

### Prerequisites

- Go 1.19 or higher (fully compatible with Go 1.24.1)
- Anki (to import the generated deck)

### Installing from source

1. Clone this repository:
   ```
   git clone https://github.com/npcnixel/leetcode-to-anki-go.git
   cd leetcode-to-anki-go
   ```

2. Build the application:
   ```
   go build
   ```

## Usage

### Saving LeetCode Pages

1. Solve problems on LeetCode
2. **Save the complete page**: 
   - LeetCode uses GraphQL which means standard browser "Save as" (Ctrl+S/Cmd+S) might not capture the code
   You may use a browser extension like [SingleFile](https://chromewebstore.google.com/detail/singlefile/mpiodijhokgodhhofbcjdecpffjipkle) to capture the fully rendered HTML page including dynamically loaded content
3. Place the saved HTML files in the `input` directory

### Generating Anki Cards

Run the application:

```
go run main.go
```

#### Command Line Options

- `-debug`: Enable detailed debugging output
  ```
  go run main.go -debug
  ```

## Directory Structure

- `input/`: Place saved LeetCode HTML files here
- `output/`: Generated Anki package will be saved here

## How It Works

1. Parses the saved HTML files to extract problem titles, descriptions, and your solutions
2. Formats the content with proper styling for readability
3. Creates an Anki deck with cards that have the problem on the front and your solution on the back
4. Packages everything into a standard Anki package (`.apkg`) format

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

TODO:
* More tests
* Add .dockerfile
* Experiment with style
