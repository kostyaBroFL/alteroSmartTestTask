// In this package stores the functions
// for parsing variables from env and flags.
// Flag more priority.
// For example if variable PORT = 80, but flag -port = 443, port will be 443.
package flagenv

import (
	"os"
	"strconv"
	"strings"

	"alteroSmartTestTask/common/errors"
)

func MustParseInt(flag *int, envName string) int {
	if value, err := ParseInt(flag, envName); err != nil {
		panic(err.Error())
	} else {
		return value
	}
}

func MustParseBool(flag *bool, envName string) bool {
	if value, err := ParseBool(flag, envName); err != nil {
		panic(err.Error())
	} else {
		return value
	}
}

func MustParseString(flag *string, envName string) string {
	if value, err := ParseString(flag, envName); err != nil {
		panic(err.Error())
	} else {
		return value
	}
}

func ParseInt(flag *int, envName string) (int, error) {
	if 0 == *flag {
		envValue := os.Getenv(envName)
		if "" == envValue {
			return 0, errors.Newf("%s is not defined\n", envName)
		}
		expectedValue, err := strconv.Atoi(envValue)
		if 0 == expectedValue {
			return 0, errors.Newf("%s is not defined\n", envName)
		}
		if err != nil {
			return 0, err
		}

		return expectedValue, nil
	}
	return *flag, nil
}

func ParseBool(flag *bool, envName string) (bool, error) {
	if false == *flag {
		envValue := strings.ToLower(os.Getenv(envName))
		if "" == envValue {
			return false, errors.Newf("%s is not defined\n", envName)
		}
		if "true" == envValue ||
			"t" == envValue ||
			"1" == envValue {
			return true, nil
		}
		return false, nil
	}
	return true, nil
}

func ParseString(flag *string, envName string) (string, error) {
	if "" == *flag {
		envValue := os.Getenv(envName)
		if "" == envValue {
			return "", errors.Newf("%s is not defined\n", envName)
		}
		return envValue, nil
	}
	return *flag, nil
}
