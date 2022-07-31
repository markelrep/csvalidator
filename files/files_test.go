package files

import (
	"fmt"
	"testing"

	"github.com/markelrep/csvalidator/config"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestTest(t *testing.T) {
	fs, err := NewFile("../samples/file.csv", config.Config{FirstIsHeader: true})
	require.NoError(t, err)
	for r := range fs.Stream() {
		fmt.Println(r.Index, ":", r.Data)
	}
}

func TestNewFile(t *testing.T) {
	_, err := NewFile("../samples/file.csv", config.Config{FirstIsHeader: true})
	require.NoError(t, err)
}

func TestFile_FirstIsHeader(t *testing.T) {
	f, err := NewFile("../samples/file.csv", config.Config{FirstIsHeader: true})
	require.NoError(t, err)
	assert.Equal(t, true, f.config.FirstIsHeader)
}

func TestFile_Path(t *testing.T) {
	f, err := NewFile("../samples/file.csv", config.Config{FirstIsHeader: true})
	require.NoError(t, err)
	assert.Equal(t, "../samples/file.csv", f.Path())
}

func TestFile_HasHeader(t *testing.T) {
	f, err := NewFile("../samples/file.csv", config.Config{FirstIsHeader: true})
	require.NoError(t, err)

	cases := []struct {
		column   string
		expected bool
	}{
		{
			column:   "id",
			expected: true,
		},
		{
			column:   "qq",
			expected: false,
		},
	}

	for _, tc := range cases {
		assert.Equal(t, tc.expected, f.HasHeader(tc.column))
	}
}

func TestFile_HeadersCount(t *testing.T) {
	f, err := NewFile("../samples/file.csv", config.Config{FirstIsHeader: true})
	require.NoError(t, err)
	assert.Equal(t, 2, f.HeadersCount())
}

func TestFile_Headers(t *testing.T) {
	f, err := NewFile("../samples/file.csv", config.Config{FirstIsHeader: true})
	require.NoError(t, err)
	assert.Equal(t, []string{"id", "comment"}, f.Headers())
}
