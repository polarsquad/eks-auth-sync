package iam

import (
	"fmt"

	intaws "gitlab.com/polarsquad/eks-auth-sync/internal/aws"
	"gitlab.com/polarsquad/eks-auth-sync/internal/mapping"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

// DefaultGroupSeparator is the string that is used for distinguishing
// where one group name ends and where the next one starts when parsing
// the Kubernetes groups from the IAM tags.
const DefaultGroupSeparator = "::"

type ScanConfig struct {
	ClusterName       string `yaml:"clusterName"`
	ClusterAccountID  string `yaml:"clusterAccountID"`
	PathPrefix        string `yaml:"pathPrefix"`
	GroupSeparatorStr string `yaml:"groupSeparator"`
	DisableUserScan   bool   `yaml:"disableUserScan"`
	DisableRoleScan   bool   `yaml:"disableRoleScan"`
}

func (c *ScanConfig) Validate() error {
	if c.ClusterName == "" {
		return fmt.Errorf("no cluster name specified")
	}
	if c.ClusterAccountID == "" {
		return fmt.Errorf("no cluster account ID specified")
	}
	return nil
}

func (c *ScanConfig) TagPrefix() string {
	return fmt.Sprintf("eks/%s/%s", c.ClusterAccountID, c.ClusterName)
}

func (c *ScanConfig) GroupSeparator() string {
	if c.GroupSeparatorStr == "" {
		return DefaultGroupSeparator
	}
	return c.GroupSeparatorStr
}

func Scan(c *ScanConfig, awsAPIs *intaws.API) (ms *mapping.All, err error) {
	ms = &mapping.All{}

	accountID, err := getAccountID(awsAPIs.STS)
	if err != nil {
		return
	}

	if !c.DisableUserScan {
		ms.Users, err = scanIAMUsers(awsAPIs.IAM, accountID, c)
		if err != nil {
			return
		}
	}

	if !c.DisableRoleScan {
		ms.Roles, err = scanIAMRoles(awsAPIs.IAM, accountID, c)
		if err != nil {
			return
		}
	}

	return
}

func getAccountID(svc stsiface.STSAPI) (string, error) {
	input := &sts.GetCallerIdentityInput{}
	result, err := svc.GetCallerIdentity(input)
	if err != nil {
		return "", err
	}
	return aws.StringValue(result.Account), nil
}
