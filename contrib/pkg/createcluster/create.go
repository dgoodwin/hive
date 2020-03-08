package createcluster

import (
	"encoding/json"
	"fmt"
	awsutils "github.com/openshift/hive/contrib/pkg/utils/aws"
	gcputils "github.com/openshift/hive/contrib/pkg/utils/gcp"
	"github.com/openshift/hive/pkg/cluster"
	"github.com/openshift/hive/pkg/gcpclient"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	"github.com/openshift/hive/pkg/apis"
	hivev1 "github.com/openshift/hive/pkg/apis/hive/v1"
	"github.com/openshift/hive/pkg/resource"
)

const longDesc = `
OVERVIEW
The hiveutil create-cluster command generates and applies the artifacts needed
to create a new Hive cluster deployment. By default, the clusterdeployment is
generated along with corresponding secrets and then applied to the current
cluster. If you don't need secrets generated, specify --include-secrets=false
in the command line. If you don't want to apply the cluster deployment and
only output it locally, specify the output flag (-o json) or (-o yaml) to
specify your output format.

IMAGES
An existing ClusterImageSet can be specified with the --image-set
flag. Otherwise, one will be generated using the images specified for the
cluster deployment. If you don't wish to use a ClusterImageSet, specify
--use-image-set=false. This will result in images only specified on the
cluster itself.


ENVIRONMENT VARIABLES
The command will use the following environment variables for its output:

PUBLIC_SSH_KEY - If present, it is used as the new cluster's public SSH key.
It overrides the public ssh key flags. If not, --ssh-public-key will be used.
If that is not specified, then --ssh-public-key-file is used.
That file's default value is %[1]s.

PULL_SECRET - If present, it is used as the cluster deployment's pull
secret and will override the --pull-secret flag. If not present, and
the --pull-secret flag is not specified, then the --pull-secret-file is
used. That file's default value is %[2]s.

AWS_SECRET_ACCESS_KEY and AWS_ACCESS_KEY_ID - Are used to determine your
AWS credentials. These are only relevant for creating a cluster on AWS. If
--creds-file is used it will take precedence over these environment
variables.

RELEASE_IMAGE - Release image to use to install the cluster. If not specified,
the --release-image flag is used. If that's not specified, a default image is
obtained from a the following URL:
https://openshift-release.svc.ci.openshift.org/api/v1/releasestream/4-stable/latest
`
const (
	hiveutilCreatedLabel = "hive.openshift.io/hiveutil-created"
	cloudAWS             = "aws"
	cloudAzure           = "azure"
	cloudGCP             = "gcp"

	testFailureManifest = `apiVersion: v1
kind: NotARealSecret
metadata:
  name: foo
  namespace: bar
type: TestFailResource
`
	azureCredFile = "osServicePrincipal.json"
)

var (
	validClouds = map[string]bool{
		cloudAWS:   true,
		cloudAzure: true,
		cloudGCP:   true,
	}
)

// Options is the set of options to generate and apply a new cluster deployment
type Options struct {
	Name                     string
	Namespace                string
	SSHPublicKeyFile         string
	SSHPublicKey             string
	SSHPrivateKeyFile        string
	BaseDomain               string
	PullSecret               string
	PullSecretFile           string
	Cloud                    string
	CredsFile                string
	CredsSecret              string
	ClusterImageSet          string
	ReleaseImage             string
	ReleaseImageSource       string
	DeleteAfter              string
	ServingCert              string
	ServingCertKey           string
	UseClusterImageSet       bool
	ManageDNS                bool
	Output                   string
	IncludeSecrets           bool
	InstallOnce              bool
	UninstallOnce            bool
	SimulateBootstrapFailure bool
	WorkerNodes              int64
	CreateSampleSyncsets     bool
	ManifestsDir             string
	Adopt                    bool
	AdoptAdminKubeConfig     string
	AdoptInfraID             string
	AdoptClusterID           string
	AdoptAdminUsername       string
	AdoptAdminPassword       string

	// Azure
	AzureBaseDomainResourceGroupName string

	homeDir string
}

