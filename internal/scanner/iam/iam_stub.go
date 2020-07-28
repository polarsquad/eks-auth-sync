package iam

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

const (
	testPathPrefix      = "/eks/"
	testUsersListMarker = "moreusers"
	testUserTagsMarker  = "moreusertags"
	testRolesListMarker  = "moreusers"
	testRoleTagsMarker  = "moreusertags"
)

type stsStub struct {
	stsiface.STSAPI
	accountID string
}

func (s *stsStub) GetCallerIdentity(input *sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error) {
	return &sts.GetCallerIdentityOutput{
		Account: aws.String(s.accountID),
	}, nil
}

type iamStub struct {
	iamiface.IAMAPI
	users []*iam.User
	roles []*iam.Role
}

func (i *iamStub) ListUsers(input *iam.ListUsersInput) (output *iam.ListUsersOutput, err error) {
	output = &iam.ListUsersOutput{
		IsTruncated: aws.Bool(false),
	}
	if *input.PathPrefix != testPathPrefix {
		return
	}
	if input.Marker != nil && *input.Marker == testUsersListMarker {
		output.Users = i.users[2:]
		return
	}
	output.IsTruncated = aws.Bool(true)
	output.Marker = aws.String(testUsersListMarker)
	output.Users = i.users[:2]
	return
}

func (i *iamStub) ListUserTags(input *iam.ListUserTagsInput) (output *iam.ListUserTagsOutput, err error) {
	output = &iam.ListUserTagsOutput{
		IsTruncated: aws.Bool(false),
	}
	for _, user := range i.users {
		if *user.UserName == *input.UserName {
			if input.Marker != nil && *input.Marker == testUserTagsMarker {
				output.Tags = user.Tags[2:]
				return
			}
			output.IsTruncated = aws.Bool(true)
			output.Marker = aws.String(testUserTagsMarker)
			output.Tags = user.Tags[:2]
		}
	}
	return
}

func (i *iamStub) ListRoles(input *iam.ListRolesInput) (output *iam.ListRolesOutput, err error) {
	output = &iam.ListRolesOutput{
		IsTruncated: aws.Bool(false),
	}
	if *input.PathPrefix != testPathPrefix {
		return
	}
	if input.Marker != nil && *input.Marker == testRolesListMarker {
		output.Roles = i.roles[2:]
		return
	}
	output.IsTruncated = aws.Bool(true)
	output.Marker = aws.String(testRolesListMarker)
	output.Roles = i.roles[:2]
	return
}

func (i *iamStub) ListRoleTags(input *iam.ListRoleTagsInput) (output *iam.ListRoleTagsOutput, err error) {
	output = &iam.ListRoleTagsOutput{
		IsTruncated: aws.Bool(false),
	}
	for _, role := range i.roles {
		if *role.RoleName == *input.RoleName {
			if input.Marker != nil && *input.Marker == testRoleTagsMarker {
				output.Tags = role.Tags[2:]
				return
			}
			output.IsTruncated = aws.Bool(true)
			output.Marker = aws.String(testRoleTagsMarker)
			output.Tags = role.Tags[:2]
		}
	}
	return
}
