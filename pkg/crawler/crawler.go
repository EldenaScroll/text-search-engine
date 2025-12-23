
package crawler

import ("fmt"; "os"; "path/filepath")

func LoadDocuments(path string) ([]string,[]string, error){
	var files []string
	// store filenames
	var filenames []string

	err := filepath.WalkDir(path, func(s string, d os.DirEntry, e error) error {
		if e != nil { return e}
		// skip directories
		if d.IsDir() { return nil}
		//keep only txt files
		if filepath.Ext(s) != ".txt" {return nil}
		content, err := os.ReadFile(s)
		if err != nil {return err}

		fmt.Printf("Loaded: %s\n", s)
		files = append(files, string(content))
		filenames = append(filenames, d.Name())
		return nil
	})
	return files, filenames, err

}