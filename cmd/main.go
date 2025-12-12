
package main

import ("fmt"; "log"; "time"; "github.com/EldenaScroll/text-search-engine/pkg/crawler")

func main() {
	dataPath := "./data"

	start := time.Now()

	docs, err := crawler.LoadDocuments(dataPath)

	if err != nil {log.Fatal(err)}
	fmt.Printf("Loaded %d documents in %v\n", len(docs), time.Since(start))
}