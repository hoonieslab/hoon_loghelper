package loghelp

var logFilePath *string

// SetLogFilePath set path to save log
func SetLogFilePath(path string) {
	logFilePath = &path
}