// NewCreateClusterCommand creates a command that generates and applies cluster deployment artifacts.
func NewCreateClusterCommand() *cobra.Command {
	opt := &Options{}

	opt.homeDir = "."

	if u, err := user.Current(); err == nil {
		opt.homeDir = u.HomeDir
	}
	defaultSSHPublicKeyFile := filepath.Join(opt.homeDir, ".ssh", "id_rsa.pub")

	defaultPullSecretFile := filepath.Join(opt.homeDir, ".pull-secret")
	if _, err := os.Stat(defaultPullSecretFile); os.IsNotExist(err) {
		defaultPullSecretFile = ""
	} else if err != nil {
		log.WithError(err).Errorf("%v can not be used", defaultPullSecretFile)
	}

	cmd := &cobra.Command{
		Use: `create-cluster CLUSTER_DEPLOYMENT_NAME
create-cluster CLUSTER_DEPLOYMENT_NAME --cloud=aws
create-cluster CLUSTER_DEPLOYMENT_NAME --cloud=azure --azure-base-domain-resource-group-name=RESOURCE_GROUP_NAME
create-cluster CLUSTER_DEPLOYMENT_NAME --cloud=gcp`,
		Short: "Creates a new Hive cluster deployment",
		Long:  fmt.Sprintf(longDesc, defaultSSHPublicKeyFile, defaultPullSecretFile),
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			log.SetLevel(log.InfoLevel)
			if err := opt.Complete(cmd, args); err != nil {
				return
			}
			if err := opt.Validate(cmd); err != nil {
				return
			}
			err := opt.Run()
			if err != nil {
				log.WithError(err).Error("Error")
				os.Exit(1)
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opt.Cloud, "cloud", cloudAWS, "Cloud provider: aws(default)|azure|gcp)")
	flags.StringVarP(&opt.Namespace, "namespace", "n", "", "Namespace to create cluster deployment in")
	flags.StringVar(&opt.SSHPrivateKeyFile, "ssh-private-key-file", "", "file name containing private key contents")
	flags.StringVar(&opt.SSHPublicKeyFile, "ssh-public-key-file", defaultSSHPublicKeyFile, "file name of SSH public key for cluster")
	flags.StringVar(&opt.SSHPublicKey, "ssh-public-key", "", "SSH public key for cluster")
	flags.StringVar(&opt.BaseDomain, "base-domain", "new-installer.openshift.com", "Base domain for the cluster")
	flags.StringVar(&opt.PullSecret, "pull-secret", "", "Pull secret for cluster. Takes precedence over pull-secret-file.")
	flags.StringVar(&opt.DeleteAfter, "delete-after", "", "Delete this cluster after the given duration. (i.e. 8h)")
	flags.StringVar(&opt.PullSecretFile, "pull-secret-file", defaultPullSecretFile, "Pull secret file for cluster")
	flags.StringVar(&opt.CredsSecret, "creds-secret", "", "Pre-existing cloud credentials secret in the target namespace to use.")
	flags.StringVar(&opt.CredsFile, "creds-file", "", "Cloud credentials file (defaults vary depending on cloud)")
	flags.StringVar(&opt.ClusterImageSet, "image-set", "", "Cluster image set to use for this cluster deployment")
	flags.StringVar(&opt.ReleaseImage, "release-image", "", "Release image to use for installing this cluster deployment")
	flags.StringVar(&opt.ReleaseImageSource, "release-image-source", "https://openshift-release.svc.ci.openshift.org/api/v1/releasestream/4-stable/latest", "URL to JSON describing the release image pull spec")
	flags.StringVar(&opt.ServingCert, "serving-cert", "", "Serving certificate for control plane and routes")
	flags.StringVar(&opt.ServingCertKey, "serving-cert-key", "", "Serving certificate key for control plane and routes")
	flags.BoolVar(&opt.ManageDNS, "manage-dns", false, "Manage this cluster's DNS. This is only available for AWS.")
	flags.BoolVar(&opt.UseClusterImageSet, "use-image-set", true, "If true(default), use a cluster image set for this cluster")
	flags.StringVarP(&opt.Output, "output", "o", "", "Output of this command (nothing will be created on cluster). Valid values: yaml,json")
	flags.BoolVar(&opt.IncludeSecrets, "include-secrets", true, "Include secrets along with ClusterDeployment")
	flags.BoolVar(&opt.InstallOnce, "install-once", false, "Run the install only one time and fail if not successful")
	flags.BoolVar(&opt.UninstallOnce, "uninstall-once", false, "Run the uninstall only one time and fail if not successful")
	flags.BoolVar(&opt.SimulateBootstrapFailure, "simulate-bootstrap-failure", false, "Simulate an install bootstrap failure by injecting an invalid manifest.")
	flags.Int64Var(&opt.WorkerNodes, "workers", 3, "Number of worker nodes to create.")
	flags.BoolVar(&opt.CreateSampleSyncsets, "create-sample-syncsets", false, "Create a set of sample syncsets for testing")
	flags.StringVar(&opt.ManifestsDir, "manifests", "", "Directory containing manifests to add during installation")

	// Flags related to adoption.
	flags.BoolVar(&opt.Adopt, "adopt", false, "Enable adoption mode for importing a pre-existing cluster into Hive. Will require additional flags for adoption info.")
	flags.StringVar(&opt.AdoptAdminKubeConfig, "adopt-admin-kubeconfig", "", "Path to a cluster admin kubeconfig file for a cluster being adopted. (required if using --adopt)")
	flags.StringVar(&opt.AdoptInfraID, "adopt-infra-id", "", "Infrastructure ID for this cluster's cloud provider. (required if using --adopt)")
	flags.StringVar(&opt.AdoptClusterID, "adopt-cluster-id", "", "Cluster UUID used for telemetry. (required if using --adopt)")
	flags.StringVar(&opt.AdoptAdminUsername, "adopt-admin-username", "", "Username for cluster web console administrator. (optional)")
	flags.StringVar(&opt.AdoptAdminPassword, "adopt-admin-password", "", "Password for cluster web console administrator. (optional)")

	// Azure flags
	flags.StringVar(&opt.AzureBaseDomainResourceGroupName, "azure-base-domain-resource-group-name", "os4-common", "Resource group where the azure DNS zone for the base domain is found")

	return cmd
}

// Complete finishes parsing arguments for the command
func (o *Options) Complete(cmd *cobra.Command, args []string) error {
	o.Name = args[0]
	return nil
}

// Validate ensures that option values make sense
func (o *Options) Validate(cmd *cobra.Command) error {
	if len(o.Output) > 0 && o.Output != "yaml" && o.Output != "json" {
		cmd.Usage()
		log.Info("Invalid value for output. Valid values are: yaml, json.")
		return fmt.Errorf("invalid output")
	}
	if !o.UseClusterImageSet && len(o.ClusterImageSet) > 0 {
		cmd.Usage()
		log.Info("If not using cluster image sets, do not specify the name of one")
		return fmt.Errorf("invalid option")
	}
	if len(o.ServingCert) > 0 && len(o.ServingCertKey) == 0 {
		cmd.Usage()
		log.Info("If specifying a serving certificate, specify a valid serving certificate key")
		return fmt.Errorf("invalid serving cert")
	}
	if !validClouds[o.Cloud] {
		cmd.Usage()
		log.Infof("Unsupported cloud: %s", o.Cloud)
		return fmt.Errorf("Unsupported cloud: %s", o.Cloud)
	}

	if o.Adopt {
		if o.AdoptAdminKubeConfig == "" || o.AdoptInfraID == "" || o.AdoptClusterID == "" {
			return fmt.Errorf("Must specify the following options when using --adopt: --adopt-admin-kube-config, --adopt-infra-id, --adopt-cluster-id")
		}

		if _, err := os.Stat(o.AdoptAdminKubeConfig); os.IsNotExist(err) {
			return fmt.Errorf("--adopt-admin-kubeconfig does not exist: %s", o.AdoptAdminKubeConfig)
		}

		// Admin username and password must both be specified if either are.
		if (o.AdoptAdminUsername != "" || o.AdoptAdminPassword != "") && !(o.AdoptAdminUsername != "" && o.AdoptAdminPassword != "") {
			return fmt.Errorf("--adopt-admin-username and --adopt-admin-password must be used together")
		}
	} else {
		if o.AdoptAdminKubeConfig != "" || o.AdoptInfraID != "" || o.AdoptClusterID != "" || o.AdoptAdminUsername != "" || o.AdoptAdminPassword != "" {
			return fmt.Errorf("Cannot use adoption options without --adopt: --adopt-admin-kube-config, --adopt-infra-id, --adopt-cluster-id, --adopt-admin-username, --adopt-admin-password")
		}
	}
	return nil
}

// Run executes the command
func (o *Options) Run() error {
	if err := apis.AddToScheme(scheme.Scheme); err != nil {
		return err
	}

	objs, err := o.GenerateObjects()
	if err != nil {
		return err
	}
	if len(o.Output) > 0 {
		var printer printers.ResourcePrinter
		if o.Output == "yaml" {
			printer = &printers.YAMLPrinter{}
		} else {
			printer = &printers.JSONPrinter{}
		}
		printObjects(objs, scheme.Scheme, printer)
		return err
	}
	rh, err := o.getResourceHelper()
	if err != nil {
		return err
	}
	if len(o.Namespace) == 0 {
		o.Namespace, err = o.defaultNamespace()
		if err != nil {
			log.Error("Cannot determine default namespace")
			return err
		}
	}
	for _, obj := range objs {
		accessor, err := meta.Accessor(obj)
		if err != nil {
			log.WithError(err).Errorf("Cannot create accessor for object of type %T", obj)
			return err
		}
		accessor.SetNamespace(o.Namespace)
		if _, err := rh.ApplyRuntimeObject(obj, scheme.Scheme); err != nil {
			return err
		}

	}
	return nil
}

func (o *Options) defaultNamespace() (string, error) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
	ns, _, err := kubeconfig.Namespace()
	return ns, err
}

