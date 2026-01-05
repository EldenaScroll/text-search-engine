
package crawler

import ("fmt";"os"; "path/filepath"; "time")

func LoadDocuments(path string) ([]string,[]string, []time.Time, error){
	var files []string
	// store filenames
	var filenames []string

	var modTimes []time.Time

	err := filepath.WalkDir(path, func(s string, d os.DirEntry, e error) error {
		if e != nil { return e}
		// skip directories
		if d.IsDir() { return nil}
		//keep only txt files
		if filepath.Ext(s) != ".txt" {return nil}
		content, err := os.ReadFile(s)
		if err != nil {return err}

		files = append(files, string(content))
		filenames = append(filenames, d.Name())

		info, err := d.Info()
		if err != nil {return err}
		modTimes = append(modTimes, info.ModTime())
		return nil
	})

	return files, filenames, modTimes, err

}

// IsStale checks if the index file is older than any data file
// returns true if we need to rebuild the index
func IsStale(indexFilename string, dataDir string) bool {
	// get index file info
	indexStat, err := os.Stat(indexFilename)
	if os.IsNotExist(err) {return true}

	indexTime := indexStat.ModTime()

	// walk through the data folder
	isStale := false
	err = filepath.WalkDir(dataDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {return err}
		// skip directories
		if d.IsDir() {return nil}
		
		// get info for the data file
		info, err := d.Info()
		if err != nil {return nil}

		// check if this is file newer than the index
		if info.ModTime().After(indexTime) {
			fmt.Printf("Detected change in: %s\n", d.Name())
			isStale = true
			return fmt.Errorf("stale") // stop walking, need to rebuild
		}
		return nil
	})

	return isStale
}