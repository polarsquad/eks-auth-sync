package iam

import (
	"testing"

	intaws "gitlab.com/polarsquad/eks-auth-sync/internal/aws"
	"gitlab.com/polarsquad/eks-auth-sync/test/stub"
	"gitlab.com/polarsquad/eks-auth-sync/test/testdata"

	"github.com/stretchr/testify/assert"
)

var (
	testAWSAPIs = &intaws.API{
		IAM: stub.NewIAM(),
		STS: stub.NewSTS(),
	}
)

func TestUserScanning(t *testing.T) {
	scanConfig := &ScanConfig{
		ClusterName:       testdata.ClusterName,
		ClusterAccountID:  testdata.AccountID,
		PathPrefix:        testdata.PathPrefix,
		GroupSeparatorStr: testdata.GroupSeparator,
		DisableRoleScan:   true,
	}

	ms, err := Scan(scanConfig, testAWSAPIs)
	if err != nil {
		t.Fatal(err)
	}

	assert.Empty(t, ms.Roles)
	assert.EqualValues(t, testdata.UserMappings, ms.Users)
}

func TestRoleScanning(t *testing.T) {
	scanConfig := &ScanConfig{
		ClusterName:       testdata.ClusterName,
		ClusterAccountID:  testdata.AccountID,
		PathPrefix:        testdata.PathPrefix,
		GroupSeparatorStr: testdata.GroupSeparator,
		DisableUserScan:   true,
	}

	ms, err := Scan(scanConfig, testAWSAPIs)
	if err != nil {
		t.Fatal(err)
	}

	assert.Empty(t, ms.Users)
	assert.EqualValues(t, testdata.RoleMappings, ms.Roles)
}

func TestRoleAndUserScanning(t *testing.T) {
	scanConfig := &ScanConfig{
		ClusterName:       testdata.ClusterName,
		ClusterAccountID:  testdata.AccountID,
		PathPrefix:        testdata.PathPrefix,
		GroupSeparatorStr: testdata.GroupSeparator,
	}

	ms, err := Scan(scanConfig, testAWSAPIs)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, testdata.UserMappings, ms.Users)
	assert.EqualValues(t, testdata.RoleMappings, ms.Roles)
}

func TestWrongPathScanning(t *testing.T) {
	scanConfig := &ScanConfig{
		ClusterName:       testdata.ClusterName,
		ClusterAccountID:  testdata.AccountID,
		PathPrefix:        testdata.PathPrefix + "wrong",
		GroupSeparatorStr: testdata.GroupSeparator,
	}

	ms, err := Scan(scanConfig, testAWSAPIs)
	if err != nil {
		t.Fatal(err)
	}

	assert.Empty(t, ms.Users)
	assert.Empty(t, ms.Roles)
}
