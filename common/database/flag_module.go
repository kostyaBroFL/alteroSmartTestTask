package database

import (
	"flag"
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
	"strconv"
	"strings"
)

var defaultHost = "localhost"

var databaseHostEnvName = "POSTGRES_HOST"
var databaseHostFlag = flag.String(
	"postgres_host",
	defaultHost,
	"This is the host for connection to database.",
)

var databasePortEnvName = "POSTGRES_PORT"
var databasePortFlag = flag.Int(
	"postgres_port",
	0,
	"This is the port of database for connection.",
)

var databaseUsernameEnvName = "POSTGRES_USERNAME"
var databaseUsernameFlag = flag.String(
	"postgres_username",
	"",
	"This is the username for log in to the database.",
)

var databasePasswordEnvName = "POSTGRES_PASSWORD"
var databasePasswordFlag = flag.String(
	"postgres_password",
	"",
	"This is the password for log in to the database.",
)

var databaseNameEnvName = "POSTGRES_DATABASE_NAME"
var databaseNameFlag = flag.String(
	"postgres_database_name",
	"",
	"This is the name of the database in which the data is stored.",
)

func MustGetNewPostgresConnectionUseFlags() *sqlx.DB {
	databaseHost := *databaseHostFlag
	if 0 == strings.Compare(databaseHost, defaultHost) {
		databaseHostEnv := os.Getenv(databaseHostEnvName)
		if "" == databaseHostEnv {
			fmt.Printf("Database host not set, so using %s.\n", defaultHost)
		} else {
			databaseHost = databaseHostEnv
		}
	}
	databasePort := int64(*databasePortFlag)
	if 0 == databasePort {
		databasePortEnv := os.Getenv(databasePortEnvName)
		if "" == databasePortEnv {
			panic("Database port must be set.")
		}
		var err error
		databasePort, err = strconv.ParseInt(databasePortEnv, 10, 32)
		if err != nil {
			panic(fmt.Sprintf("Database port env parse error. %s.", err.Error()))
		}
	}
	databaseUsername := *databaseUsernameFlag
	if "" == databaseUsername {
		databaseUsername = os.Getenv(databaseUsernameEnvName)
		if "" == databaseUsername {
			panic("Database username must be set.")
		}
	}
	databasePassword := *databasePasswordFlag
	if "" == databasePassword {
		databasePassword = os.Getenv(databasePasswordEnvName)
		if "" == databasePassword {
			panic("Database password must be set.")
		}
	}
	databaseName := *databaseNameFlag
	if "" == databaseName {
		databaseName = os.Getenv(databaseNameEnvName)
		if "" == databaseName {
			panic("Database name must be set.")
		}
	}
	return MustGetNewPostgresConnection(
		databaseHost,
		int(databasePort),
		databaseUsername,
		databasePassword,
		databaseName,
	)
}
