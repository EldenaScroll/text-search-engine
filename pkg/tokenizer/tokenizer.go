package tokenizer

import (
	"encoding/json"
	"os"
	"strings"
	"unicode"
)

// helper function to read the json file and return a set of words
func LoadStopWords(path string) (map[string]struct{}, error) {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var wordList []string
	err = json.Unmarshal(fileData, &wordList)
	if err != nil {
		return nil, err
	}

	stopWordMap := make(map[string]struct{})

	for _, word := range wordList {
		stopWordMap[word] = struct{}{}
	}
	return stopWordMap, nil
}

func Tokenize(text string) []string {

	//make it lowercase
	text = strings.ToLower(text)
	return strings.FieldsFunc(text, func(r rune) bool {

		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

}
