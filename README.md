# eks-auth-sync

A tool to help manage EKS cluster authentication configuration.

## What's the problem this tool solves?

Authentication in EKS is configured using [a single ConfigMap in the Kubernetes cluster][aws-auth] (`aws-auth`) that maps AWS IAM users and roles to users in Kubernetes.
If you only have a fixed number of users and roles assigned to a cluster, it's easy enough to just create the ConfigMap once and forget about it.
However, if the number of users and roles varies frequently (i.e. people join and leave the cluster), managing the ConfigMap can become a chore.

To help automate the ConfigMap updates, you can use eks-auth-sync to automatically pull changes from various sources and update the ConfigMap.

[aws-auth]: https://docs.aws.amazon.com/eks/latest/userguide/add-user-role.html

## How does it work?

Here's roughly what the tool does when you run it:

1. Read a given configuration file for a list of data sources (called "sccanners").
2. Read the data sources for all the available auth mappings, and join the results.
3. Update the auth mappings in your EKS cluster.

## Usage

The `eks-auth-sync` tool accepts the following parameters:

* `-config string`:
  Path to the YAML configuration file.
  Set to `-` to read the config from STDIN.
  See section below for more details on the configuration structure. 
* `-commit`:
  If enabled, the scanned results will be committed to EKS.
  Otherwise, they'll be printed to STDOUT instead.
* `-version`:
  Print the version information.

## Configuration

The `eks-auth-sync` tool is configured using a YAML file.
Here's how it's structured:

```yaml
# Configurations related to accessing
# the target Kubernetes cluster (EKS)
kubernetes:
  # If set to true, assume that the tool is run in Kubernetes
  # with configuration from Kubernetes (default: false)
  inKubeCluster: <boolean>
  # Path to a Kubernetes configuration file. (default: "")
  # Only used when `inKubeCluster` is `false`.
  kubeConfigPath: /path/to/kubeconfig

# Configurations for AWS client
aws:
  # Role to assume before interacting with AWS (default: "")
  # If left empty, no role assumption is made.
  roleARN: <string>
  # Endpoint for AWS API (default: "")
  # If left empty, the default public AWS endpoint is used.
  endpoint: <string>
  # AWS region e.g. eu-north-1 (default: "")
  # If left empty, the default region is used.
  region: <string>
  # If set to true, SSL/TLS is disabled when communicating
  # with the AWS APIs. (default: false)
  disableSSL: <boolean>
  # Maximum number of retries the AWS client makes per request.
  maxRetries: <integer>

# List of data sources to scan for authentication mappings
# See the section below for a list of scanners you can use.
# The results from all the scanners are concatenated.
scanners:
- <scanner>
```

Each scanner is structured as follows:

```yaml
# Optional name to display in error messages.
name: <string>

# Configurations for AWS client specific to this scanner.
# It has the same fields as the AWS configuration section above.
aws: <aws section from above>

# Type of the scanner to use.
# This setting defines which below sections will be used.
# Available options: static, file, iam, ssm
type: <string>

# Use the mappings specified in this section.
# Only used in combination with type=static
static:
  # IAM user to Kubernetes user mappings
  users:
  - # ARN of the IAM user to map
    userarn: <string>
    # Kubernetes username for the IAM user
    username: <string>
    # Kubernetes groups for the IAM user
    groups:
    - <string>
  # IAM role to Kubernetes user mappings
  roles:
  - # ARN of the IAM role to map
    rolearn: <string>
    # Kubernetes username for the IAM role
    username: <string>
    # Kubernetes groups for the IAM role
    groups:
    - <string>

# Read mappings from a file.
# Only used in combination with type=file
file:
  # Path to a YAML file containing the mappings.
  # The file should be structured the same as the mappings
  # in the above "static" section.
  path: <string>

# Read mappings from IAM user/role tags.
# Only used in combination with type=iam
iam:
  # Name of the EKS cluster the mappings are intended for.
  clusterName: <string>
  # The AWS account the EKS cluster is deployed in.
  clusterAccountID: <string>
  # IAM path to scan for users and roles. (default: "")
  # If not set, scan all the users and roles.
  pathPrefix: <string>
  # If set to true, skip scanning IAM users.
  # This can be used for limiting scans to IAM roles.
  disableUserScan: <bool>
  # If set to true, skip scanning IAM roles.
  # This can be used for limiting scans to IAM users.
  disableRoleScan: <bool>

# Read mappings from an SSM parameter.
# Only used in combination with type=ssm
ssm:
  # Path of the SSM Parameter to read mappings from.
  # The parameter contents should be in YAML format
  # similar to the "static" section above.
  path: <string>
```

## Development

See [Development page](DEVELOPMENT.md).

## License

Apache 2.0. See [LICENSE](LICENSE) for more information.

---

Made with ❤️ by [Polar Squad](https://polarsquad.com/)