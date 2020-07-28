package iam

import (
	"fmt"
	"testing"

	"gitlab.com/polarsquad/eks-auth-sync/internal/mapping"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/stretchr/testify/assert"
)

const (
	testClusterName     = "mycluster"
	testCluster2Name    = "notmycluster"
	testAccountID       = "123456789012"
	testAccountID2      = "098765432198"
)

var (
	testClusterTagKeyUsername  = tagKey(eksTagPrefix(testAccountID, testClusterName), tagKeyUsername)
	testClusterTagKeyGroups    = tagKey(eksTagPrefix(testAccountID, testClusterName), tagKeyGroups)
	testClusterTagKeyType      = tagKey(eksTagPrefix(testAccountID, testClusterName), tagKeyType)
	testCluster2TagKeyUsername = tagKey(eksTagPrefix(testAccountID, testCluster2Name), tagKeyUsername)
	testCluster2TagKeyGroups   = tagKey(eksTagPrefix(testAccountID, testCluster2Name), tagKeyGroups)
	testCluster2TagKeyType     = tagKey(eksTagPrefix(testAccountID, testCluster2Name), tagKeyType)
	testAWSConfig              = &AWSConfig{
		STSAPI: &stsStub{accountID: testAccountID2},
		IAMAPI: &iamStub{users: testUsers, roles: testRoles},
	}
)

var testUsers = []*iam.User{
	&iam.User{
		UserName: aws.String("jill@example.org"),
		Tags: []*iam.Tag{
			&iam.Tag{
				Key:   aws.String("role"),
				Value: aws.String("kubernetes admin"),
			},
			&iam.Tag{
				Key:   aws.String(testClusterTagKeyUsername),
				Value: aws.String("jill"),
			},
			&iam.Tag{
				Key:   aws.String(testClusterTagKeyGroups),
				Value: aws.String("admin"),
			},
		},
	},
	&iam.User{
		UserName: aws.String("john@example.org"),
		Tags: []*iam.Tag{
			&iam.Tag{
				Key:   aws.String("role"),
				Value: aws.String("team tead"),
			},
			&iam.Tag{
				Key:   aws.String(testClusterTagKeyUsername),
				Value: aws.String("john"),
			},
			&iam.Tag{
				Key:   aws.String(testClusterTagKeyGroups),
				Value: aws.String("team-x,team-y"),
			},
		},
	},
	&iam.User{
		UserName: aws.String("jack@example.org"),
		Tags: []*iam.Tag{
			&iam.Tag{
				Key:   aws.String("role"),
				Value: aws.String("developer"),
			},
			&iam.Tag{
				Key:   aws.String(testClusterTagKeyUsername),
				Value: aws.String("jack"),
			},
			&iam.Tag{
				Key:   aws.String(testClusterTagKeyGroups),
				Value: aws.String("team-x"),
			},
			&iam.Tag{
				Key:   aws.String(testCluster2TagKeyUsername),
				Value: aws.String("jack"),
			},
			&iam.Tag{
				Key:   aws.String(testCluster2TagKeyGroups),
				Value: aws.String("team-x"),
			},
		},
	},
	&iam.User{
		UserName: aws.String("jan@example.org"),
		Tags: []*iam.Tag{
			&iam.Tag{
				Key:   aws.String("role"),
				Value: aws.String("kubernetes admin"),
			},
			&iam.Tag{
				Key:   aws.String(testCluster2TagKeyUsername),
				Value: aws.String("jan"),
			},
			&iam.Tag{
				Key:   aws.String(testCluster2TagKeyGroups),
				Value: aws.String("admin"),
			},
		},
	},
}

