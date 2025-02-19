package logger

// LogArg параметр для дополнительного логирования
type LogArg struct {
	argType string
	name    string
	value   any
}

func StrArg(key string, val string) LogArg {
	return LogArg{argType: "string", name: key, value: val}
}
