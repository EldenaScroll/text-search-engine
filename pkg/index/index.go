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
