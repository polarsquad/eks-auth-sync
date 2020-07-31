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

## Documentation

* [Configuration](docs/configuration.md)
* Reading mappings
  * [Directly from configuration](docs/static.md)
  * [Files on disk](docs/file.md)
  * [IAM](docs/iam.md)
  * [SSM](docs/ssm.md)
* [Development](docs/development.md)

## License

Apache 2.0. See [LICENSE](LICENSE) for more information.

---

Made with ❤️ by [Polar Squad](https://polarsquad.com/)