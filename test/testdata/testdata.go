package testdata

import (
	"fmt"
	"strings"

	"gitlab.com/polarsquad/eks-auth-sync/internal/mapping"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
)

const (
	PathPrefix      = "/eks/"
	ClusterName     = "mycluster"
	Cluster2Name    = "notmycluster"
	AccountID       = "123456789012"
	AccountID2      = "098765432198"
	MappingsSSMPath = "/path/to/mappings"
	GroupSeparator  = "."
)

var (
	ClusterTagKeyUsername  = tagKey(eksTagPrefix(AccountID, ClusterName), "username")
	ClusterTagKeyGroups    = tagKey(eksTagPrefix(AccountID, ClusterName), "groups")
	ClusterTagKeyType      = tagKey(eksTagPrefix(AccountID, ClusterName), "type")
	Cluster2TagKeyUsername = tagKey(eksTagPrefix(AccountID, Cluster2Name), "username")
	Cluster2TagKeyGroups   = tagKey(eksTagPrefix(AccountID, Cluster2Name), "groups")
	Cluster2TagKeyType     = tagKey(eksTagPrefix(AccountID, Cluster2Name), "type")
)

var Users = []*iam.User{
	{
		UserName: aws.String("jill@example.org"),
		Tags: []*iam.Tag{
			{
				Key:   aws.String("role"),
				Value: aws.String("kubernetes admin"),
			},
			{
				Key:   aws.String(ClusterTagKeyUsername),
				Value: aws.String("jill"),
			},
			{
				Key:   aws.String(ClusterTagKeyGroups),
				Value: aws.String("admin"),
			},
		},
	},
	{
		UserName: aws.String("john@example.org"),
		Tags: []*iam.Tag{
			{
				Key:   aws.String("role"),
				Value: aws.String("team tead"),
			},
			{
				Key:   aws.String(ClusterTagKeyUsername),
				Value: aws.String("john"),
			},
			{
				Key:   aws.String(ClusterTagKeyGroups),
				Value: aws.String(strings.Join([]string{"team-x", "team-y"}, GroupSeparator)),
			},
		},
	},
	{
		UserName: aws.String("jack@example.org"),
		Tags: []*iam.Tag{
			{
				Key:   aws.String("role"),
				Value: aws.String("developer"),
			},
			{
				Key:   aws.String(ClusterTagKeyUsername),
				Value: aws.String("jack"),
			},
			{
				Key:   aws.String(ClusterTagKeyGroups),
				Value: aws.String("team-x"),
			},
			{
				Key:   aws.String(Cluster2TagKeyUsername),
				Value: aws.String("jack"),
			},
			{
				Key:   aws.String(Cluster2TagKeyGroups),
				Value: aws.String("team-x"),
			},
		},
	},
	{
		UserName: aws.String("jan@example.org"),
		Tags: []*iam.Tag{
			{
				Key:   aws.String("role"),
				Value: aws.String("kubernetes admin"),
			},
			{
				Key:   aws.String(Cluster2TagKeyUsername),
				Value: aws.String("jan"),
			},
			{
				Key:   aws.String(Cluster2TagKeyGroups),
				Value: aws.String("admin"),
			},
		},
	},
}

