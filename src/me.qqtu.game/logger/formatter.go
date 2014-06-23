package logger

type Formatter interface {
	Format(event LogEvent) ([]byte, error)
}
