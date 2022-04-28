package csvalidator

import (
	"fmt"
	"io/fs"
	"os"
)

type File struct {
	firstIsHeader bool
	headers       map[string]struct{}
	records       [][]string
}

type Files []File

func NewFiles(path string, firstHeader bool) (Files, error) {
	var files Files
	if isDir(path) {
		err := fs.WalkDir(os.DirFS(path), ".", func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("failed read %s: %w", path, err)
			}
			if d.IsDir() {
				return nil
			}
			fullPath := fmt.Sprintf("%s/%s", path, p)
			if !isCSV(fullPath) {
				return nil
			}
			records, err := readCSV(fullPath)
			if err != nil {
				return err
			}
			f := File{
				firstIsHeader: firstHeader,
				headers:       getHeaders(records),
				records:       records,
			}
			files = append(files, f)
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("failed open %s: %w", path, err)
		}
		return files, nil
	}

	records, err := readCSV(path)
	if err != nil {
		return nil, fmt.Errorf("failed create validator: %w", err)
	}
	f := File{
		firstIsHeader: firstHeader,
		headers:       getHeaders(records),
		records:       records,
	}
	files = append(files, f)
	return files, nil
}

func isDir(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return false
	}
	i, err := f.Stat()
	if err != nil {
		fmt.Println()
		return false
	}
	return i.IsDir()
}
