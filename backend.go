package logger

type backend interface {
	Close()
	Write(int, string, string, ...interface{})
}
