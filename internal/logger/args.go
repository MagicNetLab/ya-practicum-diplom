package logger

import "time"

// LogArg параметр для дополнительного логирования
type LogArg struct {
	argType string
	name    string
	value   any
}

// StrArg создает аргумент типа string
func StrArg(key string, val string) LogArg {
	return LogArg{argType: "string", name: key, value: val}
}

// IntArg создает аргумент типа int
func IntArg(key string, val int) LogArg {
	return LogArg{argType: "int", name: key, value: val}
}

// TimeDurationArg создает аргумент типа time.Duration
func TimeDurationArg(key string, val time.Duration) LogArg {
	return LogArg{argType: "duration", name: key, value: val}
}
