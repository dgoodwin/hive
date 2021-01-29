package hive

import (
	"context"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/pressly/goose"
	log "github.com/sirupsen/logrus"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	_ "github.com/openshift/hive/migration"
	hivev1 "github.com/openshift/hive/pkg/apis/hive/v1"
	"github.com/openshift/hive/pkg/operator/util"
	"github.com/openshift/hive/pkg/resource"
)

const (
	postgresqlDatabaseSecretKey = "POSTGRESQL_DATABASE"
	postgresqlUserSecretKey     = "POSTGRESQL_USER"
	postgresqlPasswordSecretKey = "POSTGRESQL_PASSWORD"
	postgresqlHostSecretKey     = "POSTGRESQL_HOST"
)

func (r *ReconcileHiveConfig) configurePostgresql(hLog log.FieldLogger, hiveNSName string, h resource.Helper, hiveCfg *hivev1.HiveConfig) error {

	// Deploy postgresql pod itself. In future we may support connecting to a separate externally
	// managed db, but for now we assume to make one for ourselves.
	namespacedAssets := []string{
		"config/postgresql/config-secret.yaml",
		"config/postgresql/pvc.yaml",
		"config/postgresql/deployment.yaml",
		"config/postgresql/svc.yaml",
	}
	for _, assetPath := range namespacedAssets {
		if err := util.ApplyAssetWithNSOverrideAndGC(h, assetPath, hiveNSName, hiveCfg); err != nil {
			hLog.WithError(err).Error("error applying postgresql asset with namespace override")
			return err
		}
		hLog.WithField("asset", assetPath).Info("applied asset with namespace override")
	}

	postgresParams, err := r.getPostgresParams(hLog, hiveCfg)
	if err != nil {
		hLog.WithError(err).Error("error getting postgres connection params")
	}

	// TODO: Remove this. Seriously
	// ##########################################################################################
	hLog.WithField("dbStr", postgresParams).Info("TODO WARNING REMOVE THIS REALLY")

	// Postgresql will take a few seconds to come up, until then we error and re-reconcile:
	hLog.Info("connecting to postgres db for migrations")
	db, err := goose.OpenDBWithDriver("postgres", postgresParams)
	if err != nil {
		hLog.WithError(err).Error("unable to connect to postgresql (may not be ready yet)")
		return err
	}
	hLog.Infof("database connection established: %v", db)

	defer func() {
		if err := db.Close(); err != nil {
			log.WithError(err).Error("error closing db connection")
		}
	}()

	hLog.Info("running db migrations")
	err = goose.Run("up", db, ".", []string{}...) // dir is unused for us as we use go migrations
	if err != nil {
		hLog.WithError(err).Error("error running db migrations")
		return err
	}
	hLog.Info("postgresql schema up to date")

	return nil
}

func (r *ReconcileHiveConfig) getPostgresParams(hLog log.FieldLogger, hiveCfg *hivev1.HiveConfig) (string, error) {
	hiveNS := getHiveNamespace(hiveCfg)
	secretName := types.NamespacedName{
		// TODO: hookup properly with configurable secret
		//	Name:      hiveCfg.Spec.StorageBackend.PostgreSQL.CredentialsSecretRef.Name,
		Name:      "hive-postgres-config",
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

	return postgresParams, nil
}

func getSecretKey(sec *corev1.Secret, key string) (string, error) {
	v := sec.Data[key]
	if len(v) == 0 {
		return "", fmt.Errorf("postgresql config Secret %s missing key %s", sec.Name, key)
	}
	return string(v), nil
}
