package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var (
	logger *zap.Logger
)

func Init() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = customTimeEncoder

	//p, _ := filepath.Abs("./logs/server.log")

	//config.OutputPaths = []string{p}
	//var err error
	logger, _ = config.Build()
	//if err != nil {
	//	config.OutputPaths = nil
	//	logger, _ = config.Build()
	//}
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("02-01-2006 15:04:05")) // DD-MM-YYYY HH:MM:SS
}

func Error(err error) {
	logger.Error("Error: ", zap.Error(err))
}

func Logger() *zap.Logger {
	return logger
}
