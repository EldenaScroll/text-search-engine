
package crawler

import ("fmt"; "os"; "path/filepath")

func LoadDocuments(path string) ([]string, error){
	var files []string

	err := filepath.WalkDir(path, func(s string, d os.DirEntry, e error) error {
		if e != nil { return e}
		// skip directories
		if d.IsDir() { return nil}

		content, err := os.ReadFile(s)
		if err != nil {return err}

		fmt.Printf("Loaded: %s\n", s)
		files = append(files, string(content))
		return nil
	})
	return files, err

}