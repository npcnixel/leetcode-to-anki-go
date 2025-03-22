# leetcode-to-anki-go

## Overview

LeetCode-to-Anki transforms solved LeetCode problems into personalized Anki flashcards, deliberately designed to reinforce **your unique problem-solving approaches**. Unlike generic solution memorization, this tool preserves the exact way *you* tackle each problem—your thought process, coding patterns, and implementation choices. By converting HTML pages of completed problems into structured flashcards with the problem on the front and *your specific solution* on the back, you create a personalized learning system that strengthens your distinctive problem-solving style.

## Features

- Captures and preserves *your personal implementation strategy* for each problem
- Extracts complete problem details including title, description, constraints, and examples
- Maintains your code comments that explain your individual thought process
- Preserves the exact syntax and structure of *your* solution with proper highlighting
- Processes both individual problems and batches of saved pages
- Creates clean, consistently formatted cards optimized for spaced repetition
- Embeds difficulty level and problem categories as searchable tags
- Automatically handles complex formatting elements like tables, math notation, and code blocks

TODO:
* HTML to .apkg
* batch os HTMLs to .apkg
* Check behaviour of uploading the same .apkg with updated cards (?) (ideally should add only new cards)
* OS-agnostic/add .dockerfile
* tests