package core

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	intaws "gitlab.com/polarsquad/eks-auth-sync/internal/aws"
	"gitlab.com/polarsquad/eks-auth-sync/internal/k8s"
	"gitlab.com/polarsquad/eks-auth-sync/internal/mapping"
	"gitlab.com/polarsquad/eks-auth-sync/test/stub"
	"gitlab.com/polarsquad/eks-auth-sync/test/testdata"
	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

const (
	testSSMPath      = "/core/test/mappings"
	testMappingsYAML = `
users:
- userarn: arn:aws:iam::123456789012:user/jessica@example.com
  username: jessica
  groups:
  - team-z
roles:
- rolearn: arn:aws:iam::098765432198:role/qa
  username: qa
  groups:
  - qa
`
)

var (
	stubSession = &session.Session{}
	stubAWSAPI  = &intaws.API{
		IAM: stub.NewIAM(),
		SSM: &stub.SSM{
			Contents: map[string]string{
				testSSMPath: testMappingsYAML,
			},
		},
		STS: stub.NewSTS(),
	}
	testMappings = mapping.All{
		Users: []*mapping.User{
			{
				UserARN:  "arn:aws:iam::123456789012:user/jessica@example.com",
				Username: "jessica",
				Groups: []string{
					"team-z",
				},
			},
		},
		Roles: []*mapping.Role{
			{
				RoleARN:  "arn:aws:iam::098765432198:role/qa",
				Username: "qa",
				Groups: []string{
					"qa",
				},
			},
		},
	}
)

func stubCore(input io.Reader, output io.Writer, kubeClient kubernetes.Interface) *Core {
	return &Core{
		AppFS:  afero.NewMemMapFs(),
		Input:  input,
		Output: output,
		AWSSession: func(c *intaws.Config) (*session.Session, error) {
			return stubSession, nil
		},
		AWSAPI: func(s *session.Session, c *intaws.Config) *intaws.API {
			return stubAWSAPI
		},
		KubeClient: func(c *k8s.Config) (kubernetes.Interface, error) {
			return kubeClient, nil
		},
	}
}

func TestIAMScanning(t *testing.T) {
	var input bytes.Buffer
	var output bytes.Buffer
	stubKubeClient := fake.NewSimpleClientset()
	core := stubCore(&input, &output, stubKubeClient)
	configPath := "/config.yaml"
	configFile := `
scanners:
- type: iam
  iam:
    clusterName: mycluster
    clusterAccountID: 123456789012
    pathPrefix: /eks/
    groupSeparator: .
`

	if err := afero.WriteFile(core.AppFS, configPath, []byte(configFile), 0666); err != nil {
		t.Fatal(err)
	}
	if err := core.Run([]string{
		"-config", configPath,
		"-commit",
	}); err != nil {
		t.Fatal(err)
	}
	awsAuth, err := getMappingsFromKube(stubKubeClient)

	assert.Equal(t, &testdata.AllMappings, awsAuth)
	assert.Nil(t, err)
}

func TestIAMAndSSMScanning(t *testing.T) {
	var input bytes.Buffer
	var output bytes.Buffer
	stubKubeClient := fake.NewSimpleClientset()
	core := stubCore(&input, &output, stubKubeClient)
	configFile := `
scanners:
- type: iam
  iam:
    clusterName: mycluster
    clusterAccountID: 123456789012
    pathPrefix: /eks/
    groupSeparator: .
- type: ssm
  ssm:
    path: /core/test/mappings
`
	var expectedMappings mapping.All
	expectedMappings.Append(&testdata.AllMappings)
	expectedMappings.Append(&testMappings)

	input.WriteString(configFile)
	if err := core.Run([]string{
		"-config", "-",
		"-commit",
	}); err != nil {
		t.Fatal(err)
	}
	awsAuth, err := getMappingsFromKube(stubKubeClient)

	assert.Nil(t, err)
	assert.Equal(t, &expectedMappings, awsAuth)
}

func TestFileScanning(t *testing.T) {
	var input bytes.Buffer
	var output bytes.Buffer
	stubKubeClient := fake.NewSimpleClientset()
	configPath := "/path/to/config.yaml"
	core := stubCore(&input, &output, stubKubeClient)
	configFile := `
scanners:
- type: file
  file:
    path: /path/to/mappings.yaml
`

	if err := afero.WriteFile(core.AppFS, configPath, []byte(configFile), 0666); err != nil {
		t.Fatal(err)
	}
	if err := afero.WriteFile(core.AppFS, "/path/to/mappings.yaml", []byte(testMappingsYAML), 0666); err != nil {
		t.Fatal(err)
	}
	if err := core.Run([]string{
		"-config", configPath,
		"-commit",
	}); err != nil {
		t.Fatal(err)
	}
	awsAuth, err := getMappingsFromKube(stubKubeClient)

	assert.Nil(t, err)
	assert.Equal(t, &testMappings, awsAuth)
}

func TestVersionInfoPrint(t *testing.T) {
	var input bytes.Buffer
	var output bytes.Buffer
	stubKubeClient := fake.NewSimpleClientset()
	core := stubCore(&input, &output, stubKubeClient)
	configFile := `
scanners:
- type: ssm
  ssm:
    path: /core/test/mappings
`

	input.WriteString(configFile)
	if err := core.Run([]string{
		"-config", "-",
		"-commit",
		"-version",
	}); err != nil {
		t.Fatal(err)
	}
	awsAuth, err := getMappingsFromKube(stubKubeClient)

	assert.Nil(t, awsAuth)
	assert.Nil(t, err)
}

func TestNoCommit(t *testing.T) {
	var input bytes.Buffer
	var output bytes.Buffer
	stubKubeClient := fake.NewSimpleClientset()
	core := stubCore(&input, &output, stubKubeClient)
	configFile := `
scanners:
- type: ssm
  ssm:
    path: /core/test/mappings
`

	input.WriteString(configFile)
	if err := core.Run([]string{
		"-config", "-",
	}); err != nil {
		t.Fatal(err)
	}
	awsAuth, err := getMappingsFromKube(stubKubeClient)

	assert.Nil(t, awsAuth)
	assert.Nil(t, err)
	assert.Contains(t, output.String(), testMappings.Users[0].Username)
}

func TestNoScanners(t *testing.T) {
	var input bytes.Buffer
	var output bytes.Buffer
	stubKubeClient := fake.NewSimpleClientset()
	core := stubCore(&input, &output, stubKubeClient)
	configFile := `scanners: []`

	input.WriteString(configFile)
	if err := core.Run([]string{
		"-config", "-",
	}); err == nil {
		t.Fatal("expected an error")
	}
	awsAuth, err := getMappingsFromKube(stubKubeClient)

	assert.Nil(t, awsAuth)
	assert.Nil(t, err)
}

func getMappingsFromKube(kubeClient kubernetes.Interface) (*mapping.All, error) {
	var users []*mapping.User
	var roles []*mapping.Role

	cm, err := kubeClient.CoreV1().ConfigMaps("kube-system").Get(context.Background(), "aws-auth", metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal([]byte(cm.Data["mapUsers"]), &users); err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal([]byte(cm.Data["mapRoles"]), &roles); err != nil {
		return nil, err
	}

	return &mapping.All{
		Users: users,
		Roles: roles,
	}, nil
}
