package logger

import "time"

// LogArg параметр для дополнительного логирования
type LogArg struct {
	argType string
	name    string
	value   any
}

func StrArg(key string, val string) LogArg {
	return LogArg{argType: "string", name: key, value: val}
}

func IntArg(key string, val int) LogArg {
	return LogArg{argType: "int", name: key, value: val}
}

func TimeDurationArg(key string, val time.Duration) LogArg {
	return LogArg{argType: "duration", name: key, value: val}
}
