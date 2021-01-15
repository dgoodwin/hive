package datastore

import (
	"database/sql"
	//"encoding/json"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"

	hiveinternal "github.com/openshift/hive/pkg/apis/hiveinternal/v1alpha1"
)

// DataStore defines the interface used for additional Hive database storage which is not
// suitable to store in etcd.
type DataStore interface {
	Get(namespace, name string) (*hiveinternal.ClusterSync, error)
	Create(clusterSync *hiveinternal.ClusterSync) error
	Update(clusterSync *hiveinternal.ClusterSync) error
}

type PostgresDataStore struct {
	db *sql.DB
}

// Get attempts to retrieve a ClusterSync for the given namespace and name.
// If the resource does not exist, nil will be returned. (with no error as well)
func (pgds *PostgresDataStore) Get(namespace, name string) (*hiveinternal.ClusterSync, error) {
	cs := &hiveinternal.ClusterSync{}
	row := pgds.db.QueryRow("SELECT data FROM clustersyncs WHERE namespace=$1 and name=$2 ORDER BY id DESC LIMIT 1", namespace, name)
	switch err := row.Scan(cs); err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		log.WithField("clustersync", cs).Info("got cluster sync")
		return cs, nil
	default:
		return nil, err
	}
}

func (pgds *PostgresDataStore) Create(cs *hiveinternal.ClusterSync) error {
	// The database driver will call the Value() method and and marshall the
	// attrs struct to JSON before the INSERT.
	_, err := pgds.db.Exec("INSERT INTO clustersyncs (name, namespace, data) VALUES($1, $2, $3)", cs.Name, cs.Namespace, cs)
	return err
}

func (pgds *PostgresDataStore) Update(cs *hiveinternal.ClusterSync) error {
	// The database driver will call the Value() method and and marshall the
	// attrs struct to JSON before the INSERT.
	_, err := pgds.db.Exec("UPDATE clustersyncs SET data=$3 WHERE namespace=$1 AND name=$2", cs.Name, cs.Namespace, cs)
	return err
}

// BuildDataStore will create a datastore based on the current environment variables.
func BuildDataStore() (DataStore, error) {
	// TODO: set postgres connection string in operator, check if env var is defined, error otherwise
	db, err := sql.Open("postgres", "user=hive dbname=hive sslmode=disable host=postgres password=helloworld")
	if err != nil {
		log.WithError(err).Error("error connecting to postgresql database")
		return nil, err
	} else {
		log.Info("database connection established")
	}
	return &PostgresDataStore{
		db: db,
	}, nil
}
