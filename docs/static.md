# Reading mappings directly from configuration

If there's a static set of mappings you want to include in your EKS authentication, you can add them using a `static` scanner in `eks-auth-sync`.

You can use the static mappings by adding a scanner entry to your configuration file with the `type` set to `static`.
The `static` section in the scanner must contain two sections: `users` and `roles`.
Those section should be in the same format as the `mapUsers` and `mapRoles` fields in `aws-auth`.

## Example

To use the static mappings, create a YAML file like this:

```yaml
kubernetes:
  kubeConfigPath: ~/.kube/config
scanners:
- type: static
  static:
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