var Roles = []*iam.Role{
	{
		RoleName: aws.String("eks-node"),
		Tags: []*iam.Tag{
			{
				Key:   aws.String(ClusterTagKeyType),
				Value: aws.String("node"),
			},
			{
				Key:   aws.String("purpose"),
				Value: aws.String("EKS node role"),
			},
			{
				Key:   aws.String("owner"),
				Value: aws.String("admins"),
			},
		},
	},
	{
		RoleName: aws.String("eks-node2"),
		Tags: []*iam.Tag{
			{
				Key:   aws.String(Cluster2TagKeyType),
				Value: aws.String("node"),
			},
			{
				Key:   aws.String("purpose"),
				Value: aws.String("EKS node role"),
			},
			{
				Key:   aws.String("owner"),
				Value: aws.String("admins"),
			},
		},
	},
	{
		RoleName: aws.String("eks-fargate-node"),
		Tags: []*iam.Tag{
			{
				Key: aws.String(ClusterTagKeyType),
				Value: aws.String("fargateNode"),
			},
			{
				Key: aws.String("purpose"),
				Value: aws.String("EKS Fargate node role"),
			},
			{
				Key: aws.String("owner"),
				Value: aws.String("admins"),
			},
		},
	},
	{
		RoleName: aws.String("invalid-role"),
		Tags: []*iam.Tag{
			{
				Key:   aws.String(ClusterTagKeyType),
				Value: aws.String("xxxx"),
			},
			{
				Key:   aws.String("purpose"),
				Value: aws.String("unknown"),
			},
			{
				Key:   aws.String("owner"),
				Value: aws.String("admins"),
			},
		},
	},
	{
		RoleName: aws.String("deployer"),
		Tags: []*iam.Tag{
			{
				Key:   aws.String(ClusterTagKeyType),
				Value: aws.String("user"),
			},
			{
				Key:   aws.String("purpose"),
				Value: aws.String("machine user for handling deployments"),
			},
			{
				Key:   aws.String("owner"),
				Value: aws.String("admins"),
			},
			{
				Key:   aws.String(ClusterTagKeyUsername),
				Value: aws.String("deployer"),
			},
			{
				Key:   aws.String(ClusterTagKeyGroups),
				Value: aws.String(strings.Join([]string{"deployer", "team-x"}, GroupSeparator)),
			},
		},
	},
	{
		RoleName: aws.String("deployer2"),
		Tags: []*iam.Tag{
			{
				Key:   aws.String(Cluster2TagKeyType),
				Value: aws.String("user"),
			},
			{
				Key:   aws.String("purpose"),
				Value: aws.String("machine user for handling deployments"),
			},
			{
				Key:   aws.String("owner"),
				Value: aws.String("admins"),
			},
			{
				Key:   aws.String(Cluster2TagKeyUsername),
				Value: aws.String("deployer"),
			},
			{
				Key:   aws.String(Cluster2TagKeyGroups),
				Value: aws.String(strings.Join([]string{"deployer", "team-x"}, GroupSeparator)),
			},
		},
	},
}

var UserMappings = []*mapping.User{
	{
		UserARN:  userARN(AccountID2, "jill@example.org"),
		Username: "jill",
		Groups:   []string{"admin"},
	},
	{
		UserARN:  userARN(AccountID2, "john@example.org"),
		Username: "john",
		Groups:   []string{"team-x", "team-y"},
	},
	{
		UserARN:  userARN(AccountID2, "jack@example.org"),
		Username: "jack",
		Groups:   []string{"team-x"},
	},
}

var RoleMappings = []*mapping.Role{
	mapping.Node(roleARN(AccountID2, "eks-node")),
	mapping.FargateNode(roleARN(AccountID2, "eks-fargate-node")),
	{
		RoleARN:  roleARN(AccountID2, "deployer"),
		Username: "deployer",
		Groups:   []string{"deployer", "team-x"},
	},
}

var AllMappings = mapping.All{
	Users: UserMappings,
	Roles: RoleMappings,
}

var MappingsYAML = `
users:
  - userarn: arn:aws:iam::098765432198:user/jill@example.org
    username: jill
    groups:
      - admin
  - userarn: arn:aws:iam::098765432198:user/john@example.org
    username: john
    groups:
      - team-x
      - team-y
  - userarn: arn:aws:iam::098765432198:user/jack@example.org
    username: jack
    groups:
      - team-x
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
  - rolearn: arn:aws:iam::098765432198:role/deployer
    username: deployer
    groups:
      - deployer
      - team-x
`

var SSMContents = map[string]string{
	MappingsSSMPath: MappingsYAML,
}

func eksTagPrefix(clusterAccountID, clusterName string) string {
	return fmt.Sprintf("eks/%s/%s", clusterAccountID, clusterName)
}

func tagKey(tagPrefix, key string) string {
	return fmt.Sprintf("%s/%s", tagPrefix, key)
}

func userARN(accountID, username string) string {
	return fmt.Sprintf("arn:aws:iam::%s:user/%s", accountID, username)
}

func roleARN(accountID, rolename string) string {
	return fmt.Sprintf("arn:aws:iam::%s:role/%s", accountID, rolename)
}
