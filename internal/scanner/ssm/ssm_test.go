package ssm

import (
	"testing"

	intaws "gitlab.com/polarsquad/eks-auth-sync/internal/aws"
	"gitlab.com/polarsquad/eks-auth-sync/test/stub"
	"gitlab.com/polarsquad/eks-auth-sync/test/testdata"

	"github.com/stretchr/testify/assert"
)

var (
	testAWSAPIs = &intaws.API{
		SSM: stub.NewSSM(),
	}
)

func TestSSMScan(t *testing.T) {
	ms, err := Scan(&ScanConfig{testdata.MappingsSSMPath}, testAWSAPIs)

	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, &testdata.AllMappings, ms)
}

func TestSSMScanFail(t *testing.T) {
	_, err := Scan(&ScanConfig{"/nope"}, testAWSAPIs)

	if err == nil {
		t.Fatal("expected an error from SSM scan")
	}
}
