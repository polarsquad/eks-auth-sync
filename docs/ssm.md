# Reading mappings from an SSM parameter

If you want to update your EKS authentication with mappings from [a parameter in SSM][ssm], you can do it using the `ssm` scanner in `eks-auth-sync`.

You can load mappings from SSM parameters by adding a scanner entry to your configuration file with the `type` set to `ssm`.
The `ssm` section must specify a path to the SSM parameter containing the mappings.

The parameter value should be a base64 encoded string that contains a YAML document with two sections: `users` and `roles`.
Those section should be in the same format as the `mapUsers` and `mapRoles` fields in `aws-auth`.
Base64 encoding is required so that we can store values that contain strings that plain SSM parameter strings don't allow.

[ssm]: https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-parameter-store.html

## Example

To use the `ssm` scanner, create a configuration YAML file like this:

```yaml
aws:
  region: eu-west-1
kubernetes:
  kubeConfigPath: ~/.kube/config
scanners:
- type: ssm
  ssm:
    path: /eks/mycluster/auth
```

This will make `eks-auth-sync` read the mappings from an SSM parameter in path `/eks/mycluster/auth`.
You can of course use whatever path works the best for you.

To place the mappings in a parameter in SSM, we can use the [AWS CLI][awscli].
First, save your mappings to a temporary file:

```
$ cat << EOF > eks-mappings.yaml
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
- rolearn: arn:aws:iam::098765432198:role/eks-fargate-node
  username: system:node:{{SessionName}}
  groups:
  - system:bootstrappers
  - system:nodes
  - system:node-proxier
EOF
```

Next, convert the mappings to Base64, and save it to an environment variable:

```
export EKS_MAPPINGS_BASE64=$(cat eks-mappings.yaml | base64)
```

Now you should be able to create the parameter using AWS CLI.

```
$ aws --region eu-west-1 ssm put-parameter \
    --name /eks/mycluster/auth \
    --type String \
    --value "$EKS_MAPPINGS_BASE64"
```

You can verify that the mappings were saved correctly by running `eks-auth-sync` without the `-commit` flag.
When the `-commit` flag is not set, `eks-auth-sync` will only print out the mappings it has gathered.

```
$ ./eks-auth-sync -config exampleconfig.yaml
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
- rolearn: arn:aws:iam::098765432198:role/eks-fargate-node
  username: system:node:{{SessionName}}
  groups:
  - system:bootstrappers
  - system:nodes
  - system:node-proxier
```

To save those mappings to your EKS cluster, provide the configuration file to `eks-auth-sync`.

```
$ eks-auth-sync -config myconfig.yaml -commit
```

[awscli]: https://aws.amazon.com/cli/