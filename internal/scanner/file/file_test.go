package file

import (
	"github.com/spf13/afero"
	"testing"

	"gitlab.com/polarsquad/eks-auth-sync/internal/mapping"

	"github.com/stretchr/testify/assert"
)

const testMappingsYAML = `
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

var (
	testMappings = &mapping.All{
		Users: []*mapping.User{
			{
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

func TestReadingFile(t *testing.T) {
	mappingsFileName := "test_mappings.yaml"
	fs := afero.NewMemMapFs()
	if err := afero.WriteFile(fs, mappingsFileName, []byte(testMappingsYAML), 0666); err != nil {
		t.Fatal(err)
	}

	ms, err := Scan(&ScanConfig{mappingsFileName}, fs)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, testMappings, ms)
}

func TestReadingNonExistingFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	_, err := Scan(&ScanConfig{"test_mappings_404.yaml"}, fs)

	if err == nil {
		t.Fatal("expected to not find a file")
	}
}
