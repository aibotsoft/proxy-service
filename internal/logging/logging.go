package logging

import (
	"log"
	"os"
	"proxy-service/internal/config"
)

//func getLevel(level string) zapcore.Level {
//	switch strings.ToLower(strings.TrimSpace(level)) {
//	case "debug":
//		return zap.DebugLevel
//	case "info":
//		return zap.InfoLevel
//	case "error":
//		return zap.ErrorLevel
//	default:
//		return zap.InfoLevel
//	}
//}

func New(cfg *config.Config) *log.Logger {
	prefix := cfg.Service.Name + ": "
	return log.New(os.Stdout, prefix, log.Ldate|log.Lmsgprefix|log.Lmicroseconds|log.Lshortfile)
}

//func New() (*zap.Logger, error) {
//	return zap.NewProduction()
//}
//func New(c *config.Config) (*zap.Logger, error) {
//	zapConfig := zap.Config{
//		//Level:             getLevel(c.Logging.Level),
//		Development:       true,
//		DisableStacktrace: true,
//		Encoding:          "console",
//		EncoderConfig:     zap.NewDevelopmentEncoderConfig(),
//		OutputPaths:       []string{"stderr"},
//		ErrorOutputPaths:  []string{"stderr"},
//	}
//	zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
//
//	return zapConfig.Build()
//}
