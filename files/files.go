package files

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/markelrep/csvalidator/csv"
)

type File struct {
	filePath      string
	firstIsHeader bool
	headers       map[string]struct{}
	headersLen    int
	Records       [][]string
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
			fullPath := path + "/" + p
			if !csv.IsCSV(fullPath) {
				return nil
			}
			records, err := csv.ReadCSV(fullPath)
			if err != nil {
				return err
			}
			h := csv.GetHeaders(records)
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

	records, err := csv.ReadCSV(path)
	if err != nil {
		return nil, fmt.Errorf("failed create validator: %w", err)
	}
	h := csv.GetHeaders(records)
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

func (f File) FirstIsHeader() bool {
	return f.firstIsHeader
}

func (f File) HeadersCount() int {
	return f.headersLen
}

func (f File) HasHeader(name string) bool {
	_, ok := f.headers[name]
	return ok
}

func (f File) Path() string {
	return f.filePath
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
