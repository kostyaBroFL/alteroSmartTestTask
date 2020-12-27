package log

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

var enableJsonEnvName = "LOG_AS_JSON"
var enableJsonFlag = flag.Bool(
	"log_as_json",
	false,
	"This flag enable logging in json mode.",
)

var logLevelEnvName = "LOG_LEVEL"
var logLevelFlag = flag.Int(
	"log_level",
	int(logrus.DebugLevel),
	"This flag for setup logging deep.",
)

func ProvideEnableJsonFlag() (enableJson bool) {
	enableJson = *enableJsonFlag
	if !enableJson {
		enableJsonEnv := os.Getenv(enableJsonEnvName)
		if "" == enableJsonEnv ||
			0 == strings.Compare("false", strings.ToLower(enableJsonEnv)) {
			return false
		} else if 0 == strings.Compare("true", strings.ToLower(enableJsonEnv)) {
			return true
		} else {
			panic(fmt.Sprintf("%s must be true or false", enableJsonEnvName))
		}
	}
	return enableJson
}

func ProvideLogLevelFlag() logrus.Level {
	logLevel := int64(*logLevelFlag)
	if 0 == logLevel {
		logLevelEnv := os.Getenv(logLevelEnvName)
		if "" == logLevelEnv {
			panic(fmt.Sprintf("need to provide %s", logLevelEnvName))
		}
		var err error
		logLevel, err = strconv.ParseInt(logLevelEnv, 10, 32)
		if err != nil {
			panic(fmt.Sprintf(
				"%s env parse error. %s.",
				logLevelEnvName, err.Error(),
			))
		}
	}
	logrus.NewEntry(logrus.New()).
		WithField("log_level_from_flags", logLevel).
		Info(logrus.Level(logLevel).String())
	return logrus.Level(logLevel)
}

func ProvideLogrusLogger(
	logAsJson bool,
	logLevel logrus.Level,
) *logrus.Logger {
	logger := logrus.New()
	if logAsJson {
		logger.Formatter = &logrus.JSONFormatter{}
	}
	logger.SetLevel(logLevel)
	return logger
}

func ProvideLogrusLoggerUseFlags() *logrus.Logger {
	return ProvideLogrusLogger(
		ProvideEnableJsonFlag(),
		ProvideLogLevelFlag(),
	)
}

// TODO[#371](https://team.cron.global/issue/wallet-371): create DI methods
