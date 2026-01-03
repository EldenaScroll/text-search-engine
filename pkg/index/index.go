package index

type Index struct{ store map[string][]int }

// return an empty Index
func NewIndex() *Index { return &Index{store: make(map[string][]int)} }

// add indexes a document
func (idx *Index) Add(docID int, tokens []string) {
	for _, token := range tokens {
		ids := idx.store[token]
		if len(ids) > 0 && ids[len(ids)-1] == docID {
			continue
		}
		idx.store[token] = append(ids, docID)
	}

}

// Search takes a list of query tokens and returns IDs that contain all tokens.
func (idx *Index) Search(tokens []string) []int {
	// edge case
	if len(tokens) == 0 {
		return nil
	}

	firstWord := tokens[0]
	finalIDs, ok := idx.store[firstWord]
	if !ok {
		return nil
	}

	// check the rest of the words
	for i := 1; i < len(tokens); i++ {
		nextWord := tokens[i]
		nextIDs, ok := idx.store[nextWord]
		if !ok {
			return nil
		}
		
		// Intersect current results with this new list
		finalIDs = intersection(finalIDs, nextIDs)
		
		// If result becomes empty, stop early
		if len(finalIDs) == 0 {
			return nil
		}
	}

	return finalIDs
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