func (o *Options) getResourceHelper() (*resource.Helper, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		log.WithError(err).Error("Cannot get client config")
		return nil, err
	}
	helper := resource.NewHelperFromRESTConfig(cfg, log.WithField("command", "create-cluster"))
	return helper, nil
}

// GenerateObjects generates resources for a new cluster deployment
func (o *Options) GenerateObjects() ([]runtime.Object, error) {

	pullSecret, err := o.getPullSecret()
	if err != nil {
		return nil, err
	}

	sshPrivateKey, err := o.getSSHPrivateKey()
	if err != nil {
		return nil, err
	}

	sshPublicKey, err := o.getSSHPublicKey()
	if err != nil {
		return nil, err
	}

	if len(o.ServingCert) == 0 {
		return nil, nil
	}

	servingCert, err := ioutil.ReadFile(o.ServingCert)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %v", o.ServingCert, err)
	}
	servingCertKey, err := ioutil.ReadFile(o.ServingCertKey)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %v", o.ServingCertKey, err)
	}

	// Load installer manifest files:
	manifestFileData, err := o.getManifestFileBytes()
	if err != nil {
		return nil, err
	}

	generator := &cluster.Generator{
		Name:             o.Name,
		Namespace:        o.Namespace,
		WorkerNodesCount: o.WorkerNodes,
		PullSecret:       pullSecret,
		SSHPrivateKey:    sshPrivateKey,
		SSHPublicKey:     sshPublicKey,
		InstallOnce:      o.InstallOnce,
		BaseDomain:       o.BaseDomain,
		ManageDNS:        o.ManageDNS,
		DeleteAfter:      o.DeleteAfter,
		ServingCert:      string(servingCert),
		ServingCertKey:   string(servingCertKey),
		Labels: map[string]string{
			hiveutilCreatedLabel: "true",
		},
		InstallerManifests: manifestFileData,
	}
	if o.Adopt {
		kubeconfigBytes, err := ioutil.ReadFile(o.AdoptAdminKubeConfig)
		if err != nil {
			return nil, err
		}
		generator.Adopt = o.Adopt
		generator.AdoptInfraID = o.AdoptInfraID
		generator.AdoptClusterID = o.AdoptClusterID
		generator.AdoptAdminKubeconfig = kubeconfigBytes
		generator.AdoptAdminUsername = o.AdoptAdminUsername
		generator.AdoptAdminPassword = o.AdoptAdminPassword
	}

	switch o.Cloud {
	case cloudAWS:
		defaultCredsFilePath := filepath.Join(o.homeDir, ".aws", "credentials")
		accessKeyID, secretAccessKey, err := awsutils.GetAWSCreds(o.CredsFile, defaultCredsFilePath)
		if err != nil {
			return nil, err
		}
		awsProvider := &cluster.AWSCloudProvider{
			AccessKeyID:     accessKeyID,
			SecretAccessKey: secretAccessKey,
		}
		if o.CredsSecret != "" {
			awsProvider.ReuseCredsSecret = &corev1.LocalObjectReference{Name: o.CredsSecret}
		}
		generator.CloudProvider = awsProvider
	case cloudAzure:
		credsFilePath := filepath.Join(os.Getenv("HOME"), ".azure", azureCredFile)
		if l := os.Getenv("AZURE_AUTH_LOCATION"); l != "" {
			credsFilePath = l
		}
		if o.CredsFile != "" {
			credsFilePath = o.CredsFile
		}
		log.Infof("Loading Azure service principal from: %s", credsFilePath)
		spFileContents, err := ioutil.ReadFile(credsFilePath)
		if err != nil {
			return nil, err
		}

		azureProvider := &cluster.AzureCloudProvider{
			ServicePrincipal:            spFileContents,
			BaseDomainResourceGroupName: o.AzureBaseDomainResourceGroupName,
		}
		if o.CredsSecret != "" {
			azureProvider.ReuseCredsSecret = &corev1.LocalObjectReference{Name: o.CredsSecret}
		}
		generator.CloudProvider = azureProvider
	case cloudGCP:
		creds, err := gcputils.GetCreds(o.CredsFile)
		if err != nil {
			return nil, err
		}
		projectID, err := gcpclient.ProjectID(creds)
		if err != nil {
			return nil, err
		}

		gcpProvider := &cluster.GCPCloudProvider{
			ProjectID:      projectID,
			ServiceAccount: creds,
		}
		if o.CredsSecret != "" {
			gcpProvider.ReuseCredsSecret = &corev1.LocalObjectReference{Name: o.CredsSecret}
		}
		generator.CloudProvider = gcpProvider
	}

	imageSet, err := o.configureImages(generator)
	if err != nil {
		return nil, err
	}

	result, err := generator.GenerateAll()
	if err != nil {
		return result, err
	}

	// Add some additional objects we don't yet want to move to the cluster generator library.
	if imageSet != nil {
		result = append(result, imageSet)
	}

	if o.CreateSampleSyncsets {
		result = append(result, o.generateSampleSyncSets()...)
	}

	return result, nil
}

