# LangApp - Language Learning Made Easy with Flashcards

LangApp is a tool designed to help you learn new languages using flashcards. When reading a book in an unfamiliar language, you may come across many words you don't recognize. LangApp helps bridge that gap by generating flashcards containing words and their meanings. The generated flashcards are saved as a `.txt` file, which can be easily imported into Anki to assist in memorization.

## How It Works

LangApp leverages a comprehensive dictionary file to create a database of words and their meanings. Once the database is set up, you can generate flashcards based on any text you provide.

### Step 1: Parse the Dictionary

To get started, you'll need to parse a dictionary file from [kaikki.org](https://kaikki.org/dictionary/English/kaikki.org-dictionary-English.jsonl), which contains words in various languages along with their English meanings. This process takes around 3-5 hours.

**Command to parse the dictionary**:

```bash
go run cmd/parser/main.go -db ./assets/meaning/ -dict ./assets/kaikki.org-dictionary-English.jsonl -t 10 -p true

```

- **`db`**: Path where the parsed database will be stored.
- **`dict`**: Path to the dictionary file.
- **`t`**: Number of threads to use for faster processing.
- **`p`**: Enable progress output (set to `true`).

*Note*: The dictionary file should be stored in the `./assets` folder of LangApp.

### Step 2: Generate Flashcards

Once the dictionary has been parsed and the database is created, you can generate flashcards from any text.

**Command to generate flashcards**:

```bash
go run cmd/flashcard/main.go -t "Yo vivo en Granada, una ciudad pequeña que tiene monumentos muy importantes como la Alhambra. Aquí la comida es deliciosa y son famosos el gazpacho, el rebujito y el salmorejo."

```

- **`t`**: Text input containing words for which you want flashcards.

This command will generate a `cards.txt` file containing the words and their meanings. The file will be saved in the current directory and can be imported into Anki to create your personalized deck.

### Example Usage

1. **Parse the Dictionary**:
    
    ```bash
    go run cmd/parser/main.go -db ./assets/meaning/ -dict ./assets/kaikki.org-dictionary-English.jsonl -t 10 -p true
    
    ```
    
2. **Generate Flashcards**:
    
    ```bash
    go run cmd/flashcard/main.go -t "Bonjour, je m'appelle Pierre et j'aime voyager."
    
    ```
    
    This will create a `cards.txt` file with flashcards for words like "Bonjour" and "voyager" along with their English meanings.
    

### Happy Learning!

LangApp makes it easy to expand your vocabulary while reading in a new language. Use it to build your Anki decks and enhance your language learning journey.
