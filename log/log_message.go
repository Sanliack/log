package log

type LogMessage struct {
	Level    int
	Message  string
	Warnflag bool
}

func CreatLogMessage(level int, msg string) *LogMessage {
	logmes := LogMessage{
		level,
		msg,
		false,
	}
	if level >= LevelWarn {
		logmes.Warnflag = true
	}
	return &logmes
}
