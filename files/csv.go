package files

import (
	"bytes"
	"path"
)

// BOM is the pattern of BOM bytes that can contains CSV file
var BOM = []byte{239, 187, 191}

// removeBOM removes BOM bytes pattern from the file
func removeBOM(records []string) {
	if len(records) != 0 {
		firstRecord := []byte(records[0])
		if bytes.Equal(BOM, firstRecord[:3]) {
			records[0] = string(firstRecord[3:])
		}
	}
}

// isCSV makes sure file is CSV
func isCSV(filePath string) bool {
	return path.Ext(filePath) == ".csv"
}
