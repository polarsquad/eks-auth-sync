# Deployment using CronJobs

As the name of the tool indicates, the `eks-auth-sync` tool is meant for synchronizing auth information from external sources to EKS.
However, the tool only synchronizes the information once per call.
In order synchronize periodically, we can use a [Kubernetes CronJob][k8s-cronjob] to run the tool at any schedule we like.

This page guides you through the process of automating EKS auth synchronization using a CronJob.
We'll create a CronJob that runs the `eks-auth-sync` tool within an EKS cluster.
The tool is configured to update the auth configuration located in the cluster.

[k8s-cronjob]: https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/

## Prerequisites

This guide assumes you have an EKS cluster running, and that you have access administrate it.
The name `mycluster` is used as a placeholder for the cluster name.

## IAM policy

If you want `eks-auth-sync` to read auth information from IAM or SSM parameters, you need to grant the tool enough permissions in AWS to do so.
The following IAM policy document grants enough access to read from both sources.

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "ListRoleTags",
            "Effect": "Allow",
            "Action": [
                "iam:ListRoleTags"
            ],
            "Resource": [
                "arn:aws:iam::123456789012:role/*"
            ]
        },
        {
            "Sid": "ListUserTags",
            "Effect": "Allow",
            "Action": [
                "iam:ListUserTags"
            ],
            "Resource": [
                "arn:aws:iam::123456789012:user/*"
            ]
        },
        {
            "Sid": "GetParamter",
            "Effect": "Allow",
            "Action": [
                "ssm:GetParameter"
            ],
            "Resource": [
                "arn:aws:ssm:eu-central-1:123456789012:parameter/eks/mycluster/auth"
            ]
        },
        {
            "Sid": "ListRolesAndUsers",
            "Effect": "Allow",
            "Action": [
                "iam:ListRoles",
                "iam:ListUsers"
            ],
            "Resource": "*"
        }
    ]
}
```

Feel free to customize the permissions to fit your use-case.
You should at least change the account ID (`123456789012`) and IAM user/role and SSM parameter paths to fit your use-case.

You can also remove statements you know you don't need:
For example, if you only want to read auth information from SSM, you can drop the IAM user/role statements.

## IAM role and service account

We'll need a way to pass the IAM role to the Pods that will run `eks-auth-sync`.
The recommended way to do this is to use [IAM roles for service accounts][iam-role-sa]:
a [Kubernetes Service Account] is attached to an IAM role, which means that all AWS API calls done via a Pod using the service account will automatically assume the attached IAM role.
To get started, you'll first need to [enable IAM roles for service accounts on your cluster][iam-role-sa-enable].

Next, [create an IAM role and a service account][iam-role-sa-create] for `eks-auth-sync`.
The IAM role should have the above policy attached to it, and it must be made accessible from the same namespace where `eks-auth-sync` is deployed to.
The rest of this guide will assume you use `eks-auth-sync` as the name for the service account, and `kube-system` as the namespace.

It's also possible to use [kube2iam][] or [kiam][] instead of IAM roles for service accounts, but those are left out of scope for this guide.

[iam-role-sa]: https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html
[iam-role-sa-enable]: https://docs.aws.amazon.com/eks/latest/userguide/enable-iam-roles-for-service-accounts.html
[iam-role-sa-create]: https://docs.aws.amazon.com/eks/latest/userguide/create-service-account-iam-policy-and-role.html
[sa]: https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/
[kube2iam]: https://github.com/jtblin/kube2iam
[kiam]: https://github.com/uswitch/kiam

## RBAC role and permissions

The `eks-auth-sync` tool also needs access to edit the `aws-auth` ConfigMap in the `kube-system` namespace.
We can grant this access using the following RBAC role.

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: aws-auth-editor
  namespace: kube-system
rules:
- apiGroups:
  - ""
  resourceNames:
  - aws-auth
  resources:
  - configmaps
  verbs:
  - get
  - create
  - update
```

Alternatively, you can create the role directly using `kubectl`.

```shell
kubectl -n kube-system \
    create role aws-auth-editor \
    --verb get,create,update \
    --resource=configmaps \
    --resource-name=aws-auth
```

We can grant this access to the service account we created earlier by binding above role to it.

```shell
kubectl -n kube-system \
    create rolebinding eks-auth-sync \
    --role=aws-auth-editor \
    --serviceaccount=kube-system:eks-auth-sync
```

Note that the `kube-system` in the `--serviceaccount` option must match the namespace name where the `eks-auth-sync` service account is located.

## ConfigMap

Next, we need a place to host the configuration file for `eks-auth-sync`.
Let's write the configuration as a [Kubernetes ConfigMap][configmap]

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: eks-auth-sync
  namespace: kube-system
data:
  config.yaml: |
    kubernetes:
      inKubeCluster: true
    scanners:
    - type: iam
      iam:
        clusterName: mycluster
        clusterAccountID: "123456789012"
        pathPrefix: /
    - type: ssm
      ssm:
        path: /eks/mycluster/auth
```

Make sure to replace details such as the `clusterAccountID` to match your environment.
The ConfigMap must be created in the same namespace where the `eks-auth-sync` service account was created.

[configmap]: https://kubernetes.io/docs/concepts/configuration/configmap/

## CronJob

Finally, we can create the CronJob for running the `eks-auth-sync` tool.

```yaml
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: eks-auth-sync
  namespace: kube-system
spec:
  schedule: "*/15 * * * *"
  jobTemplate:
    spec:
      parallelism: 1
      completions: 1
      backoffLimit: 3
      activeDeadlineSeconds: 120
      ttlSecondsAfterFinished: 300
      template:
        spec:
          containers:
          - name: eks-auth-sync
            image: registry.gitlab.com/polarsquad/eks-auth-sync
            imagePullPolicy: Always
            args:
            - -config
            - /etc/eks-auth-sync/config.yaml
            - -commit
            env:
            - name: AWS_REGION
              value: eu-central-1
            resources:
              requests:
                cpu: 100m
                memory: 128Mi
              limits:
                cpu: 500m
                memory: 128Mi
            volumeMounts:
            - name: config
              mountPath: /etc/eks-auth-sync
              readOnly: true
          restartPolicy: Never
          serviceAccountName: eks-auth-sync
          securityContext:
            runAsNonRoot: true
            runAsUser: 10101
            runAsGroup: 10101
            fsGroup: 10101
          volumes:
          - name: config
            configMap:
              name: eks-auth-sync
```

Notes:

* The namespace must the same as where the `eks-auth-sync` service account is located.
* The sync job is scheduled to run every 15 minutes.
  Check out [crontab guru](https://crontab.guru/), if you're new to the cron format.
* Typically, the sync task is fairly quick, so the 2 minute (120 second) deadline should be generous.
* The example here uses the `latest` Docker image tag.
  You should to pin that to a specific tag instead.
  You can find the image tags from the [Gitlab Docker registry](https://gitlab.com/polarsquad/eks-auth-sync/container_registry).
* The AWS region is fed via an environment variable.
  Alternatively, we could feed it via the configuration file in the ConfigMap.
* The resource requests and limits in the example have not been benchmarked.
  Lower settings may work.
* Security context is used for pinning the user and group in the container to [a non-root user and group](https://opensource.com/article/18/3/just-say-no-root-containers).

**WARNING!** Existing `aws-auth` configurations will be overridden by this CronJob, which may lead to you losing some of the existing mappings.
Run the tool *without* the `-commit` flag to print out what mappings it will produce and compare them to the existing mappings in `aws-auth`.
