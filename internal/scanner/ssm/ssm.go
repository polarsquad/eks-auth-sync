package ssm

import (
	"gitlab.com/polarsquad/eks-auth-sync/internal/mapping"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

type ScanConfig struct {
	SSMPath string
}

type AWSConfig struct {
	SSMAPI ssmiface.SSMAPI
}

func AWSConfigFromSession(s *session.Session, c *aws.Config) *AWSConfig {
	return &AWSConfig{
		SSMAPI: ssm.New(s, c),
	}
}

func Scan(c *ScanConfig, awsConfig *AWSConfig) (ms *mapping.All, err error) {
	ms = &mapping.All{}
	var output *ssm.GetParameterOutput
	output, err = awsConfig.SSMAPI.GetParameter(&ssm.GetParameterInput{
		Name: aws.String(c.SSMPath),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return
	}
	err = ms.FromYAML([]byte(*output.Parameter.Value))
	return
}
