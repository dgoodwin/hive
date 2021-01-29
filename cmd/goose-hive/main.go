package main

import (
	"flag"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	log "github.com/sirupsen/logrus"

	_ "github.com/openshift/hive/migration"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	// Unused for Hive, we only use go migrations.
	dir = flags.String("dir", ".", "directory with migration files")
)

func main() {
	flags.Parse(os.Args[1:])
	args := flags.Args()

	if len(args) < 2 {
		flags.Usage()
		return
	}

	arguments := []string{}
	if len(args) > 3 {
		arguments = append(arguments, args[3:]...)
	}

	dbStr, command := args[0], args[1]
	runGooseDBMigrations(dbStr, command, arguments)

}

func runGooseDBMigrations(dbStr, command string, arguments []string) error {
	log.Info("connecting to postgres db for migrations")
	db, err := goose.OpenDBWithDriver("postgres", dbStr)
	if err != nil {
		return err
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.WithError(err).Error("error closing db connection")
		}
	}()

	log.Info("running db migrations")
	err = goose.Run(command, db, *dir, arguments...)
	return err
}
