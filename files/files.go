package files

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/markelrep/csvalidator/csv"
)

// File represent CSV file
type File struct {
	filePath      string
	firstIsHeader bool
	headers       map[string]struct{}
	headersLen    int
	Records       [][]string
}

// Files slice of File
type Files []File

// NewFiles create a new Files from path
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
			fullPath := path + "/" + p
			records, h, err := csv.ReadCSV(fullPath, firstHeader)
			if err != nil {
				return err
			}
			f := File{
				filePath:      fullPath,
				firstIsHeader: firstHeader,
				headers:       h,
				headersLen:    len(h),
				Records:       records,
			}
			files = append(files, f)
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("failed open %s: %w", path, err)
		}
		return files, nil
	}

	records, h, err := csv.ReadCSV(path, firstHeader)
	if err != nil {
		return nil, fmt.Errorf("failed create validator: %w", err)
	}
	f := File{
		filePath:      path,
		firstIsHeader: firstHeader,
		headers:       h,
		headersLen:    len(h),
		Records:       records,
	}
	files = append(files, f)
	return files, nil
}

// FirstIsHeader gets information about first line, is it header or not
func (f File) FirstIsHeader() bool {
	return f.firstIsHeader
}

// HeadersCount return len of headers map
func (f File) HeadersCount() int {
	return f.headersLen
}

// HasHeader checks if specific header exists or not
func (f File) HasHeader(name string) bool {
	_, ok := f.headers[name]
	return ok
}

// Path returns path to file
func (f File) Path() string {
	return f.filePath
}

// isDir helps to recognize dir by path
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
