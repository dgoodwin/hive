package migration

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up20210129105627, Down20210129105627)
}

func Up20210129105627(tx *sql.Tx) error {
	_, err := tx.Exec(`CREATE TABLE clustersyncs (
id serial NOT NULL PRIMARY KEY,
	name varchar(255) NOT NULL,
	namespace varchar(255) NOT NULL,
	data jsonb NOT NULL,
	unique (name, namespace)
);`)
	if err != nil {
		return err
	}
	return nil
}

func Down20210129105627(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE clustersyncs;`)
	if err != nil {
		return err
	}
	return nil
}
