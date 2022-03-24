package log

const (
	LevelDebug = iota
	LevelTrace
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var Levelmap = map[int]string{
	LevelDebug: "Debug",
	LevelTrace: "Trace",
	LevelInfo:  "Info",
	LevelWarn:  "Warn",
	LevelError: "Error",
	LevelFatal: "Fatal",
}
