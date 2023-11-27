package logging

import "go.uber.org/zap"

type ElasticLogger struct {
	ecsLogger *zap.Logger
}

func NewElasticLogger(ecsLogger *zap.Logger) Logger {

	return &ElasticLogger{
		ecsLogger: ecsLogger,
	}

}

func (log *ElasticLogger) Info(msg string) {
	log.ecsLogger.Info(msg)
}

func (log *ElasticLogger) Debug(msg string) {
	log.ecsLogger.Debug(msg)
}

func (log *ElasticLogger) Warn(msg string) {
	log.ecsLogger.Warn(msg)
}

func (log *ElasticLogger) Error(msg string) {
	log.ecsLogger.Error(msg)
}
