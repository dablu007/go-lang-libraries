package logger

import (
	logrs "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"path"
	"sync"
)

var logger *logrs.Logger
var lock = &sync.Mutex{}
var SugarLogger *zap.SugaredLogger

func InitLogger() {
	writerSyncer := getLogWriter()

	encoder := getEncoder()

	core := zapcore.NewCore(encoder, writerSyncer, zapcore.DebugLevel)

	logger := zap.New(core, zap.AddCaller())

	SugarLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	logFilePath := viper.GetString("log.file.path")
	logFileName := viper.GetString("log.file.name")
	logFileMaxSize := viper.GetInt("log.file.maxsize")
	logFileMaxBackups := viper.GetInt("log.file.maxbackup")
	logFileMaxAge := viper.GetInt("log.file.maxage")
	logFile := path.Join(logFilePath, logFileName)
	lumberJackLogger := &lumberjack.Logger{
		Filename: logFile,
		MaxSize: logFileMaxSize,
		MaxBackups: logFileMaxBackups,
		MaxAge: logFileMaxAge,
		Compress: true,
		LocalTime:true,
	}
	return zapcore.AddSync(lumberJackLogger)
}