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
	dataPath := "./data"
	stopWordsPath := "./stop-words-english.json"
	start := time.Now()

	stopWords, err := tokenizer.LoadStopWords(stopWordsPath)
	if err != nil {
		log.Fatal(err)
	}

	files, filenames, err := crawler.LoadDocuments(dataPath)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Loaded %d documents in %v\n", len(files), time.Since(start))

	idx := index.NewIndex()

	if len(files) > 0 {

		for docID, content := range files {
			tokens := tokenizer.Tokenize(content)

			var cleanTokens []string
			for _, t := range tokens {
				if _, isStopWord := stopWords[t]; !isStopWord {
					cleanTokens = append(cleanTokens, t)
				}

			}
			fmt.Print(cleanTokens, "\n")
			idx.Add(docID, cleanTokens)
			fmt.Print(idx, "\n")
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
				if i >= 10 {
					break
				}
				if result.DocID < len(filenames) {
					fmt.Printf("%d. %s (Score: %.2f)\n", i+1, filenames[result.DocID], result.Score)
				}
			}

		}
	}

}
