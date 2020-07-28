package ssm

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

type ssmStub struct {
	ssmiface.SSMAPI
	contents map[string]string
}

func (s *ssmStub) GetParameter(input *ssm.GetParameterInput) (output *ssm.GetParameterOutput, err error) {
	output = &ssm.GetParameterOutput{}
	value, ok := s.contents[*input.Name]
	if !ok {
		err = awserr.New(ssm.ErrCodeParameterNotFound, "", nil)
		return
	}
	output.Parameter = &ssm.Parameter{
		Name: input.Name,
		Value: aws.String(value),
	}
	return
}

