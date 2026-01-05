package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/EldenaScroll/text-search-engine/pkg/crawler"
	"github.com/EldenaScroll/text-search-engine/pkg/index"
	"github.com/EldenaScroll/text-search-engine/pkg/tokenizer"
)

func main() {
	indexPath := "index.gob"
	dataPath := "./data"
	stopWordsPath := "./stop-words-english.json"
	start := time.Now()

	stopWords, err := tokenizer.LoadStopWords(stopWordsPath)
	if err != nil {
		log.Fatal(err)
	}

	files, filenames, modTimes, err := crawler.LoadDocuments(dataPath)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Loaded %d documents in %v\n", len(files), time.Since(start))

	idx := index.NewIndex()
	_, err = os.Stat(indexPath)
	indexLoaded := false

	if err == nil {
		err := idx.Load(indexPath)
		if err != nil {log.Fatal(err)}
		if idx.IsStale(filenames, modTimes) {
			fmt.Println("Changes detected. Rebuilding index...")
			idx = index.NewIndex()
		} else {
			fmt.Println("Index is up to date.")
			indexLoaded = true

		}
	} 
	if !indexLoaded{
		for docID, content := range files {
			tokens := tokenizer.Tokenize(content)

			var cleanTokens []string
			for _, t := range tokens {
				if _, isStopWord := stopWords[t]; !isStopWord {
					cleanTokens = append(cleanTokens, t)
				}

			}
			idx.Add(docID, filenames[docID], cleanTokens, modTimes[docID])
		}
		if err := idx.Save(indexPath); err != nil {
			log.Printf("Warning: Failed to save index: %v", err)
		} else {
			fmt.Println("Index saved to disk.")
		}
	}
	fmt.Printf("Total setup time: %v\n", time.Since(start))
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\nSearch (or exit to quit)-> ")
		if !scanner.Scan() {
			break
		}
		query := scanner.Text()
		if query == "exit" {
			break
		}
		queryTokens := tokenizer.Tokenize(query)
		if len(queryTokens) == 0 {
			continue
		}

		var searchTokens []string

		for _, t := range queryTokens {
			if _, isStopWord := stopWords[t]; !isStopWord {
				searchTokens = append(searchTokens, t)
			}
		}
		matchedIDs := idx.Search(searchTokens)
		if len(matchedIDs) == 0 {
			fmt.Print("Not Found")
		}

		for i, result := range matchedIDs {
			if i >= 5 {
				break
			}
			if result.DocID < len(filenames) {
				snippet := index.ExtractSnippet(files[result.DocID], searchTokens[0])
				fmt.Printf("%d. %s (Score: %.2f)\n", i+1, filenames[result.DocID], result.Score)
				fmt.Printf("    %s\n\n", snippet)
			}
		}

	}
}