func (o *Options) getPullSecret() (string, error) {
	pullSecret := os.Getenv("PULL_SECRET")
	if len(pullSecret) > 0 {
		return pullSecret, nil
	}
	if len(o.PullSecret) > 0 {
		return o.PullSecret, nil
	}
	if len(o.PullSecretFile) > 0 {
		data, err := ioutil.ReadFile(o.PullSecretFile)
		if err != nil {
			log.Error("Cannot read pull secret file")
			return "", err
		}
		pullSecret = strings.TrimSpace(string(data))
		return pullSecret, nil
	}
	return "", nil
}

func (o *Options) getSSHPublicKey() (string, error) {
	sshPublicKey := os.Getenv("PUBLIC_SSH_KEY")
	if len(sshPublicKey) > 0 {
		return sshPublicKey, nil
	}
	if len(o.SSHPublicKey) > 0 {
		return o.SSHPublicKey, nil
	}
	if len(o.SSHPublicKeyFile) > 0 {
		data, err := ioutil.ReadFile(o.SSHPublicKeyFile)
		if err != nil {
			log.Error("Cannot read SSH public key file")
			return "", err
		}
		sshPublicKey = strings.TrimSpace(string(data))
		return sshPublicKey, nil
	}

	log.Error("Cannot determine SSH key to use")
	return "", nil
}

