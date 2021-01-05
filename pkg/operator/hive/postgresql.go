package hive

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	hivev1 "github.com/openshift/hive/pkg/apis/hive/v1"
)

const (
	postgresqlDatabaseSecretKey = "POSTGRESQL_DATABASE"
	postgresqlUserSecretKey     = "POSTGRESQL_USER"
	postgresqlPasswordSecretKey = "POSTGRESQL_PASSWORD"
	postgresqlHostSecretKey     = "POSTGRESQL_HOST"
)

func (r *ReconcileHiveConfig) configurePostgresql(hLog log.FieldLogger, hiveCfg *hivev1.HiveConfig) error {
	hLog.Info("configuring postgresql storage")

	postgresParams, err := r.getPostgresParams(hLog, hiveCfg)

	db, err := sql.Open("postgres", postgresParams)
	if err != nil {
		log.WithError(err).Fatal("error connecting to database")
	}
	log.Infof("database connection established: %v", db)

	hLog.Info("postgresql storage configured")
	return nil
}

func (r *ReconcileHiveConfig) getPostgresParams(hLog log.FieldLogger, hiveCfg *hivev1.HiveConfig) (string, error) {
	hiveNS := getHiveNamespace(hiveCfg)
	secretName := types.NamespacedName{
		Name:      hiveCfg.Spec.StorageBackend.PostgreSQL.CredentialsSecretRef.Name,
		Namespace: hiveNS,
	}
	postgresConfigSecret := &corev1.Secret{}
	err := r.Client.Get(context.Background(), secretName, postgresConfigSecret)
	if err != nil {
		return "", errors.Wrap(err, "error looking up postgresql config secret")
	}

	host, err := getSecretKey(postgresConfigSecret, postgresqlHostSecretKey)
	if err != nil {
		return "", err
	}

	user, err := getSecretKey(postgresConfigSecret, postgresqlUserSecretKey)
	if err != nil {
		return "", err
	}
	hLog.Infof("got postgres user: %s", user)

	password, err := getSecretKey(postgresConfigSecret, postgresqlPasswordSecretKey)
	if err != nil {
		return "", err
	}

	dbName, err := getSecretKey(postgresConfigSecret, postgresqlDatabaseSecretKey)
	if err != nil {
		return "", err
	}

	postgresParams := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbName)
	// TODO: definitely remove this
	hLog.WithField("postgresParams", postgresParams).Info("built postgres params string")

	return "", nil
}

func getSecretKey(sec *corev1.Secret, key string) (string, error) {
	v := sec.Data[key]
	if len(v) == 0 {
		return "", fmt.Errorf("postgresql config Secret %s missing key %s", sec.Name, key)
	}
	return string(v), nil
}
