package csv

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path"
)

// BOM is the pattern of BOM bytes that can contains CSV file
var BOM = []byte{239, 187, 191}

// removeBOM removes BOM bytes pattern from the file
func removeBOM(records [][]string) {
	if len(records) != 0 {
		firstRecord := []byte(records[0][0])
		if bytes.Equal(BOM, firstRecord[:3]) {
			records[0][0] = string(firstRecord[3:])
		}
	}
}

//ReadCSV reads CSV file by the path and returns two-dimensional array of raw data
func ReadCSV(filePath string, firstIsHeader bool) ([][]string, map[string]struct{}, error) {
	if !isCSV(filePath) {
		return nil, nil, nil
	}
	b, err := os.Open(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed read csv file from %v: %w", filePath, err)
	}
	r := csv.NewReader(b)
	records, err := r.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("failed parse csv from %v: %w", filePath, err)
	}
	removeBOM(records)
	if firstIsHeader {
		return records, getHeaders(records), nil
	}
	return records, nil, nil
}

// getHeaders gets the first line of file
func getHeaders(records [][]string) map[string]struct{} {
	headers := make(map[string]struct{})
	for i, row := range records {
		if i != 0 {
			break
		}
		for _, h := range row {
			headers[h] = struct{}{}
		}
	}
	return headers
}

// isCSV makes sure file is CSV
func isCSV(filePath string) bool {
	return path.Ext(filePath) == ".csv"
}

func NewReader(r io.Reader) *csv.Reader {
	return csv.NewReader(r)
}
