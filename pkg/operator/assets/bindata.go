package assets

import (
	"fmt"
	"strings"
)

var _config_crds_hive_v1alpha1_clusterdeployment_yaml = []byte(`apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  creationTimestamp: null
  labels:
    controller-tools.k8s.io: "1.0"
  name: clusterdeployments.hive.openshift.io
spec:
  group: hive.openshift.io
  names:
    kind: ClusterDeployment
    plural: clusterdeployments
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          properties:
            baseDomain:
              description: BaseDomain is the base domain to which the cluster should
                belong.
              type: string
            clusterName:
              description: ClusterName is the friendly name of the cluster. It is
                used for subdomains, some resource tagging, and other instances where
                a friendly name for the cluster is useful.
              type: string
            compute:
              description: Compute is the list of MachinePools containing compute
                nodes that need to be installed.
              items:
                properties:
                  name:
                    description: Name is the name of the machine pool.
                    type: string
                  platform:
                    description: Platform is configuration for machine pool specific
                      to the platfrom.
                    properties:
                      aws:
                        description: AWS is the configuration used when installing
                          on AWS.
                        properties:
                          iamRoleName:
                            description: IAMRoleName defines the IAM role associated
                              with the ec2 instance.
                            type: string
                          rootVolume:
                            description: EC2RootVolume defines the storage for ec2
                              instance.
                            properties:
                              iops:
                                description: IOPS defines the iops for the instance.
                                format: int64
                                type: integer
                              size:
                                description: Size defines the size of the instance.
                                format: int64
                                type: integer
                              type:
                                description: Type defines the type of the instance.
                                type: string
                            required:
                            - iops
                            - size
                            - type
                            type: object
                          type:
                            description: InstanceType defines the ec2 instance type.
                              eg. m4-large
                            type: string
                          zones:
                            description: Zones is list of availability zones that
                              can be used.
                            items:
                              type: string
                            type: array
                        required:
                        - type
                        - iamRoleName
                        - rootVolume
                        type: object
                      libvirt:
                        description: Libvirt is the configuration used when installing
                          on libvirt.
                        properties:
                          image:
                            description: Image is the URL to the OS image. E.g. "http://aos-ostree.rhev-ci-vms.eng.rdu2.redhat.com/rhcos/images/cloud/latest/rhcos-qemu.qcow2.gz"
                            type: string
                          imagePool:
                            description: ImagePool is the name of the libvirt storage
                              pool to which the storage volume containing the OS image
                              belongs.
                            type: string
                          imageVolume:
                            description: ImageVolume is the name of the libvirt storage
                              volume containing the OS image.
                            type: string
                        required:
                        - image
                        type: object
                      openstack:
                        description: OpenStack is the configuration used when installing
                          on OpenStack.
                        properties:
                          rootVolume:
                            description: OpenStackRootVolume defines the storage for
                              Nova instance.
                            properties:
                              iops:
                                description: IOPS defines the iops for the instance.
                                format: int64
                                type: integer
                              size:
                                description: Size defines the size of the instance.
                                format: int64
                                type: integer
                              type:
                                description: Type defines the type of the instance.
                                type: string
                            required:
                            - iops
                            - size
                            - type
                            type: object
                          type:
                            description: FlavorName defines the OpenStack Nova flavor.
                              eg. m1.large
                            type: string
                        required:
                        - type
                        - rootVolume
                        type: object
                    type: object
                  replicas:
                    description: Replicas is the count of machines for this machine
                      pool. Default is 1.
                    format: int64
                    type: integer
                required:
                - name
                - replicas
                - platform
                type: object
              type: array
            controlPlane:
              description: ControlPlane is the MachinePool containing control plane
                nodes that need to be installed.
              properties:
                name:
                  description: Name is the name of the machine pool.
                  type: string
                platform:
                  description: Platform is configuration for machine pool specific
                    to the platfrom.
                  properties:
                    aws:
                      description: AWS is the configuration used when installing on
                        AWS.
                      properties:
                        iamRoleName:
                          description: IAMRoleName defines the IAM role associated
                            with the ec2 instance.
                          type: string
                        rootVolume:
                          description: EC2RootVolume defines the storage for ec2 instance.
                          properties:
                            iops:
                              description: IOPS defines the iops for the instance.
                              format: int64
                              type: integer
                            size:
                              description: Size defines the size of the instance.
                              format: int64
                              type: integer
                            type:
                              description: Type defines the type of the instance.
                              type: string
                          required:
                          - iops
                          - size
                          - type
                          type: object
                        type:
                          description: InstanceType defines the ec2 instance type.
                            eg. m4-large
                          type: string
                        zones:
                          description: Zones is list of availability zones that can
                            be used.
                          items:
                            type: string
                          type: array
                      required:
                      - type
                      - iamRoleName
                      - rootVolume
                      type: object
                    libvirt:
                      description: Libvirt is the configuration used when installing
                        on libvirt.
                      properties:
                        image:
                          description: Image is the URL to the OS image. E.g. "http://aos-ostree.rhev-ci-vms.eng.rdu2.redhat.com/rhcos/images/cloud/latest/rhcos-qemu.qcow2.gz"
                          type: string
                        imagePool:
                          description: ImagePool is the name of the libvirt storage
                            pool to which the storage volume containing the OS image
                            belongs.
                          type: string
                        imageVolume:
                          description: ImageVolume is the name of the libvirt storage
                            volume containing the OS image.
                          type: string
                      required:
                      - image
                      type: object
                    openstack:
                      description: OpenStack is the configuration used when installing
                        on OpenStack.
                      properties:
                        rootVolume:
                          description: OpenStackRootVolume defines the storage for
                            Nova instance.
                          properties:
                            iops:
                              description: IOPS defines the iops for the instance.
                              format: int64
                              type: integer
                            size:
                              description: Size defines the size of the instance.
                              format: int64
                              type: integer
                            type:
                              description: Type defines the type of the instance.
                              type: string
                          required:
                          - iops
                          - size
                          - type
                          type: object
                        type:
                          description: FlavorName defines the OpenStack Nova flavor.
                            eg. m1.large
                          type: string
                      required:
                      - type
                      - rootVolume
                      type: object
                  type: object
                replicas:
                  description: Replicas is the count of machines for this machine
                    pool. Default is 1.
                  format: int64
                  type: integer
              required:
              - name
              - replicas
              - platform
              type: object
            images:
              description: Images allows overriding the default images used to provision
                and manage the cluster.
              properties:
                hiveImage:
                  description: HiveImage is the image used in the sidecar container
                    to manage execution of openshift-install.
                  type: string
                hiveImagePullPolicy:
                  description: HiveImagePullPolicy is the pull policy for the installer
                    image.
                  type: string
                installerImage:
                  description: InstallerImage is the image containing the openshift-install
                    binary that will be used to install.
                  type: string
                installerImagePullPolicy:
                  description: InstallerImagePullPolicy is the pull policy for the
                    installer image.
                  type: string
                releaseImage:
                  description: ReleaseImage is the image containing metadata for all
                    components that run in the cluster, and is the primary and best
                    way to specify what specific version of OpenShift you wish to
                    install.
                  type: string
              type: object
            networking:
              description: Networking defines the pod network provider in the cluster.
              properties:
                clusterNetworks:
                  description: ClusterNetworks is the IP address space from which
                    to assign pod IPs.
                  items:
                    properties:
                      cidr:
                        type: string
                      hostSubnetLength:
                        format: int32
                        type: integer
                    required:
                    - cidr
                    - hostSubnetLength
                    type: object
                  type: array
                machineCIDR:
                  description: MachineCIDR is the IP address space from which to assign
                    machine IPs.
                  type: string
                serviceCIDR:
                  description: ServiceCIDR is the IP address space from which to assign
                    service IPs.
                  type: string
                type:
                  description: Type is the network type to install
                  type: string
              required:
              - machineCIDR
              - type
              - serviceCIDR
              type: object
            platform:
              description: Platform is the configuration for the specific platform
                upon which to perform the installation.
              properties:
                aws:
                  description: AWS is the configuration used when installing on AWS.
                  properties:
                    defaultMachinePlatform:
                      description: DefaultMachinePlatform is the default configuration
                        used when installing on AWS for machine pools which do not
                        define their own platform configuration.
                      properties:
                        iamRoleName:
                          description: IAMRoleName defines the IAM role associated
                            with the ec2 instance.
                          type: string
                        rootVolume:
                          description: EC2RootVolume defines the storage for ec2 instance.
                          properties:
                            iops:
                              description: IOPS defines the iops for the instance.
                              format: int64
                              type: integer
                            size:
                              description: Size defines the size of the instance.
                              format: int64
                              type: integer
                            type:
                              description: Type defines the type of the instance.
                              type: string
                          required:
                          - iops
                          - size
                          - type
                          type: object
                        type:
                          description: InstanceType defines the ec2 instance type.
                            eg. m4-large
                          type: string
                        zones:
                          description: Zones is list of availability zones that can
                            be used.
                          items:
                            type: string
                          type: array
                      required:
                      - type
                      - iamRoleName
                      - rootVolume
                      type: object
                    region:
                      description: Region specifies the AWS region where the cluster
                        will be created.
                      type: string
                    userTags:
                      description: UserTags specifies additional tags for AWS resources
                        created for the cluster.
                      type: object
                  required:
                  - region
                  type: object
                libvirt:
                  description: Libvirt is the configuration used when installing on
                    libvirt.
                  properties:
                    URI:
                      description: URI is the identifier for the libvirtd connection.  It
                        must be reachable from both the host (where the installer
                        is run) and the cluster (where the cluster-API controller
                        pod will be running).
                      type: string
                    defaultMachinePlatform:
                      description: DefaultMachinePlatform is the default configuration
                        used when installing on AWS for machine pools which do not
                        define their own platform configuration.
                      properties:
                        image:
                          description: Image is the URL to the OS image. E.g. "http://aos-ostree.rhev-ci-vms.eng.rdu2.redhat.com/rhcos/images/cloud/latest/rhcos-qemu.qcow2.gz"
                          type: string
                        imagePool:
                          description: ImagePool is the name of the libvirt storage
                            pool to which the storage volume containing the OS image
                            belongs.
                          type: string
                        imageVolume:
                          description: ImageVolume is the name of the libvirt storage
                            volume containing the OS image.
                          type: string
                      required:
                      - image
                      type: object
                    masterIPs:
                      description: MasterIPs
                      items:
                        format: byte
                        type: string
                      type: array
                    network:
                      description: Network
                      properties:
                        if:
                          description: IfName is the name of the network interface.
                          type: string
                        ipRange:
                          description: IPRange is the range of IPs to use.
                          type: string
                        name:
                          description: Name is the name of the nework.
                          type: string
                      required:
                      - name
                      - if
                      - ipRange
                      type: object
                  required:
                  - URI
                  - network
                  - masterIPs
                  type: object
              type: object
            platformSecrets:
              description: PlatformSecrets contains credentials and secrets for the
                cluster infrastructure.
              properties:
                aws:
                  properties:
                    credentials:
                      description: Credentials refers to a secret that contains the
                        AWS account access credentials.
                      type: object
                  required:
                  - credentials
                  type: object
              type: object
            preserveOnDelete:
              description: PreserveOnDelete allows the user to disconnect a cluster
                from Hive without deprovisioning it
              type: boolean
            pullSecret:
              description: PullSecret is the reference to the secret to use when pulling
                images.
              type: object
            sshKey:
              description: SSHKey is the reference to the secret that contains a public
                key to use for access to compute instances.
              type: object
          required:
          - clusterName
          - baseDomain
          - networking
          - controlPlane
          - compute
          - platform
          - pullSecret
          - platformSecrets
          type: object
        status:
          properties:
            adminKubeconfigSecret:
              description: AdminKubeconfigSecret references the secret containing
                the admin kubeconfig for this cluster.
              type: object
            adminPasswordSecret:
              description: AdminPasswordSecret references the secret containing the
                admin username/password which can be used to login to this cluster.
              type: object
            apiURL:
              description: APIURL is the URL where the cluster's API can be accessed.
              type: string
            clusterID:
              description: ClusterID is a unique identifier for this cluster generated
                during installation.
              type: string
            clusterVersionStatus:
              description: ClusterVersionStatus will hold a copy of the remote cluster's
                ClusterVersion.Status
              properties:
                availableUpdates:
                  description: availableUpdates contains the list of updates that
                    are appropriate for this cluster. This list may be empty if no
                    updates are recommended, if the update service is unavailable,
                    or if an invalid channel has been specified.
                  items:
                    properties:
                      payload:
                        description: payload is a container image location that contains
                          the update. When this field is part of spec, payload is
                          optional if version is specified and the availableUpdates
                          field contains a matching version.
                        type: string
                      version:
                        description: version is a semantic versioning identifying
                          the update version. When this field is part of spec, version
                          is optional if payload is specified.
                        type: string
                    required:
                    - version
                    - payload
                    type: object
                  type: array
                conditions:
                  description: conditions provides information about the cluster version.
                    The condition "Available" is set to true if the desiredUpdate
                    has been reached. The condition "Progressing" is set to true if
                    an update is being applied. The condition "Failing" is set to
                    true if an update is currently blocked by a temporary or permanent
                    error. Conditions are only valid for the current desiredUpdate
                    when metadata.generation is equal to status.generation.
                  items:
                    properties:
                      lastTransitionTime:
                        description: lastTransitionTime is the time of the last update
                          to the current status object.
                        format: date-time
                        type: string
                      message:
                        description: message provides additional information about
                          the current condition. This is only to be consumed by humans.
                        type: string
                      reason:
                        description: reason is the reason for the condition's last
                          transition.  Reasons are CamelCase
                        type: string
                      status:
                        description: status of the condition, one of True, False,
                          Unknown.
                        type: string
                      type:
                        description: type specifies the state of the operator's reconciliation
                          functionality.
                        type: string
                    required:
                    - type
                    - status
                    - lastTransitionTime
                    type: object
                  type: array
                desired:
                  description: desired is the version that the cluster is reconciling
                    towards. If the cluster is not yet fully initialized desired will
                    be set with the information available, which may be a payload
                    or a tag.
                  properties:
                    payload:
                      description: payload is a container image location that contains
                        the update. When this field is part of spec, payload is optional
                        if version is specified and the availableUpdates field contains
                        a matching version.
                      type: string
                    version:
                      description: version is a semantic versioning identifying the
                        update version. When this field is part of spec, version is
                        optional if payload is specified.
                      type: string
                  required:
                  - version
                  - payload
                  type: object
                generation:
                  description: generation reports which version of the spec is being
                    processed. If this value is not equal to metadata.generation,
                    then the current and conditions fields have not yet been updated
                    to reflect the latest request.
                  format: int64
                  type: integer
                history:
                  description: history contains a list of the most recent versions
                    applied to the cluster. This value may be empty during cluster
                    startup, and then will be updated when a new update is being applied.
                    The newest update is first in the list and it is ordered by recency.
                    Updates in the history have state Completed if the rollout completed
                    - if an update was failing or halfway applied the state will be
                    Partial. Only a limited amount of update history is preserved.
                  items:
                    properties:
                      completionTime:
                        description: completionTime, if set, is when the update was
                          fully applied. The update that is currently being applied
                          will have a null completion time. Completion time will always
                          be set for entries that are not the current update (usually
                          to the started time of the next update).
                        format: date-time
                        type: string
                      payload:
                        description: payload is a container image location that contains
                          the update. This value is always populated.
                        type: string
                      startedTime:
                        description: startedTime is the time at which the update was
                          started.
                        format: date-time
                        type: string
                      state:
                        description: state reflects whether the update was fully applied.
                          The Partial state indicates the update is not fully applied,
                          while the Completed state indicates the update was successfully
                          rolled out at least once (all parts of the update successfully
                          applied).
                        type: string
                      version:
                        description: version is a semantic versioning identifying
                          the update version. If the requested payload does not define
                          a version, or if a failure occurs retrieving the payload,
                          this value may be empty.
                        type: string
                    required:
                    - state
                    - startedTime
                    - completionTime
                    - version
                    - payload
                    type: object
                  type: array
                versionHash:
                  description: versionHash is a fingerprint of the content that the
                    cluster will be updated with. It is used by the operator to avoid
                    unnecessary work and is for internal use only.
                  type: string
              required:
              - desired
              - history
              - generation
              - versionHash
              - conditions
              - availableUpdates
              type: object
            federated:
              description: Federated is true if the cluster deployment has been federated
                with the host cluster.
              type: boolean
            installed:
              description: Installed is true if the installer job has successfully
                completed for this cluster.
              type: boolean
            webConsoleURL:
              description: WebConsoleURL is the URL for the cluster's web console
                UI.
              type: string
          required:
          - installed
          type: object
  version: v1alpha1
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func config_crds_hive_v1alpha1_clusterdeployment_yaml() ([]byte, error) {
	return _config_crds_hive_v1alpha1_clusterdeployment_yaml, nil
}

var _config_crds_hive_v1alpha1_dnszone_yaml = []byte(`apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  creationTimestamp: null
  labels:
    controller-tools.k8s.io: "1.0"
  name: dnszones.hive.openshift.io
spec:
  group: hive.openshift.io
  names:
    kind: DNSZone
    plural: dnszones
  scope: Namespaced
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          properties:
            aws:
              description: AWS specifies AWS-specific cloud configuration
              properties:
                accountSecret:
                  description: AccountSecret contains a reference to a secret that
                    contains AWS credentials for CRUD operations
                  type: object
                region:
                  description: Region specifies the region-specific API endpoint to
                    use
                  type: string
              required:
              - accountSecret
              - region
              type: object
            zone:
              description: Zone is the DNS zoneto host
              type: string
          required:
          - zone
          type: object
        status:
          properties:
            lastSyncGeneration:
              description: LastSyncGeneration is the generation of the zone resource
                that was last sync'd. This is used to know if the Object has changed
                and we should sync immediately.
              format: int64
              type: integer
            lastSyncTimestamp:
              description: LastSyncTimestamp is the time that the zone was last sync'd.
              format: date-time
              type: string
          required:
          - lastSyncGeneration
          type: object
  version: v1alpha1
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func config_crds_hive_v1alpha1_dnszone_yaml() ([]byte, error) {
	return _config_crds_hive_v1alpha1_dnszone_yaml, nil
}

var _config_crds_hive_v1alpha1_hiveadmissionconfig_yaml = []byte(`apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  creationTimestamp: null
  labels:
    controller-tools.k8s.io: "1.0"
  name: hiveadmissionconfigs.hive.openshift.io
spec:
  group: hive.openshift.io
  names:
    kind: HiveAdmissionConfig
    plural: hiveadmissionconfigs
  scope: Namespaced
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          properties:
            image:
              description: Image controls the image used for HiveAdmission pods. Default
                is to use latest master images.
              type: string
          type: object
        status:
          type: object
  version: v1alpha1
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func config_crds_hive_v1alpha1_hiveadmissionconfig_yaml() ([]byte, error) {
	return _config_crds_hive_v1alpha1_hiveadmissionconfig_yaml, nil
}

var _config_crds_hive_v1alpha1_hiveconfig_yaml = []byte(`apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  creationTimestamp: null
  labels:
    controller-tools.k8s.io: "1.0"
  name: hiveconfigs.hive.openshift.io
spec:
  group: hive.openshift.io
  names:
    kind: HiveConfig
    plural: hiveconfigs
  scope: Namespaced
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          properties:
            image:
              description: Image controls the image used for the Hive controllers,
                as well as deprovision pods. The default if left empty is to use latest
                master.
              type: string
          type: object
        status:
          type: object
  version: v1alpha1
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func config_crds_hive_v1alpha1_hiveconfig_yaml() ([]byte, error) {
	return _config_crds_hive_v1alpha1_hiveconfig_yaml, nil
}

var _config_hiveadmission_apiservice_yaml = []byte(`---
# register as aggregated apiserver; this has a number of benefits:
#
# - allows other kubernetes components to talk to the the admission webhook using the ` + "`" + `kubernetes.default.svc` + "`" + ` service
# - allows other kubernetes components to use their in-cluster credentials to communicate with the webhook
# - allows you to test the webhook using kubectl
# - allows you to govern access to the webhook using RBAC
# - prevents other extension API servers from leaking their service account tokens to the webhook
#
# for more information, see: https://kubernetes.io/blog/2018/01/extensible-admission-is-beta
apiVersion: apiregistration.k8s.io/v1
kind: APIService
metadata:
  name: v1alpha1.admission.hive.openshift.io
  annotations:
    service.alpha.openshift.io/inject-cabundle: "true"
spec:
  group: admission.hive.openshift.io
  groupPriorityMinimum: 1000
  versionPriority: 15
  service:
    name: hiveadmission
    namespace: openshift-hive
  version: v1alpha1
`)

func config_hiveadmission_apiservice_yaml() ([]byte, error) {
	return _config_hiveadmission_apiservice_yaml, nil
}

var _config_hiveadmission_clusterdeployment_webhook_yaml = []byte(`---
apiVersion: admissionregistration.k8s.io/v1beta1
kind: ValidatingWebhookConfiguration
metadata:
  name: clusterdeployments.admission.hive.openshift.io
webhooks:
- name: clusterdeployments.admission.hive.openshift.io
  clientConfig:
    service:
      # reach the webhook via the registered aggregated API
      namespace: default
      name: kubernetes
      path: /apis/admission.hive.openshift.io/v1alpha1/clusterdeployments
  rules:
  - operations:
    - CREATE
    - UPDATE
    apiGroups:
    - hive.openshift.io
    apiVersions:
    - v1alpha1
    resources:
    - clusterdeployments
  failurePolicy: Fail
`)

func config_hiveadmission_clusterdeployment_webhook_yaml() ([]byte, error) {
	return _config_hiveadmission_clusterdeployment_webhook_yaml, nil
}

var _config_hiveadmission_daemonset_yaml = []byte(`---
# to create the namespace-reservation-server
apiVersion: apps/v1
kind: DaemonSet
metadata:
  namespace: openshift-hive
  name: hiveadmission
  labels:
    app: hiveadmission
    hiveadmission: "true"
spec:
  selector:
    matchLabels:
      app: hiveadmission
      hiveadmission: "true"
  updateStrategy:
    rollingUpdate:
      maxUnavailable: 1
    type: RollingUpdate
  template:
    metadata:
      name: hiveadmission
      labels:
        app: hiveadmission
        hiveadmission: "true"
    spec:
      serviceAccountName: hiveadmission
      containers:
      - name: hiveadmission
        image: registry.svc.ci.openshift.org/openshift/hive-v4.0:hive
        imagePullPolicy: Always
        command:
        - "/opt/services/hiveadmission"
        - "--secure-port=9443"
        - "--audit-log-path=-"
        - "--tls-cert-file=/var/serving-cert/tls.crt"
        - "--tls-private-key-file=/var/serving-cert/tls.key"
        - "--v=8"
        ports:
        - containerPort: 9443
        volumeMounts:
        - mountPath: /var/serving-cert
          name: serving-cert
        readinessProbe:
          httpGet:
            path: /healthz
            port: 9443
            scheme: HTTPS
      nodeSelector:
        node-role.kubernetes.io/master: ""
      tolerations:
      - effect: NoSchedule
        key: node-role.kubernetes.io/master
        operator: Exists
      volumes:
      - name: serving-cert
        secret:
          defaultMode: 420
          secretName: hiveadmission-serving-cert
`)

func config_hiveadmission_daemonset_yaml() ([]byte, error) {
	return _config_hiveadmission_daemonset_yaml, nil
}

var _config_hiveadmission_dnszones_webhook_yaml = []byte(`---
# register to intercept DNSZone object creates and updates
apiVersion: admissionregistration.k8s.io/v1beta1
kind: ValidatingWebhookConfiguration
metadata:
  name: dnszones.admission.hive.openshift.io
webhooks:
- name: dnszones.admission.hive.openshift.io
  clientConfig:
    service:
      # reach the webhook via the registered aggregated API
      namespace: default
      name: kubernetes
      path: /apis/admission.hive.openshift.io/v1alpha1/dnszones
  rules:
  - operations:
    - CREATE
    - UPDATE
    apiGroups:
    - hive.openshift.io
    apiVersions:
    - v1alpha1
    resources:
    - dnszones
  failurePolicy: Fail
`)

func config_hiveadmission_dnszones_webhook_yaml() ([]byte, error) {
	return _config_hiveadmission_dnszones_webhook_yaml, nil
}

var _config_hiveadmission_service_account_yaml = []byte(`---
# to be able to assign powers to the hiveadmission process
apiVersion: v1
kind: ServiceAccount
metadata:
  namespace: openshift-hive
  name: hiveadmission
`)

func config_hiveadmission_service_account_yaml() ([]byte, error) {
	return _config_hiveadmission_service_account_yaml, nil
}

var _config_hiveadmission_service_yaml = []byte(`---
apiVersion: v1
kind: Service
metadata:
  namespace: openshift-hive
  name: hiveadmission
  annotations:
    service.alpha.openshift.io/serving-cert-secret-name: hiveadmission-serving-cert
spec:
  selector:
    app: hiveadmission
  ports:
  - port: 443
    targetPort: 9443
`)

func config_hiveadmission_service_yaml() ([]byte, error) {
	return _config_hiveadmission_service_yaml, nil
}

var _config_manager_deployment_yaml = []byte(`---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hive-controllers
  namespace: openshift-hive
  labels:
    control-plane: controller-manager
    controller-tools.k8s.io: "1.0"
spec:
  selector:
    matchLabels:
      control-plane: controller-manager
      controller-tools.k8s.io: "1.0"
  replicas: 1
  revisionHistoryLimit: 4
  template:
    metadata:
      labels:
        control-plane: controller-manager
        controller-tools.k8s.io: "1.0"
    spec:
      containers:
      # By default we will use the latest CI images published from hive master:
      - image: registry.svc.ci.openshift.org/openshift/hive-v4.0:hive
        imagePullPolicy: Always
        name: manager
        resources:
          limits:
            cpu: 100m
            memory: 256Mi
          requests:
            cpu: 100m
            memory: 75Mi
        command:
          - /opt/services/manager
          - --log-level
          - debug
      terminationGracePeriodSeconds: 10
`)

func config_manager_deployment_yaml() ([]byte, error) {
	return _config_manager_deployment_yaml, nil
}

var _config_manager_service_yaml = []byte(`---
apiVersion: v1
kind: Service
metadata:
  name: hive-controllers-service
  namespace: openshift-hive
  labels:
    control-plane: controller-manager
    controller-tools.k8s.io: "1.0"
spec:
  selector:
    control-plane: controller-manager
    controller-tools.k8s.io: "1.0"
  ports:
  - port: 443
`)

func config_manager_service_yaml() ([]byte, error) {
	return _config_manager_service_yaml, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() ([]byte, error){
	"config/crds/hive_v1alpha1_clusterdeployment.yaml":    config_crds_hive_v1alpha1_clusterdeployment_yaml,
	"config/crds/hive_v1alpha1_dnszone.yaml":              config_crds_hive_v1alpha1_dnszone_yaml,
	"config/crds/hive_v1alpha1_hiveadmissionconfig.yaml":  config_crds_hive_v1alpha1_hiveadmissionconfig_yaml,
	"config/crds/hive_v1alpha1_hiveconfig.yaml":           config_crds_hive_v1alpha1_hiveconfig_yaml,
	"config/hiveadmission/apiservice.yaml":                config_hiveadmission_apiservice_yaml,
	"config/hiveadmission/clusterdeployment-webhook.yaml": config_hiveadmission_clusterdeployment_webhook_yaml,
	"config/hiveadmission/daemonset.yaml":                 config_hiveadmission_daemonset_yaml,
	"config/hiveadmission/dnszones-webhook.yaml":          config_hiveadmission_dnszones_webhook_yaml,
	"config/hiveadmission/service-account.yaml":           config_hiveadmission_service_account_yaml,
	"config/hiveadmission/service.yaml":                   config_hiveadmission_service_yaml,
	"config/manager/deployment.yaml":                      config_manager_deployment_yaml,
	"config/manager/service.yaml":                         config_manager_service_yaml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func     func() ([]byte, error)
	Children map[string]*_bintree_t
}

var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"config": {nil, map[string]*_bintree_t{
		"crds": {nil, map[string]*_bintree_t{
			"hive_v1alpha1_clusterdeployment.yaml":   {config_crds_hive_v1alpha1_clusterdeployment_yaml, map[string]*_bintree_t{}},
			"hive_v1alpha1_dnszone.yaml":             {config_crds_hive_v1alpha1_dnszone_yaml, map[string]*_bintree_t{}},
			"hive_v1alpha1_hiveadmissionconfig.yaml": {config_crds_hive_v1alpha1_hiveadmissionconfig_yaml, map[string]*_bintree_t{}},
			"hive_v1alpha1_hiveconfig.yaml":          {config_crds_hive_v1alpha1_hiveconfig_yaml, map[string]*_bintree_t{}},
		}},
		"hiveadmission": {nil, map[string]*_bintree_t{
			"apiservice.yaml":                {config_hiveadmission_apiservice_yaml, map[string]*_bintree_t{}},
			"clusterdeployment-webhook.yaml": {config_hiveadmission_clusterdeployment_webhook_yaml, map[string]*_bintree_t{}},
			"daemonset.yaml":                 {config_hiveadmission_daemonset_yaml, map[string]*_bintree_t{}},
			"dnszones-webhook.yaml":          {config_hiveadmission_dnszones_webhook_yaml, map[string]*_bintree_t{}},
			"service-account.yaml":           {config_hiveadmission_service_account_yaml, map[string]*_bintree_t{}},
			"service.yaml":                   {config_hiveadmission_service_yaml, map[string]*_bintree_t{}},
		}},
		"manager": {nil, map[string]*_bintree_t{
			"deployment.yaml": {config_manager_deployment_yaml, map[string]*_bintree_t{}},
			"service.yaml":    {config_manager_service_yaml, map[string]*_bintree_t{}},
		}},
	}},
}}