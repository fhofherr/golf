package log

type ErrorWriter struct {
	Err error
}

func (w ErrorWriter) Write(p []byte) (n int, err error) {
	return 0, w.Err
}
