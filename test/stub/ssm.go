package stub

import (
	"encoding/base64"
	"gitlab.com/polarsquad/eks-auth-sync/test/testdata"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

type SSM struct {
	ssmiface.SSMAPI
	Contents map[string]string
}

func NewSSM() *SSM {
	return &SSM{
		Contents: testdata.SSMContents,
	}
}

func (s *SSM) GetParameter(input *ssm.GetParameterInput) (output *ssm.GetParameterOutput, err error) {
	output = &ssm.GetParameterOutput{}

	value, ok := s.Contents[*input.Name]
	if !ok {
		err = awserr.New(ssm.ErrCodeParameterNotFound, "", nil)
		return
	}

	valueBase64 := base64.StdEncoding.EncodeToString([]byte(value))
	output.Parameter = &ssm.Parameter{
		Name:  input.Name,
		Value: aws.String(valueBase64),
	}
	return
}
