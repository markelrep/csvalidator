package config

type Config struct {
	// FilePath is path to file or directory with CSV files which to be checked
	FilePath string `short:"p" long:"path" description:"Path to CSV files"`
	// FirstIsHeader true if first line of CSV is a header, else false
	FirstIsHeader bool `short:"f" long:"firstheader" description:"First is header true or false"`
	// SchemaPath path to schema
	SchemaPath string `short:"s" long:"schema" description:"Path to schema"`

	// Comma is the field delimiter.
	Comma rune
	// CommaString uses to define delimiter in CLI
	CommaString string `short:"c" long:"comma" description:"CSV delimiter"`
	// If LazyQuotes is true, a quote may appear in an unquoted field and a
	// non-doubled quote may appear in a quoted field.
	LazyQuotes bool `short:"l" long:"lazyquotes" description:"Lazy Quotes"`

	// WorkerPoolSize is amount of workers which will be in pool ready to check a file.
	WorkerPoolSize int `short:"w" long:"workers" description:"Worker pool size"`

	ErrFilePath string `short:"e" long:"errpath" description:"Path to the error file"`
}
