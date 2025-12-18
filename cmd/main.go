package main

import (
	"fmt"
	"log"
	"time"

	"github.com/EldenaScroll/text-search-engine/pkg/crawler"
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

	files, err := crawler.LoadDocuments(dataPath)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Loaded %d documents in %v\n", len(files), time.Since(start))
	
	if len(files) > 0 {
		rawContent := files[0]
		tokens := tokenizer.Tokenize(rawContent)

		// filter out stop words
		var cleanTokens []string
		for _, token := range tokens {
			
			if _, isStopWord := stopWords[token]; !isStopWord {
				cleanTokens = append(cleanTokens, token)
			}
		}
		fmt.Print(cleanTokens)

	}
	fmt.Printf("Total setup time: %v\n", time.Since(start))
	
}
