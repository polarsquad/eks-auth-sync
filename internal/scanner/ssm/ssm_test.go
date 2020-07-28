package ssm

import (
	"testing"

	"gitlab.com/polarsquad/eks-auth-sync/internal/mapping"

	"github.com/stretchr/testify/assert"
)

const (
	testSSMPath      = "/path/to/mappings"
	testMappingsYAML = `
users:
- userarn: arn:aws:iam::098765432198:user/john@example.org
  username: john
  groups:
  - admin
roles:
- rolearn: arn:aws:iam::098765432198:role/eks-node
  username: system:node:{{EC2PrivateDNSName}}
  groups:
  - system:bootstrappers
  - system:nodes
`
)

var (
	testAWSConfig = &AWSConfig{
		SSMAPI: &ssmStub{
			contents: map[string]string{
				testSSMPath: testMappingsYAML,
			},
		},
	}
	testMappings = &mapping.All{
		Users: []*mapping.User{
			&mapping.User{
				UserARN:  "arn:aws:iam::098765432198:user/john@example.org",
				Username: "john",
				Groups:   []string{"admin"},
			},
		},
		Roles: []*mapping.Role{
			mapping.Node("arn:aws:iam::098765432198:role/eks-node"),
		},
	}
)

func TestSSMScan(t *testing.T) {
	ms, err := Scan(&ScanConfig{testSSMPath}, testAWSConfig)

	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, testMappings, ms)
}

func TestSSMScanFail(t *testing.T) {
	_, err := Scan(&ScanConfig{testSSMPath + "/nope"}, testAWSConfig)

	if err == nil {
		t.Fatal("expected an error from SSM scan")
	}
}
