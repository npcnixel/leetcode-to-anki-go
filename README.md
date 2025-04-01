# leetcode-to-anki-go

## Overview
**leetcode-to-anki-go** transforms solved LeetCode problems into personalized Anki flashcards, deliberately designed to reinforce **your unique problem-solving approaches**.

<<<<<<< Updated upstream
<<<<<<< Updated upstream
<<<<<<< Updated upstream
By converting HTML pages of completed problems into structured flashcards with the problem on the front and *your personal solution* on the back, you create a personalized learning system that strengthens your distinctive problem-solving style.
=======
=======
>>>>>>> Stashed changes
=======
>>>>>>> Stashed changes
## Features

- **Complete HTML to Anki Conversion**: Transforms saved LeetCode pages into ready-to-import Anki packages
- **Personalized Learning**: Captures your unique solution approaches with original comments and code style
- **Image Preservation**: Maintains diagrams and illustrations from problem descriptions
- **Batch Processing**: Process multiple problems at once by adding HTML files to the input folder
- **Incremental Updates**: Adding new problems only creates cards for content not already in your Anki collection
- **Beautiful Formatting**: Dark-themed cards with proper syntax highlighting for better readability
- **Debug Mode**: Detailed logging with `-debug` flag to troubleshoot extraction issues
<<<<<<< Updated upstream
<<<<<<< Updated upstream
>>>>>>> Stashed changes
=======
>>>>>>> Stashed changes
=======
>>>>>>> Stashed changes

## Installation

1. Clone this repository:
   ```
   git clone https://github.com/npcnixel/leetcode-to-anki-go.git
   cd leetcode-to-anki-go
   ```

2. Run the application:
   ```
   go run main.go 
   ```

### Option 2: Docker (recommended for Windows users)

<details>
<summary><b>Using Docker (recommended for Windows users)</b> - Click to expand</summary>

If you're on Windows or prefer not to install Go locally, you can use Docker instead:

1. Install [Docker Desktop](https://www.docker.com/products/docker-desktop/)

2. Clone this repository:
   ```
   git clone https://github.com/npcnixel/leetcode-to-anki-go.git
   cd leetcode-to-anki-go
   ```

3. Build the Docker image:
   ```
   docker build -t leetcode-to-anki .
   ```

4. Run the container:
   ```
   docker run --rm -v "$(pwd)/input:/app/input" -v "$(pwd)/output:/app/output" leetcode-to-anki
   ```

   For Windows CMD:
   ```
   docker run --rm -v "%cd%/input:/app/input" -v "%cd%/output:/app/output" leetcode-to-anki
   ```

   For Windows PowerShell:
   ```
   docker run --rm -v "${PWD}/input:/app/input" -v "${PWD}/output:/app/output" leetcode-to-anki
   ```

   To run with debug mode, add the `-debug` flag:
   ```
   docker run --rm -v "$(pwd)/input:/app/input" -v "$(pwd)/output:/app/output" leetcode-to-anki -debug
   ```

5. The output will be available in the `output` directory, just as with the local installation

</details>

## Usage

1. LeetCode uses GraphQL which means standard browser "Save as" (Ctrl+S/Cmd+S) might not capture all content. Hence, use a browser extension like [SingleFile](https://chromewebstore.google.com/detail/singlefile/mpiodijhokgodhhofbcjdecpffjipkle) to capture the fully rendered page with all content

2. Place all saved HTML files in the `input` directory. You may add multiple problems at once for batch processing.

3. Run the application from the command line:

```
go run main.go
```
<details>
<summary><b>Import into Anki</b> - Click to expand</summary>

1. Locate the generated `.apkg` file in the `output` directory
2. Open Anki and select "File > Import" (or press Ctrl+Shift+I / Cmd+Shift+I)
3. Select the `.apkg` file and click "Open"
4. The cards will be added to your Anki collection

Note: Only new problems will be added as cards. If you've previously imported some problems, they won't be duplicated.
</details>


## Features

- **Complete HTML to Anki Conversion**: Transforms saved LeetCode pages into ready-to-import Anki packages
- **Image Preservation**: Maintains diagrams and illustrations from problem descriptions
- **Batch Processing**: Process multiple problems at once by adding HTML files to the input folder
- **Incremental Updates**: Adding new problems and importing .apkg only creates cards that didn't exist in previously imported Anki collection
- **Beautiful Formatting**: Dark-themed cards with proper syntax highlighting for better readability
- **Debug Mode**: Detailed logging with `-debug` flag to troubleshoot extraction issues
- **Cross-platform**: You can use Dockerfile for Windows

<details>
<summary><b>How It Works</b> - Click to expand</summary>

### Directory Structure

- `input/`: Place saved LeetCode HTML files here
- `output/`: Generated Anki package will be saved here

### Process

1. **HTML Parsing**: The application scans the `input` directory and parses all HTML files
2. **Content Extraction**: For each file, it extracts:
   - Problem title and unique identifier
   - Complete problem description with examples
   - Images and diagrams (embedded in the HTML)
   - Your solution code with comments and formatting
3. **Card Generation**: Creates Anki cards with:
   - Front: Problem statement with all examples and constraints
   - Back: Your complete solution with syntax highlighting
4. **Package Creation**: Builds an Anki package (`.apkg`) with all extracted content
5. **Incremental Updates**: When you import the package into Anki:
   - Only new problems are added as new cards
   - Existing problems are not duplicated
   - Your existing collection remains intact
6. **Debugging Support**: With the `-debug` flag, detailed logs show exactly what's being extracted and how it's being processed
<<<<<<< Updated upstream
<<<<<<< Updated upstream
<<<<<<< Updated upstream
7. Heavily relies on [genanki-go](https://github.com/npcnixel/genanki-go) library for generating notes, decks, package and so on.
</details>

## TODO:
* More tests
* Include pictures
* Experiment with style
* Review and test Dockerfile
=======
=======
>>>>>>> Stashed changes
=======
>>>>>>> Stashed changes

## TODO:
* More tests
* Add .dockerfile
* Experiment with style
<<<<<<< Updated upstream
<<<<<<< Updated upstream
>>>>>>> Stashed changes
=======
>>>>>>> Stashed changes
=======
>>>>>>> Stashed changes

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

<<<<<<< Updated upstream
<<<<<<< Updated upstream
<<<<<<< Updated upstream
This project is licensed under the MIT License - see the LICENSE file for details.
=======
This project is licensed under the MIT License - see the LICENSE file for details.
>>>>>>> Stashed changes
=======
This project is licensed under the MIT License - see the LICENSE file for details.
>>>>>>> Stashed changes
=======
This project is licensed under the MIT License - see the LICENSE file for details.
>>>>>>> Stashed changes
