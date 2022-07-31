package files

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"

	"github.com/markelrep/csvalidator/config"
)

// File represent CSV file
type File struct {
	stream     chan Row
	config     config.Config
	headersMap map[string]struct{}
	headers    []string
	headersLen int
	path       string
}

// Row represents row from csv file as slice of string
type Row struct {
	Data  []string
	Index int
}

func NewFile(path string, config config.Config) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	ch := make(chan Row)
	csvReader := csv.NewReader(f)
	csvReader.LazyQuotes = config.LazyQuotes
	if config.Comma != 0 {
		csvReader.Comma = config.Comma
	}

	var headers []string
	headersMap := make(map[string]struct{})
	if config.FirstIsHeader {
		line, err := csvReader.Read()
		removeBOM(line)
		if err != nil {
			return nil, err
		}
		for _, l := range line {
			headersMap[l] = struct{}{}
			headers = append(headers, l)
		}
	}

	fileStream := &File{
		stream:     ch,
		config:     config,
		path:       path,
		headersMap: headersMap,
		headers:    headers,
		headersLen: len(headers),
	}

	go fileStream.run(csvReader)

	return fileStream, nil
}

func (f *File) run(reader *csv.Reader) {
	var line int
	for {
		data, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				close(f.stream)
				return
			}
			log.Println(err)
			close(f.stream)
			return
		}
		if f.config.FirstIsHeader && line == 0 {
			line++
		}
		if !f.config.FirstIsHeader && line == 0 {
			removeBOM(data)
		}
		f.stream <- Row{Data: data, Index: line}
		line++
	}
}

// Stream returns chanel of row.
func (f *File) Stream() chan Row {
	return f.stream
}

// HasHeader checks if specific header exists or not
func (f *File) HasHeader(value string) bool {
	_, ok := f.headersMap[value]
	return ok
}

// FirstIsHeader gets information about first line, is it header or not
func (f *File) FirstIsHeader() bool {
	return f.config.FirstIsHeader
}

// Path returns path to file
func (f *File) Path() string {
	return f.path
}

// HeadersCount return len of headers map
func (f *File) HeadersCount() int {
	return f.headersLen
}

// Headers returns list of headers
func (f *File) Headers() []string {
	return f.headers
}

// Files slice of File
type Files []*File

// NewFiles create a new Files from path
func NewFiles(config config.Config) (Files, error) {
	path := config.FilePath
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
			f, err := NewFile(fullPath, config)
			if err != nil {
				return err
			}
			files = append(files, f)
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("failed open %s: %w", path, err)
		}
		return files, nil
	}

	f, err := NewFile(path, config)
	if err != nil {
		return nil, err
	}
	files = append(files, f)
	return files, nil
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
