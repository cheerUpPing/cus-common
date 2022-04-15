package c_log

import (
	"fmt"
	c_common "github.com/cheerUpPing/cus-common"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"os"
	"path/filepath"
	"time"
)

type logger struct {
	fileLogger    zerolog.Logger
	consoleLogger zerolog.Logger
}

var log = logger{}

func init() {
	time.LoadLocation("Asia/Shanghai") //上海
	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "l"
	zerolog.MessageFieldName = "m"
	zerolog.CallerFieldName = "c"
	zerolog.TimeFieldFormat = c_common.TIME_FORMAT
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	consoleLogger()
	fileLogger()
}

func LogInfo(traceId, msg string) {
	log.consoleLogger.Info().CallerSkipFrame(1).Msg(fmt.Sprintf("%s | %s", traceId, msg))
	log.fileLogger.Info().CallerSkipFrame(1).Msg(fmt.Sprintf("%s | %s", traceId, msg))
}

func LogError(traceId string, err error) {
	log.consoleLogger.Error().Stack().Err(errors.WithStack(err)).Msg(fmt.Sprintf("%s | %s", traceId, err))
	log.fileLogger.Error().Stack().Err(errors.WithStack(err)).Msg(fmt.Sprintf("%s | %s", traceId, err))
}

func fileLogger() {
	val, ex := os.LookupEnv(c_common.LOG_PATH)
	if !ex {
		log.consoleLogger.Error().Msg("LOG_PATH env not config")
		os.Exit(1)
	}
	_, err := os.Stat(val)
	if err != nil {
		err = os.MkdirAll(val, os.ModePerm)
		if err != nil {
			log.consoleLogger.Error().Msg("create log dir fail, err: " + fmt.Sprint(err))
			os.Exit(1)
		}
	}
	logDir, err := filepath.Abs(val)
	if err != nil {
		log.consoleLogger.Error().Msg("get log dir path fail, err: " + fmt.Sprint(err))
		os.Exit(1)
	}
	logFilePath := logDir + string(os.PathSeparator) + c_common.LOG_FILE_NAME_PREFIX
	logFile, err := rotatelogs.New(
		logFilePath+c_common.LOG_TIME_FORMAT,
		rotatelogs.WithLinkName(logFilePath),
		rotatelogs.WithMaxAge(time.Hour*24*90),
		rotatelogs.WithRotationTime(time.Hour*24))
	if err != nil {
		log.consoleLogger.Error().Msg("create symlink fail, err: " + fmt.Sprint(err))
		os.Exit(1)
	}
	fileLogger := zerolog.New(logFile).With().Timestamp().Caller().Logger()
	log.fileLogger = fileLogger
}

func consoleLogger() {
	consoleLogger := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
	log.consoleLogger = consoleLogger
}
