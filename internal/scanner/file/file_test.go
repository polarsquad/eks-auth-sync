package file

import (
	"testing"

	"gitlab.com/polarsquad/eks-auth-sync/internal/mapping"

	"github.com/stretchr/testify/assert"
)

var (
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

func TestReadingFile(t *testing.T) {
	ms, err := Scan(&ScanConfig{"test_mappings.yaml"})

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, testMappings, ms)
}

func TestReadingNonExistingFile(t *testing.T) {
	_, err := Scan(&ScanConfig{"test_mappings_404.yaml"})

	if err == nil {
		t.Fatal("expected to not find a file")
	}
}
