package stub

import (
	"gitlab.com/polarsquad/eks-auth-sync/test/testdata"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

type STS struct {
	stsiface.STSAPI
	AccountID string
}

func NewSTS() *STS {
	return &STS{
		AccountID: testdata.AccountID2,
	}
}

func (s *STS) GetCallerIdentity(_ *sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error) {
	return &sts.GetCallerIdentityOutput{
		Account: aws.String(s.AccountID),
	}, nil
}
