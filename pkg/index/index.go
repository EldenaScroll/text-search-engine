package index

import (
	"math"
	"sort"
	"strings"
	"encoding/gob"
	"os"
	"time"
)

type Posting struct {
	DocID int
	Count int
	
}

type Index struct {
	Store     map[string][]Posting
	TotalDocs int
	FileMeta map[string]time.Time
}

type SearchResult struct {
	DocID int
	Score float64
}

func NewIndex() *Index { 
	return &Index{
		Store: make(map[string][]Posting),
		TotalDocs: 0, 
		FileMeta: make(map[string]time.Time),
	} 
}

// add indexes a document
func (idx *Index) Add(docID int, filename string, tokens []string, modTime time.Time) {
	idx.FileMeta[filename] = modTime
	idx.TotalDocs++
	
	// map word to frequency
	termFreq := make(map[string]int)

	for _, token := range tokens {
		termFreq[token]++
	}

	for token, count := range termFreq {
		posting := Posting{DocID: docID, Count: count}
		idx.Store[token] = append(idx.Store[token], posting)
	}

}

// Search takes a list of query tokens and returns result
func (idx *Index) Search(tokens []string) []SearchResult {
	// edge case
	if len(tokens) == 0 {
		return nil
	}

	firstWord := tokens[0]
	postings, ok := idx.Store[firstWord]
	if !ok {
		return nil
	}

	docScores := make(map[int]float64)
	//initialization
	idf := math.Log(float64(idx.TotalDocs) / float64(len(postings)))

	for _, p := range postings {
		docScores[p.DocID] = float64(p.Count) * idf
	}

	// check the rest of the words
	for i := 1; i < len(tokens); i++ {
		nextWord := tokens[i]
		nextIDs, ok := idx.Store[nextWord]
		if !ok {
			return nil
		}

		nextDocScores := make(map[int]float64)

		for _, p := range nextIDs {

			if currentScore, exists := docScores[p.DocID]; exists {
				nextDocScores[p.DocID] = currentScore + (float64(p.Count) * idf)
			}
		}

		docScores = nextDocScores

		// If result becomes empty, stop early
		if len(docScores) == 0 {
			return nil
		}
	}

	var results []SearchResult
	for docID, score := range docScores {
		results = append(results, SearchResult{DocID: docID, Score: score})
	}

	// sort by Score (High to Low)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}

func ExtractSnippet(text string, term string) string {
	// find the term (Case Insensitive)
	lowerText := strings.ToLower(text)
	lowerTerm := strings.ToLower(term)

	idx := strings.Index(lowerText, lowerTerm)

	// define the Window
	start := idx - 30
	end := idx + len(term) + 50

	prefix := "..."
	suffix := "..."

	if start <= 0 {
		start = 0
		prefix = ""
	}
	
	if end >= len(text) {
		end = len(text)
		suffix = ""
	}

	snippet := text[start:end]
	
	return prefix + snippet + suffix
}

// write the index to a binary file
func (idx *Index) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(idx); err != nil {
		return err
	}
	
	return nil
}

// load reads the index from a binary file
func (idx *Index) Load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	
	// decode directly into the current Index struct
	if err := decoder.Decode(idx); err != nil {
		return err
	}
	
	return nil
}

// IsStale checks if the on-disk files have changed compared to the index
func (idx *Index) IsStale(currentFileNames []string, currentModTimes []time.Time) bool {
	// different number of files (added or deleted)
	if len(currentFileNames) != len(idx.FileMeta) {
		return true
	}

	// check every file
	for i, filename := range currentFileNames {
		storedTime, exists := idx.FileMeta[filename]
		
		// new file found (Rename/Swap)
		if !exists {
			return true
		}

		// file modified (timestamp changed)
		if !storedTime.Equal(currentModTimes[i]) {
			return true
		}
	}

	// everything matches
	return false
}

func intersection(a, b []int) []int {
	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}

	// avoid resizing overhead
	r := make([]int, 0, maxLen)

	var i, j int
	for i < len(a) && j < len(b) {
		if a[i] < b[j] {
			i++
		} else if a[i] > b[j] {
			j++
		} else {
			// match found
			r = append(r, a[i])
			i++
			j++
		}
	}
	return r
}
