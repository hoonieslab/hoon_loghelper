package loghelp

var logFilePath *string

func SetLogFilePath(path string) {
	logFilePath = &path
}
