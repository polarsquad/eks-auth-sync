package mapping

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testAccountID1 = "098765432198"
	testAccountID2 = "123456789012"
)

var testMappings = All{
	Users: []*User{
		{
			UserARN:  userARN(testAccountID1, "john@example.org"),
			Username: "john",
			Groups:   []string{"admin"},
		},
		{
			UserARN:  userARN(testAccountID2, "jill@example.org"),
			Username: "jill",
			Groups:   []string{"team-x"},
		},
		{
			UserARN:  userARN(testAccountID2, "jack@example.org"),
			Username: "jack",
			Groups:   []string{"team-x"},
		},
	},
	Roles: []*Role{
		Node(roleARN(testAccountID1, "eks-node")),
		{
			RoleARN:  roleARN(testAccountID2, "deployer"),
			Username: "deployer",
			Groups:   []string{"team-x", "deployer"},
		},
	},
}

var testUsersYAML = strings.TrimSpace(`
- userarn: arn:aws:iam::098765432198:user/john@example.org
  username: john
  groups:
  - admin
- userarn: arn:aws:iam::123456789012:user/jill@example.org
  username: jill
  groups:
  - team-x
- userarn: arn:aws:iam::123456789012:user/jack@example.org
  username: jack
  groups:
  - team-x
`)

var testRolesYAML = strings.TrimSpace(`
- rolearn: arn:aws:iam::098765432198:role/eks-node
  username: system:node:{{EC2PrivateDNSName}}
  groups:
  - system:bootstrappers
  - system:nodes
- rolearn: arn:aws:iam::123456789012:role/deployer
  username: deployer
  groups:
  - team-x
  - deployer
`)

var testMappingsYAML = []byte(fmt.Sprintf(
	"users:\n%s\nroles:\n%s\n",
	testUsersYAML,
	testRolesYAML,
))

func TestConfigMapGeneration(t *testing.T) {
	cm, err := testMappings.ToConfigMap()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "aws-auth", cm.Name)
	assert.Equal(t, "kube-system", cm.Namespace)
	assert.Equal(t, testUsersYAML, strings.TrimSpace(cm.Data["mapUsers"]))
	assert.Equal(t, testRolesYAML, strings.TrimSpace(cm.Data["mapRoles"]))
}

func TestReadingFromYAML(t *testing.T) {
	var ms All

	if err := ms.FromYAML(testMappingsYAML); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, testMappings, ms)
}

func TestAppend(t *testing.T) {
	var all All
	ms1 := &All{
		Users: testMappings.Users[:2],
	}
	ms2 := &All{
		Roles: testMappings.Roles[:1],
	}
	ms3 := &All{
		Users: testMappings.Users[2:],
		Roles: testMappings.Roles[1:],
	}

	all.Append(ms1)
	all.Append(ms2)
	all.Append(ms3)

	assert.EqualValues(t, &testMappings, &all)
}

func TestIsEmpty(t *testing.T) {
	assert.True(t, (&All{}).IsEmpty())
	assert.False(t, (&All{Users: testMappings.Users}).IsEmpty())
	assert.False(t, (&All{Roles: testMappings.Roles}).IsEmpty())
	assert.False(t, testMappings.IsEmpty())
}

func userARN(accountID, username string) string {
	return fmt.Sprintf("arn:aws:iam::%s:user/%s", accountID, username)
}

func roleARN(accountID, rolename string) string {
	return fmt.Sprintf("arn:aws:iam::%s:role/%s", accountID, rolename)
}
