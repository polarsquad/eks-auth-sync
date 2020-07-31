# Reading mappings from a file

If you want to update your EKS authentication with mappings from a file, you can do it using the `file` scanner in `eks-auth-sync`.
This could be useful in situations where you already have a separate process for collecting the mappings, but you just want to use `eks-auth-sync` to deliver them to your EKS cluster.

You can load mappings from a file by adding a scanner entry to your configuration file with the `type` set to `file`.
The `file` section must specify a path to the file containing the mappings.

The mappings file should be in YAML format with two sections: `users` and `roles`.
Those section should be in the same format as the `mapUsers` and `mapRoles` fields in `aws-auth`.

## Example

To use the `file` scanner, create a configuration YAML file like this:

```yaml
kubernetes:
  kubeConfigPath: ~/.kube/config
scanners:
- type: file
  file:
    path: mappings.yaml
```

Here's an example of what the mappings file could contain.

```yaml
users:
- userarn: arn:aws:iam::098765432198:user/john
  username: john
  groups:
  - admin
roles:
- rolearn: arn:aws:iam::098765432198:role/eks-node
  username: system:node:{{EC2PrivateDNSName}}
  groups:
  - system:bootstrappers
  - system:nodes
```

To save those mappings to your EKS cluster, provide the configuration file to `eks-auth-sync`.

```
$ eks-auth-sync -config myconfig.yaml -commit
```
