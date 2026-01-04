package index

import (
	"math"
	"sort"
)

type Posting struct {
	DocID int
	Count int
}

type Index struct {
	store     map[string][]Posting
	totalDocs int
}

type SearchResult struct {
	DocID int
	Score float64
}

func NewIndex() *Index { return &Index{store: make(map[string][]Posting), totalDocs: 0} }

// add indexes a document
func (idx *Index) Add(docID int, tokens []string) {
	idx.totalDocs++

	// map word to frequency
	termFreq := make(map[string]int)

	for _, token := range tokens {
		termFreq[token]++
	}

	for token, count := range termFreq {
		posting := Posting{DocID: docID, Count: count}
		idx.store[token] = append(idx.store[token], posting)
	}

}

// Search takes a list of query tokens and returns result
func (idx *Index) Search(tokens []string) []SearchResult {
	// edge case
	if len(tokens) == 0 {
		return nil
	}

	firstWord := tokens[0]
	postings, ok := idx.store[firstWord]
	if !ok {
		return nil
	}

	docScores := make(map[int]float64)
	//initialization
	idf := math.Log(float64(idx.totalDocs) / float64(len(postings)))

	for _, p := range postings {
		docScores[p.DocID] = float64(p.Count) * idf
	}

	// check the rest of the words
	for i := 1; i < len(tokens); i++ {
		nextWord := tokens[i]
		nextIDs, ok := idx.store[nextWord]
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
