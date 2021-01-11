// This tools find *.up.sql files and run migrations in order.
// TODO[#5]: Move code for migration to another package.
// TODO[#6]: Add migration by specific path.
package main

import (
	"io/ioutil"
	"sort"
	"strconv"
	"strings"

	"alteroSmartTestTask/common/database"
)

func main() {
	postgresConnection := database.MustGetNewPostgresConnectionUseFlags()
	listOfMigration := readAllFilesWithUpMigration()
	for _, migration := range listOfMigration {
		_, err := postgresConnection.Exec(migration)
		if err != nil {
			panic(err)
		}
	}
}

func readAllFilesWithUpMigration() (output []string) {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		panic(err)
	}
	sort.Slice(files, func(i, j int) bool {
		if !strings.HasSuffix(files[i].Name(), ".up.sql") {
			return i > j
		}
		if !strings.HasSuffix(files[j].Name(), ".up.sql") {
			return i < j
		}
		idI, err := strconv.Atoi(strings.Split(files[i].Name(), "_")[0])
		if err != nil {
			panic(err)
		}
		idJ, err := strconv.Atoi(strings.Split(files[j].Name(), "_")[0])
		if err != nil {
			panic(err)
		}
		return idI < idJ
	})
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".up.sql") {
			bs, err := ioutil.ReadFile(file.Name())
			if err != nil {
				panic(err)
			}
			output = append(output, string(bs))
		}
	}
	return output
}
