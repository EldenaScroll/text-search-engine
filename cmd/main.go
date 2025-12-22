package main

import (
	"fmt"
	"log"
	"time"

	"github.com/EldenaScroll/text-search-engine/pkg/crawler"
	"github.com/EldenaScroll/text-search-engine/pkg/index"
	"github.com/EldenaScroll/text-search-engine/pkg/tokenizer"
	//"github.com/EldenaScroll/text-search-engine/pkg/index"
)

func main() {
	dataPath := "./data"
	stopWordsPath := "./stop-words-english.json"
	start := time.Now()

	stopWords, err := tokenizer.LoadStopWords(stopWordsPath)
	if err != nil {log.Fatal(err)}

	files, err := crawler.LoadDocuments(dataPath)

	if err != nil {log.Fatal(err)}
	fmt.Printf("Loaded %d documents in %v\n", len(files), time.Since(start))
	
	idx := index.NewIndex()

	if len(files) > 0 {

		for docID, content := range files{
			tokens := tokenizer.Tokenize(content)

			var cleanTokens []string
			for _, t := range tokens{
				if _, isStopWord := stopWords[t]; !isStopWord{
					cleanTokens = append(cleanTokens, t)
				}

			}
			fmt.Print(cleanTokens,"\n")
			idx.Add(docID, cleanTokens)
			fmt.Print(idx,"\n")
		}
		fmt.Printf("Total setup time: %v\n", time.Since(start))
	}
	
	
}
