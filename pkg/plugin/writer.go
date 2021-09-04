package plugin

import (
	"io"

	"arhat.dev/pkg/log"
)

var _ io.Writer = (*logWriter)(nil)

type logWriter struct {
	msg    string
	key    string
	logger log.Interface
}

func (w *logWriter) Write(p []byte) (n int, err error) {
	w.logger.D(w.msg, log.String(w.key, string(p)))
	return len(p), nil
}

func LogWriter(logger log.Interface, msg, key string) io.Writer {
	return &logWriter{
		msg:    msg,
		key:    key,
		logger: logger,
	}
}
