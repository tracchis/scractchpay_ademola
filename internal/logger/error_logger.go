package logger

// ErrorLoggerFunc implements the ErrorLogger interface
type ErrorLoggerFunc func(error)

// Error logs an error message
func (f ErrorLoggerFunc) Error(err error) {
	if f != nil {
		f(err)
	}
}
