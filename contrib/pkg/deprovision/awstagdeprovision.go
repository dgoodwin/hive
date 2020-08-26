package deprovision

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/openshift/installer/pkg/destroy/aws"

	"github.com/openshift/hive/pkg/constants"
)

// NewDeprovisionAWSWithTagsCommand is the entrypoint to create the 'aws-tag-deprovision' subcommand
// TODO: Port to a sub-command of deprovision.
func NewDeprovisionAWSWithTagsCommand() *cobra.Command {
	opt := &aws.ClusterUninstaller{}
	var logLevel string
	cmd := &cobra.Command{
		Use:   "aws-tag-deprovision KEY=VALUE ...",
		Short: "Deprovision AWS assets (as created by openshift-installer) with the given tag(s)",
		Long:  "Deprovision AWS assets (as created by openshift-installer) with the given tag(s).  A resource matches the filter if any of the key/value pairs are in its tags.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := completeAWSUninstaller(opt, logLevel, args); err != nil {
				log.WithError(err).Error("Cannot complete command")
				return
			}
			// Parse credentials from files mounted in as a secret volume and set as env vars. We use the file
			// instead of passing env vars directly so we can monitor for changes and restart the pod if an admin
			// modifies the creds secret after the uninstall job has launched.
			data, err := ioutil.ReadFile(filepath.Join(constants.AWSCredsMount, "aws_access_key_id"))
			if err != nil {
				log.WithError(err).Fatal("error reading aws_access_key_id file")
			}
			os.Setenv("AWS_ACCESS_KEY_ID", string(data))
			data, err = ioutil.ReadFile(filepath.Join(constants.AWSCredsMount, "aws_secret_access_key"))
			if err != nil {
				log.WithError(err).Fatal("error reading aws_secret_access_key file")
			}
			os.Setenv("AWS_SECRET_ACCESS_KEY", string(data))

			// creates a new file watcher
			watcher, err := fsnotify.NewWatcher()
			if err != nil {
				log.WithError(err).Fatal("error creating fsnotify watcher")
			}
			defer watcher.Close()
			done := make(chan bool)
			go func() {
				for {
					select {
					case event := <-watcher.Events:
						log.WithField("event", event).Fatalf("file changes detected in %s, restarting pod", constants.AWSCredsMount)
					case err := <-watcher.Errors:
						log.WithError(err).Fatalf("%s file watch error, restarting pod", constants.AWSCredsMount)
					}
				}
			}()

			// out of the box fsnotify can watch a single file, or a single directory
			if err := watcher.Add("/etc/aws-creds"); err != nil {
				log.WithError(err).Fatal("error establishing watch on /etc/aws-creds")
			}

			if err := opt.Run(); err != nil {
				log.WithError(err).Fatal("Runtime error")
			}
			<-done
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&logLevel, "loglevel", "info", "log level, one of: debug, info, warn, error, fatal, panic")
	flags.StringVar(&opt.Region, "region", "us-east-1", "AWS region to use")
	return cmd
}

func completeAWSUninstaller(o *aws.ClusterUninstaller, logLevel string, args []string) error {

	for _, arg := range args {
		filter := aws.Filter{}
		err := parseFilter(filter, arg)
		if err != nil {
			return fmt.Errorf("cannot parse filter %s: %v", arg, err)
		}
		o.Filters = append(o.Filters, filter)
	}

	// Set log level
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.WithError(err).Error("cannot parse log level")
		return err
	}

	o.Logger = log.NewEntry(&log.Logger{
		Out: os.Stdout,
		Formatter: &log.TextFormatter{
			FullTimestamp: true,
		},
		Hooks: make(log.LevelHooks),
		Level: level,
	})

	return nil
}

func parseFilter(filterMap aws.Filter, str string) error {
	parts := strings.SplitN(str, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("incorrectly formatted filter")
	}

	filterMap[parts[0]] = parts[1]

	return nil
}
