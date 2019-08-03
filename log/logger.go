package log

import "go.uber.org/zap"

var (
	logger, _     = zap.NewProduction()
	sugaredLogger = logger.Sugar()
)

func GetLogger() *zap.SugaredLogger {
	return sugaredLogger
}
