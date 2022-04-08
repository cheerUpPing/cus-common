package log

import (
	"fmt"
	cus_common "github.com/cheerUpPing/cus-common"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type logger struct {
	fileLogger    zerolog.Logger
	consoleLogger zerolog.Logger
}

var log = logger{}

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
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
	val, ex := os.LookupEnv(cus_common.LOG_PATH)
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
	logFilePath := logDir + string(os.PathSeparator) + cus_common.LOG_FILE_NAME_PREFIX
	logFile, err := rotatelogs.New(
		logFilePath+cus_common.LOG_TIME_FORMAT,
		rotatelogs.WithLinkName(logFilePath),
		rotatelogs.WithMaxAge(time.Hour*24*365),
		rotatelogs.WithRotationTime(time.Minute*1))
	if err != nil {
		log.consoleLogger.Error().Msg("create symlink fail, err: " + fmt.Sprint(err))
		os.Exit(1)
	}
	fileLogger := zerolog.New(logFile).With().Timestamp().Stack().CallerWithSkipFrameCount(2).Logger()
	log.fileLogger = fileLogger
}

func consoleLogger() {
	console := zerolog.ConsoleWriter{
		Out:          os.Stdout,
		NoColor:      false,
		TimeFormat:   "",
		PartsOrder:   nil,
		PartsExclude: nil,
		FormatTimestamp: func(i interface{}) string {
			return time.Now().Format(cus_common.TIME_FORMAT)
		},
		FormatLevel: func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("| %-6s", i))
		},
		FormatCaller: nil,
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("| %s ", i)
		},
		FormatFieldName: func(i interface{}) string {
			return fmt.Sprintf("| %s: ", i)
		},
		FormatFieldValue: func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("| %s ", i))
		},
		FormatErrFieldName:  nil,
		FormatErrFieldValue: nil,
	}
	consoleLogger := zerolog.New(console).With().Timestamp().Stack().CallerWithSkipFrameCount(2).Logger()
	log.consoleLogger = consoleLogger
}
