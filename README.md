# leetcode-to-anki-go

## TL;DR
| Save as SingleFile.html | Input Folder | Output Folder |
|:-----:|:-------:|:------:|
| **leetcode.com/problem-x** | **SingleFile.html** | **anki_deck.apkg** |

Transform your solved LeetCode problems into personalized Anki flashcards, deliberately designed to reinforce **your unique problem-solving approaches**.

By converting HTML pages of completed problems into structured flashcards with the problem on the front and *your specific solution* on the back, you create a personalized learning system that strengthens your distinctive problem-solving style.

<img width="1180" alt="SCR-20250402-uanu" src="https://github.com/user-attachments/assets/98cb99f3-a584-4048-9aab-7f1418fc1b57" />

## ✨ Features

- **Complete HTML to Anki Conversion**: Transforms saved LeetCode pages into ready-to-import Anki packages
- **Image Preservation**: Maintains diagrams and illustrations from problem descriptions
- **Batch Processing**: Process multiple problems at once by adding HTML files to the input folder
- **Incremental Updates**: Adding new problems only creates cards for content not already in your Anki collection
- **Theme Support**: Cards automatically adapt to Anki's light/dark theme settings for comfortable studying in any environment
- **Debug Mode**: Detailed logging with `-debug` flag to troubleshoot extraction issues

## 🚀 Installation

<details>
<summary><b>Option 1: Local Installation</b> - Click to expand</summary>

1. Clone this repository:
   ```bash
   git clone https://github.com/npcnixel/leetcode-to-anki-go.git
   cd leetcode-to-anki-go
   ```

2. Run the application:
   ```bash
   go run main.go 
   ```
</details>

<details>
<summary><b>Option 2: Docker (recommended for Windows users)</b> - Click to expand</summary>

1. Install [Docker Desktop](https://www.docker.com/products/docker-desktop/)

2. Clone this repository:
   ```bash
   git clone https://github.com/npcnixel/leetcode-to-anki-go.git
   cd leetcode-to-anki-go
   ```

3. Build the Docker image:
   ```bash
   docker build -t leetcode-to-anki-go .
   ```

4. Run the container:
   ```bash
   # For macOS/Linux:
   docker run --rm -v "$(pwd)/input:/app/input" -v "$(pwd)/output:/app/output" leetcode-to-anki-go

   # For Windows CMD:
   docker run --rm -v "%cd%/input:/app/input" -v "%cd%/output:/app/output" leetcode-to-anki-go
   
   # For Windows PowerShell:
   docker run --rm -v "${PWD}/input:/app/input" -v "${PWD}/output:/app/output" leetcode-to-anki-go
   ```

5. The output will be available in the `output` directory, just as with the local installation
</details>

## Usage

1. **Save LeetCode Problems**: 
   - LeetCode uses GraphQL which means standard browser "Save as" (Ctrl+S/Cmd+S) might not capture all content.
   - ‼️ Use a browser extension like [SingleFile](https://chromewebstore.google.com/detail/singlefile/mpiodijhokgodhhofbcjdecpffjipkle) to capture the fully rendered page with all content ‼️

2. **Prepare Input Files**:
   - Place all saved HTML files in the `input` directory

3. **Generate Anki Cards**:
   - Run the application from the command line:
     ```bash
     go run main.go
     ```

4. **Import into Anki**:
   - Locate the generated `.apkg` file in the `output` directory
   - Open Anki and select "File > Import" (or press Ctrl+Shift+I / Cmd+Shift+I)
   - Select the `.apkg` file and click "Open"
   - The cards will be added to your Anki collection

   Note: Only new problems will be added as cards. If you've previously imported some problems, they won't be duplicated.

## 🔧 How It Works

1. Saved single-file HTML files are parsed to extract problem titles, descriptions, and your solutions
2. Content is parsed with proper styling for readability
3. Anki deck is created with cards that have the problem on the front and your solution on the back
4. Everything is packaged into a standard Anki package (`.apkg`) format

### Directory Structure

- `input/`: Place saved LeetCode HTML files here
- `output/`: Generated Anki package will be saved here

## 🤙 Dependencies

This tool heavily relies on [genanki-go](https://github.com/npcnixel/genanki-go) for comprehensive ANKI deck construction, complete customization and professional flashcard generation in Go with full control over deck structure and formatting.

## 🙌 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.