var testRoles = []*iam.Role{
	&iam.Role{
		RoleName: aws.String("eks-node"),
		Tags: []*iam.Tag{
			&iam.Tag{
				Key:   aws.String(testClusterTagKeyType),
				Value: aws.String("node"),
			},
			&iam.Tag{
				Key:   aws.String("purpose"),
				Value: aws.String("EKS node role"),
			},
			&iam.Tag{
				Key:   aws.String("owner"),
				Value: aws.String("admins"),
			},
		},
	},
	&iam.Role{
		RoleName: aws.String("eks-node2"),
		Tags: []*iam.Tag{
			&iam.Tag{
				Key:   aws.String(testCluster2TagKeyType),
				Value: aws.String("node"),
			},
			&iam.Tag{
				Key:   aws.String("purpose"),
				Value: aws.String("EKS node role"),
			},
			&iam.Tag{
				Key:   aws.String("owner"),
				Value: aws.String("admins"),
			},
		},
	},
	&iam.Role{
		RoleName: aws.String("invalid-role"),
		Tags: []*iam.Tag{
			&iam.Tag{
				Key:   aws.String(testClusterTagKeyType),
				Value: aws.String("xxxx"),
			},
			&iam.Tag{
				Key:   aws.String("purpose"),
				Value: aws.String("unknown"),
			},
			&iam.Tag{
				Key:   aws.String("owner"),
				Value: aws.String("admins"),
			},
		},
	},
	&iam.Role{
		RoleName: aws.String("deployer"),
		Tags: []*iam.Tag{
			&iam.Tag{
				Key:   aws.String(testClusterTagKeyType),
				Value: aws.String("user"),
			},
			&iam.Tag{
				Key:   aws.String("purpose"),
				Value: aws.String("machine user for handling deployments"),
			},
			&iam.Tag{
				Key:   aws.String("owner"),
				Value: aws.String("admins"),
			},
			&iam.Tag{
				Key:   aws.String(testClusterTagKeyUsername),
				Value: aws.String("deployer"),
			},
			&iam.Tag{
				Key:   aws.String(testClusterTagKeyGroups),
				Value: aws.String("deployer,team-x"),
			},
		},
	},
	&iam.Role{
		RoleName: aws.String("deployer2"),
		Tags: []*iam.Tag{
			&iam.Tag{
				Key:   aws.String(testCluster2TagKeyType),
				Value: aws.String("user"),
			},
			&iam.Tag{
				Key:   aws.String("purpose"),
				Value: aws.String("machine user for handling deployments"),
			},
			&iam.Tag{
				Key:   aws.String("owner"),
				Value: aws.String("admins"),
			},
			&iam.Tag{
				Key:   aws.String(testCluster2TagKeyUsername),
				Value: aws.String("deployer"),
			},
			&iam.Tag{
				Key:   aws.String(testCluster2TagKeyGroups),
				Value: aws.String("deployer,team-x"),
			},
		},
	},
}

var testUserMappings = []*mapping.User{
	&mapping.User{
		UserARN:  userARN(testAccountID2, "jill@example.org"),
		Username: "jill",
		Groups:   []string{"admin"},
	},
	&mapping.User{
		UserARN:  userARN(testAccountID2, "john@example.org"),
		Username: "john",
		Groups:   []string{"team-x", "team-y"},
	},
	&mapping.User{
		UserARN:  userARN(testAccountID2, "jack@example.org"),
		Username: "jack",
		Groups:   []string{"team-x"},
	},
}

var testRoleMappings = []*mapping.Role{
	mapping.Node(roleARN(testAccountID2, "eks-node")),
	&mapping.Role{
		RoleARN: roleARN(testAccountID2, "deployer"),
		Username: "deployer",
		Groups: []string{"deployer", "team-x"},
	},
}

func TestUserScanning(t *testing.T) {
	scanConfig := &ScanConfig{
		ClusterName:      testClusterName,
		ClusterAccountID: testAccountID,
		PathPrefix:       testPathPrefix,
		DisableRoleScan:  true,
	}

	ms, err := Scan(scanConfig, testAWSConfig)
	if err != nil {
		t.Fatal(err)
	}

	assert.Empty(t, ms.Roles)
	assert.EqualValues(t, testUserMappings, ms.Users)
}

func TestRoleScanning(t *testing.T) {
	scanConfig := &ScanConfig{
		ClusterName:      testClusterName,
		ClusterAccountID: testAccountID,
		PathPrefix:       testPathPrefix,
		DisableUserScan:  true,
	}

	ms, err := Scan(scanConfig, testAWSConfig)
	if err != nil {
		t.Fatal(err)
	}

	assert.Empty(t, ms.Users)
	assert.EqualValues(t, testRoleMappings, ms.Roles)
}

func TestRoleAndUserScanning(t *testing.T) {
	scanConfig := &ScanConfig{
		ClusterName:      testClusterName,
		ClusterAccountID: testAccountID,
		PathPrefix:       testPathPrefix,
	}

	ms, err := Scan(scanConfig, testAWSConfig)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, testUserMappings, ms.Users)
	assert.EqualValues(t, testRoleMappings, ms.Roles)
}

func TestWrongPathScanning(t *testing.T) {
	scanConfig := &ScanConfig{
		ClusterName:      testClusterName,
		ClusterAccountID: testAccountID,
		PathPrefix:       testPathPrefix + "wrong",
	}

	ms, err := Scan(scanConfig, testAWSConfig)
	if err != nil {
		t.Fatal(err)
	}

	assert.Empty(t, ms.Users)
	assert.Empty(t, ms.Roles)
}

func userARN(accountID, username string) string {
	return fmt.Sprintf("arn:aws:iam::%s:user/%s", accountID, username)
}

func roleARN(accountID, rolename string) string {
	return fmt.Sprintf("arn:aws:iam::%s:role/%s", accountID, rolename)
}