func (o *Options) getSSHPrivateKey() (string, error) {
	if len(o.SSHPrivateKeyFile) > 0 {
		data, err := ioutil.ReadFile(o.SSHPrivateKeyFile)
		if err != nil {
			log.Error("Cannot read SSH private key file")
			return "", err
		}
		sshPrivateKey := strings.TrimSpace(string(data))
		return sshPrivateKey, nil
	}
	log.Debug("No private SSH key file provided")
	return "", nil
}

func (o *Options) getManifestFileBytes() (map[string][]byte, error) {
	if o.ManifestsDir == "" && !o.SimulateBootstrapFailure {
		return nil, nil
	}
	fileData := map[string][]byte{}
	if o.ManifestsDir != "" {
		files, err := ioutil.ReadDir(o.ManifestsDir)
		if err != nil {
			return nil, errors.Wrap(err, "could not read manifests directory")
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			data, err := ioutil.ReadFile(filepath.Join(o.ManifestsDir, file.Name()))
			if err != nil {
				return nil, errors.Wrapf(err, "could not read manifest file %q", file.Name())
			}
			fileData[file.Name()] = data
		}
	}
	if o.SimulateBootstrapFailure {
		fileData["failure-test.yaml"] = []byte(testFailureManifest)
	}
	return fileData, nil
}

