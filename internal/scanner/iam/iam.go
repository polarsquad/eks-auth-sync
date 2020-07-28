package iam

import (
	"fmt"

	"gitlab.com/polarsquad/eks-auth-sync/internal/mapping"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

type ScanConfig struct {
	ClusterName      string
	ClusterAccountID string
	PathPrefix       string
	DisableUserScan  bool
	DisableRoleScan  bool
}

type AWSConfig struct {
	STSAPI stsiface.STSAPI
	IAMAPI iamiface.IAMAPI
}

func AWSConfigFromSession(s *session.Session, c *aws.Config) *AWSConfig {
	return &AWSConfig{
		STSAPI: sts.New(s, c),
		IAMAPI: iam.New(s, c),
	}
}

func (c *ScanConfig) TagPrefix() string {
	return eksTagPrefix(c.ClusterAccountID, c.ClusterName)
}

func eksTagPrefix(clusterAccountID, clusterName string) string {
	return fmt.Sprintf("eks/%s/%s", clusterAccountID, clusterName)
}

func Scan(c *ScanConfig, awsConfig *AWSConfig) (ms *mapping.All, err error) {
	ms = &mapping.All{}

	accountID, err := getAccountID(awsConfig.STSAPI)
	if err != nil {
		return
	}

	if !c.DisableUserScan {
		ms.Users, err = scanIAMUsers(awsConfig.IAMAPI, accountID, c)
		if err != nil {
			return
		}
	}

	if !c.DisableRoleScan {
		ms.Roles, err = scanIAMRoles(awsConfig.IAMAPI, accountID, c)
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
