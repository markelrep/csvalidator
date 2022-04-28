package csvalidator

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"path"
)

var BOM = []byte{239, 187, 191}

func removeBOM(records [][]string) {
	if len(records) != 0 {
		firstRecord := []byte(records[0][0])
		if bytes.Equal(BOM, firstRecord[:3]) {
			records[0][0] = string(firstRecord[3:])
		}
	}
}

func readCSV(filePath string) ([][]string, error) {
	b, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed read csv file from %v: %w", filePath, err)
	}
	r := csv.NewReader(b)
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed parse csv from %v: %w", filePath, err)
	}
	removeBOM(records)
	return records, nil
}

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

func isCSV(filePath string) bool {
	return path.Ext(filePath) == ".csv"
}
