package csvalidator

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
)

var BOM = []byte{239, 187, 191}

func removeBOM(records [][]string) {
	if len(records) != 0 {
		firstRecord := []byte(records[0][0])
		if bytes.Compare(BOM, firstRecord[:3]) == 0 {
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
