package logger

type Logger interface {
	Error()
	Warnf()
	Infof()
	Debugf()
}

type WebberLogger struct {
	Logger Logger
}

func (wl *WebberLogger) Error() {
	wl.Logger.Error()
}

func (wl *WebberLogger) Warnf() {
	wl.Logger.Warnf()
}

func (wl *WebberLogger) Infof() {
	wl.Logger.Infof()
}

func (wl *WebberLogger) Debugf() {
	wl.Logger.Debugf()
}
