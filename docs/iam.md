# Reading mappings from a IAM user/role tags

The `eks-auth-sync` tool can read mappings from IAM user/role tags.
This can be useful when you want to store the mappings as close to the IAM users and roles as possible.

You can scan mappings from IAM user/role tags by adding a scanner entry to your configuration file with the `type` set to `iam`.
The `iam` section must specify the name of the EKS cluster and the ID of the AWS account where the cluster is located.

## Tags

Here's all the tags that the tool scans for:

* `eks/{account id}/{cluster name}/username`:
  Username to map the user/role to in Kubernetes
* `eks/{account id}/{cluster name}/groups`:
  A comma-separated list of Kubernetes groups the user/role is assigned to
* `eks/{account id}/{cluster name}/type`:
  When set to `node`, the role is interpreted as an EKS worker node role.
  When set to `user`, the role is interpreted as a normal Kubernetes user.
  Only applicable for roles.

Where:

* `{cluster name}` is the name of the EKS cluster the mappings should be saved to.
* `{account id}` is the ID of the AWS account where the EKS cluster is located.

## Example

To use the `iam` scanner, create a configuration YAML file like this:

```yaml
kubernetes:
  kubeConfigPath: ~/.kube/config
scanners:
- type: iam
  iam:
    clusterName: mycluster
    clusterAccountID: 123456789012
    pathPrefix: /eks/
```

This will make `eks-auth-sync` scan IAM user/role tags for entries that contain the given cluster name and account ID.
Make sure to replace those details with the details of your EKS cluster.
Using the path prefix, we can limit the scan to users/roles that start with path `/eks/`.

Let's create some users and roles to test the configuration with.
First, we'll create a normal user `johannes` that's assigned to a Kubernetes group called `admin`.

```
$ aws iam create-user \
    --user-name johannes \
    --path /eks/ \
    --tags \
        Key=eks/123456789012/mycluster/username,Value=johannes \
        Key=eks/123456789012/mycluster/groups,Value=admin
```

Next, we'll create a EKS worker node role:

```
$ aws iam create-role \
    --role-name eks-node \
    --path /eks/ \
    --tags Key=eks/123456789012/mycluster/type,Value=node \
    --assume-role-policy-document file://eks-assume-role-policy-doc.json
```

See the [AWS documentation][eks-node-iam-role] on what kind of assume role policy you need for your EKS worker.

With the user and role created, we can verify that we can find the mappings as expected by running `eks-auth-sync` without the `-commit` flag.
When the `-commit` flag is not set, `eks-auth-sync` will only print out the mappings it has gathered.

```
$ ./eks-auth-sync -config exampleconfig.yaml
users:
- userarn: arn:aws:iam::123456789012:user/johannes
  username: johannes
  groups:
  - admin
roles:
- rolearn: arn:aws:iam::123456789012:role/eks-node
  username: system:node:{{EC2PrivateDNSName}}
  groups:
  - system:bootstrappers
  - system:nodes
```

[eks-node-iam-role]: https://docs.aws.amazon.com/eks/latest/userguide/worker_node_IAM_role.html