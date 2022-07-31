package csvalidator

type Config struct {
	// FilePath is path to file or directory with CSV files which to be checked
	FilePath string
	// FirstIsHeader true if first line of CSV is a header, else false
	FirstIsHeader bool
	// SchemaPath path to schema
	SchemaPath string

	// Comma is the field delimiter.
	Comma rune
	// If LazyQuotes is true, a quote may appear in an unquoted field and a
	// non-doubled quote may appear in a quoted field.
	LazyQuotes bool

	// WorkerPoolSize is amount of workers which will be in pool ready to check a file.
	WorkerPoolSize int
}