func (o *Options) configureImages(generator *cluster.Generator) (*hivev1.ClusterImageSet, error) {
	if len(o.ClusterImageSet) > 0 {
		generator.ImageSet = o.ClusterImageSet
		return nil, nil
	}
	// TODO: move release image lookup code to the cluster library
	if o.ReleaseImage == "" {
		if o.ReleaseImageSource == "" {
			return nil, fmt.Errorf("Specify either a release image or a release image source")
		}
		var err error
		o.ReleaseImage, err = determineReleaseImageFromSource(o.ReleaseImageSource)
		if err != nil {
			return nil, fmt.Errorf("Cannot determine release image: %v", err)
		}
	}
	if !o.UseClusterImageSet {
		o.ReleaseImage = o.ReleaseImage
		return nil, nil
	}

	imageSet := &hivev1.ClusterImageSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-imageset", o.Name),
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterImageSet",
			APIVersion: hivev1.SchemeGroupVersion.String(),
		},
		Spec: hivev1.ClusterImageSetSpec{
			ReleaseImage: o.ReleaseImage,
		},
	}
	return imageSet, nil
}

func (o *Options) generateSampleSyncSets() []runtime.Object {
	syncsets := []runtime.Object{}
	for i := range [10]int{} {
		syncsets = append(syncsets, sampleSyncSet(fmt.Sprintf("%s-sample-syncset%d", o.Name, i), o.Namespace, o.Name))
		syncsets = append(syncsets, sampleSelectorSyncSet(fmt.Sprintf("sample-selector-syncset%d", i), o.Name))
	}
	return syncsets
}

func sampleSyncSet(name, namespace, cdName string) *hivev1.SyncSet {
	return &hivev1.SyncSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "SyncSet",
			APIVersion: hivev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
		Spec: hivev1.SyncSetSpec{
			ClusterDeploymentRefs: []corev1.LocalObjectReference{
				{
					Name: cdName,
				},
			},
			SyncSetCommonSpec: hivev1.SyncSetCommonSpec{
				ResourceApplyMode: hivev1.SyncResourceApplyMode,
				Resources: []runtime.RawExtension{
					{
						Object: sampleCM(fmt.Sprintf("%s-configmap", name)),
					},
				},
			},
		},
	}
}

func sampleSelectorSyncSet(name, cdName string) *hivev1.SelectorSyncSet {
	return &hivev1.SelectorSyncSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "SelectorSyncSet",
			APIVersion: hivev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: hivev1.SelectorSyncSetSpec{
			ClusterDeploymentSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{hiveutilCreatedLabel: "true"},
			},
			SyncSetCommonSpec: hivev1.SyncSetCommonSpec{
				ResourceApplyMode: hivev1.SyncResourceApplyMode,
				Resources: []runtime.RawExtension{
					{
						Object: sampleCM(fmt.Sprintf("%s-configmap", name)),
					},
				},
			},
		},
	}
}

func sampleCM(name string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Data: map[string]string{
			"foo": "bar",
		},
	}
}

func printObjects(objects []runtime.Object, scheme *runtime.Scheme, printer printers.ResourcePrinter) {
	typeSetterPrinter := printers.NewTypeSetter(scheme).ToPrinter(printer)
	switch len(objects) {
	case 0:
		return
	case 1:
		typeSetterPrinter.PrintObj(objects[0], os.Stdout)
	default:
		list := &metav1.List{
			TypeMeta: metav1.TypeMeta{
				Kind:       "List",
				APIVersion: corev1.SchemeGroupVersion.String(),
			},
			ListMeta: metav1.ListMeta{},
		}
		meta.SetList(list, objects)
		typeSetterPrinter.PrintObj(list, os.Stdout)
	}
}

type releasePayload struct {
	PullSpec string `json:"pullSpec"`
}

func determineReleaseImageFromSource(sourceURL string) (string, error) {
	resp, err := http.Get(sourceURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	payload := &releasePayload{}
	err = json.Unmarshal(data, payload)
	if err != nil {
		return "", err
	}
	return payload.PullSpec, nil
}
