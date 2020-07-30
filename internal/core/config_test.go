package core

import (
	"bytes"
	"github.com/spf13/afero"
	"gitlab.com/polarsquad/eks-auth-sync/internal/mapping"
	"gitlab.com/polarsquad/eks-auth-sync/internal/scanner/file"
	"gitlab.com/polarsquad/eks-auth-sync/internal/scanner/iam"
	"gitlab.com/polarsquad/eks-auth-sync/internal/scanner/ssm"
	"testing"

	intaws "gitlab.com/polarsquad/eks-auth-sync/internal/aws"
	"gitlab.com/polarsquad/eks-auth-sync/internal/k8s"
	"gitlab.com/polarsquad/eks-auth-sync/internal/scanner"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
)

const testYamlConfig = `
kubernetes:
  inKubeCluster: true
  kubeConfigPath: /path/to/kubeconfig
aws:
  roleARN: arn:aws:iam::098765432198:role/eks-admin
  endpoint: https://localhost:9876
  region: eu-north-1
  disableSSL: true
  maxRetries: 6
scanners:
- name: local
  type: file
  file:
    path: /path/to/mappings.yaml
- name: hardcoded
  type: hardcoded
  mappings:
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
- name: my-account-iam
  type: iam
  iam:
    clusterName: mycluster
    clusterAccountID: 098765432198
    pathPrefix: /eks/
- name: other-account-ssm
  type: ssm
  ssm:
    path: /path/to/mappings
  aws:
    roleARN: arn:aws:iam::123456789012:role/ssm-reader
`

var testConfig = Config{
	Kubernetes: k8s.Config{
		InKubeCluster:  true,
		KubeConfigPath: "/path/to/kubeconfig",
	},
	AWS: intaws.Config{
		RoleARN:    "arn:aws:iam::098765432198:role/eks-admin",
		Endpoint:   "https://localhost:9876",
		Region:     "eu-north-1",
		DisableSSL: aws.Bool(true),
		MaxRetries: aws.Int(6),
	},
	Scanners: []*scanner.Scanner{
		{
			Name: "local",
			Type: "file",
			File: file.ScanConfig{Path: "/path/to/mappings.yaml"},
		},
		{
			Name: "hardcoded",
			Type: "hardcoded",
			Mappings: mapping.All{
				Users: []*mapping.User{
					{
						UserARN:  "arn:aws:iam::098765432198:user/john",
						Username: "john",
						Groups:   []string{"admin"},
					},
				},
				Roles: []*mapping.Role{
					mapping.Node("arn:aws:iam::098765432198:role/eks-node"),
				},
			},
		},
		{
			Name: "my-account-iam",
			Type: "iam",
			IAM: iam.ScanConfig{
				ClusterName:      "mycluster",
				ClusterAccountID: "098765432198",
				PathPrefix:       "/eks/",
			},
		},
		{
			Name: "other-account-ssm",
			Type: "ssm",
			SSM:  ssm.ScanConfig{Path: "/path/to/mappings"},
			AWS: intaws.Config{
				RoleARN: "arn:aws:iam::123456789012:role/ssm-reader",
			},
		},
	},
}

func TestReadingFromYAML(t *testing.T) {
	var config Config
	var buf bytes.Buffer
	buf.WriteString(testYamlConfig)

	if err := config.FromYAML(&buf); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, testConfig, config)
}

func TestReadingYAMLFromFile(t *testing.T) {
	configFileName := "myconfig.yaml"
	fs := afero.NewMemMapFs()
	if err := afero.WriteFile(fs, configFileName, []byte(testYamlConfig), 0666); err != nil {
		t.Fatal(err)
	}
	var config Config

	if err := config.FromYAMLFile(fs, configFileName); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, testConfig, config)
}
