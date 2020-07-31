package ssm

import (
	"encoding/base64"
	"fmt"

	intaws "gitlab.com/polarsquad/eks-auth-sync/internal/aws"
	"gitlab.com/polarsquad/eks-auth-sync/internal/mapping"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type ScanConfig struct {
	Path string `yaml:"path"`
}

func (s *ScanConfig) Validate() error {
	if s.Path == "" {
		return fmt.Errorf("no path specified")
	}
	return nil
}

func Scan(c *ScanConfig, awsAPIs *intaws.API) (ms *mapping.All, err error) {
	ms = &mapping.All{}

	var output *ssm.GetParameterOutput
	output, err = awsAPIs.SSM.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(c.Path),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return
	}

	var bs []byte
	bs, err = base64.StdEncoding.DecodeString(*output.Parameter.Value)
	if err != nil {
		return
	}

	err = ms.FromYAML(bs)
	return
}
